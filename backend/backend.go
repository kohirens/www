package backend

import (
	"net/http"

	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/awslambda"
	"github.com/kohirens/www/storage"
)

type App interface {
	AddRoute(endpoint string, handler Route)
	AddService(key string, service interface{})
	AuthManager() AuthManager
	Decrypt(message []byte) ([]byte, error)
	Encrypt(message []byte) ([]byte, error)
	LoadGPG()
	Name() string
	RouteNotFound(handler Route)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ServeLambda(event *awslambda.Input) (*awslambda.Output, error)
	Service(key string) (interface{}, error)
	ServiceManager() ServiceManager
	TmplManager() TemplateManager
}

const (
	KeyGoogleProvider = "gp"
	KeySessionManager = "sm"

	// MetaRefresh HTML template to redirect the client.
	MetaRefresh = `<!DOCTYPE html>
<html>
	<head><meta http-equiv="refresh" content="0; url='%s'"></head>
	<body></body>
</html>`
	TmplSuffix = "tmpl"

	KeyAccountManager = "am"
	PrefixAccounts    = "accounts"
	PrefixSecrets     = "secrets"
	skLoggedIn        = "loggedIn"
)

var (
	Log     = &logger.Standard{}
	TmplDir = "templates"

	// PublicPages list pages that do not require a client to be logged in.
	PublicPages = []string{
		"/",
		"/api/sign-in",
		"/api/sign-out",
		"/favicon.ico",
	}
	LoginPage = "/"
)

// New A nNew initialized application instance.
func New(
	name string,
	router RouteManager,
	serviceManager ServiceManager,
	tmpl TemplateManager,
	authManager AuthManager,
	store storage.Storage,
) App {
	return &Api{
		name:           name,
		serviceManager: serviceManager,
		router:         router,
		tmplManager:    tmpl,
		authManager:    authManager,
		storage:        store,
	}
}

func NewWithDefaults(name string, store storage.Storage) App {
	return New(
		name,
		NewRouteManager(),
		NewServiceManager(),
		NewTemplateManager(store, TmplDir, TmplSuffix),
		NewAuthManager(),
		store,
	)
}
func NewAccountExec(store storage.Storage) *AccountExec {
	return &AccountExec{store: store}
}
