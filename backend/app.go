package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/kohirens/sso"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www"
	"github.com/kohirens/www/awslambda"
	"github.com/kohirens/www/session"
	"github.com/kohirens/www/storage"
	"net/http"
)

const (
	cnJWT             = "__JWT__"
	KeyGoogleProvider = "gp"
	KeySessionManager = "sm"
	KeyStorage        = "store"

	// MetaRefresh HTML template to redirect the client.
	MetaRefresh = `<!DOCTYPE html>
<html>
	<head><meta http-equiv="refresh" content="0; url='%s'"></head>
	<body></body>
</html>`
	TmplSuffix = "tmpl"
)

// Api serves as the backend server for managing routes (a.k.a endpoints),
// services, authentication providers, and a template engine.
// These components are available to your routes (which you define). Also,
// these components are replaceable as long as the meet the interface
// requirements.
type Api struct {
	authManager    AuthManager
	router         RouteManager
	serviceManager ServiceManager
	storage        storage.Storage
	tmplManager    TemplateManager
}

type App interface {
	AddRoute(endpoint string, handler Route)
	AddService(key string, service interface{})
	AuthManager() AuthManager
	ServiceManager() ServiceManager
	TmplManager() TemplateManager
	RouteNotFound(handler Route)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ServeLambda(event *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error)
	Service(key string) (interface{}, error)
}

var (
	Log     = &logger.Standard{}
	TmplDir = "templates"
)

// New A nNew initialized application instance.
func New(
	router RouteManager,
	serviceManager ServiceManager,
	tmpl TemplateManager,
	authManager AuthManager,
) App {
	return &Api{
		serviceManager: serviceManager,
		router:         router,
		tmplManager:    tmpl,
		authManager:    authManager,
	}
}

func NewWithDefaults(store storage.Storage) App {
	return &Api{
		authManager:    NewAuthManager(),
		router:         NewRouteManager(),
		serviceManager: NewServiceManager(),
		tmplManager:    NewTemplateManager(store, TmplDir, TmplSuffix),
	}
}

func (a *Api) AddService(key string, service interface{}) {
	a.serviceManager.Add(key, service)
}

// AuthManager Return the authentication manager.
func (a *Api) AuthManager() AuthManager {
	return a.authManager
}

// AddProvider Wrapper method that adds an auth provider to the AuthManager
// for retrieval during request handling.
func (a *Api) AddProvider(key string, provider sso.OIDCProvider) {
	a.authManager.Add(key, provider)
}

// AuthProvider Retrieve an authentication provider from the authentication
// manager.
func (a *Api) AuthProvider(authProvider string) interface{} {
	p, e1 := a.authManager.Get(authProvider)
	if e1 != nil {
		Log.Errf(stderr.AuthProviderLookup, e1.Error())
		return nil
	}
	return p
}

// AddRoute Maps a function to a http.HandlerFunc so that it will respond when
// the route (a.k.a endpoint) is requested.
func (a *Api) AddRoute(endpoint string, handler Route) {
	a.router.Add(endpoint, handler)
}

// RouteNotFound Add a http.HandlerFunc to return a response when a route is
// not found.
func (a *Api) RouteNotFound(handler Route) {
	a.router.NotFound(handler)
}

func (a *Api) Service(key string) (interface{}, error) {
	return a.serviceManager.Get(key)
}

// ServiceManager Get the handler for retrieving services.
func (a *Api) ServiceManager() ServiceManager {
	return a.serviceManager
}

// Session Get the session manager.
func (a *Api) Session() (*session.Manager, error) {
	x, e1 := a.serviceManager.Get(KeySessionManager)
	if e1 != nil {
		return nil, e1
	}
	return x.(*session.Manager), nil
}

// TmplManager Template engine that renders templates.
func (a *Api) TmplManager() TemplateManager {
	return a.tmplManager
}

