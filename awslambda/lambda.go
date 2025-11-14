package awslambda

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www"
	"os"
	"strings"
)

type Handler struct {
	PageSource PageSource
}

type PageSource interface {
	Load(pagePath string) ([]byte, error)
}

const (
	envHttpMethods = "HTTP_METHODS_ALLOWED"
	envRedirectTo  = "REDIRECT_TO"
	headerAltHost  = "viewer-host"
	headerCfDomain = "distribution-domain"
)

var (
	Log = logger.Standard{}
)

func PreliminaryChecks(event *events.LambdaFunctionURLRequest) *events.LambdaFunctionURLResponse {
	method := event.RequestContext.HTTP.Method
	httpAllowedMethods, ok := os.LookupEnv(envHttpMethods)

	if !ok {
		Log.Errf(stderr.MissingEnv, envHttpMethods)
		return www.Respond500().ToLambdaResponse()
	}

	supportedMethods := strings.Split(httpAllowedMethods, ",")

	if strings.ToUpper(method) == "OPTIONS" {
		return www.ResponseOptions(httpAllowedMethods).ToLambdaResponse()
	}

	if www.NotImplemented(method, supportedMethods) {
		return www.Respond501().ToLambdaResponse()
	}

	host := www.GetHeader(event.Headers, headerAltHost)

	doIt, e1 := www.ShouldRedirect(host)
	if e1 != nil {
		Log.Errf("%v", e1.Error())
		return www.Respond500().ToLambdaResponse()
	}

	if doIt {
		serverHost, _ := os.LookupEnv(envRedirectTo)
		if !strings.Contains(serverHost, "https://") {
			serverHost = "https://" + serverHost
		}
		switch method {
		case "POST":
			return www.Respond308(serverHost).ToLambdaResponse()
		}
		return www.Respond301(serverHost).ToLambdaResponse()
	}

	distributionDomain := www.GetHeader(event.Headers, headerCfDomain)

	Log.Infof(stdout.DistDomain, distributionDomain)

	if host == distributionDomain {
		Log.Errf(stderr.DistroRequest, distributionDomain)
		return www.Respond401().ToLambdaResponse()
	}

	Log.Infof("%v", stdout.PreChecks)

	return nil
}