// ServeHTTP Will be called for every request to this server. There is no need
// to register individual handlers for each pattern or use confusing middleware
// logic. Its responsibilities:
//  1. Initialize/Load an HTTP session for client requests.
//  2. Load logic to process a request and write a response.
//  3. Save the session before sending an HTTP response.
func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Path
	Log.Infof("request %v %v", r.Method, rawPath)

	if e := a.RestoreSessionData(w, r); e != nil {
		Log.Errf(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Add common variables to the template manager.
	a.tmplManager.AppendVars(Variables{
		"HTTP_Method": r.Method,
		"URL_Path":    rawPath,
		"URL_Query":   r.URL.RawQuery,
	})

	// Find the route to respond to the request.
	fn := a.router.Find(rawPath)

	e1 := fn(w, r, a)
	if e1 != nil {
		Log.Errf(e1.Error())

		switch e := e1.(type) {
		case *ReferralError:
			w.Header().Set("Location", e.Location)
			w.WriteHeader(http.StatusSeeOther)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	Log.Infof(stdout.PageDone)

	if e := a.SaveSessionData(w, r); e != nil {
		Log.Errf(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ServeLambda Provide an HTTP response for an AWS Lambda function. Same as
// ServerHTTP but a little extra because the AWS Go library containing
// *events.LambdaFunctionURLRequest it not interchangeable with Go's
// http.Request, and *events.LambdaFunctionURLResponse is not compatible with
// Go's http.Response. They are different patterns that you have to account for.
func (a *Api) ServeLambda(event *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	Log.Infof("handler started")

	if errRes := awslambda.PreliminaryChecks(event); errRes != nil {
		return errRes, nil
	}

	method := event.RequestContext.HTTP.Method
	rawPath := event.RawPath
	w := &www.Response{
		Headers: http.Header{},
	}

	r, e1 := www.NewRequestFromLambdaFunctionURLRequest(event)
	if e1 != nil {
		Log.Errf(e1.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return w.ToLambdaResponse(), nil
	}

	Log.Infof("request %v %v", method, rawPath)

	if e := a.RestoreSessionData(w, r.Request); e != nil {
		Log.Errf(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return w.ToLambdaResponse(), nil
	}

	fn := a.router.Find(rawPath)
	if e := fn(w, r.Request, a); e != nil {
		Log.Errf(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return w.ToLambdaResponse(), nil

	}

	Log.Infof("done loading page")

	if e := a.SaveSessionData(w, r.Request); e != nil {
		Log.Errf(e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return w.ToLambdaResponse(), nil
	}

	return w.ToLambdaResponse(), nil
}

func (a *Api) RestoreSessionData(w http.ResponseWriter, r *http.Request) error {
	sm, e1 := a.Session()
	if e1 != nil {
		return e1
	}

	sm.Load(w, r)

	// TODO pull from the cookie which provider the client chose.
	gp, e2 := a.authManager.Get(KeyGoogleProvider)
	if e2 != nil {
		return e2
	}

	gpData := sm.Get(sso.SessionTokenGoogle)
	if gpData != nil { // restore from the saved session.
		// TODO: Test if you can overwrite members of an initialized struct from a json.Unmarshal.
		//var savedGp *sso.GoogleProvider
		if e := json.Unmarshal(gpData, &gp); e != nil {
			var je *json.UnmarshalTypeError
			if errors.As(e, &je) {
				return fmt.Errorf("json unmarshall error %v %v, offset: %v", je.Field, je.Value, je.Offset)
			}
			return fmt.Errorf(stderr.DecodeJSON, e.Error())
		}
		//gp = savedGp
	}

	return nil
}

func (a *Api) SaveSessionData(w http.ResponseWriter, r *http.Request) error {
	sm, e1 := a.Session()
	if e1 != nil {
		return e1
	}

	// TODO set in the cookie which provider the client chose.

	authProvider, e2 := a.authManager.Get(KeyGoogleProvider)
	if authProvider == nil {
		return e2
	}

	gpData, e3 := json.Marshal(authProvider)
	if e3 != nil {
		return e3
	}

	// When you restore the Google provider from the session the previous token
	// should also be restored.
	sm.Set(sso.SessionTokenGoogle, gpData)

	if e := sm.Save(); e != nil {
		return e
	}

	return nil
}

// Storage Retrieve the storage service from the service manager.
func (a *Api) Storage() storage.Storage {
	return a.storage
}
