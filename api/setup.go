package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/routes"
	"github.com/shellhub-io/shellhub/api/routes/middlewares"
	svc "github.com/shellhub-io/shellhub/api/services"
)

func InitializeRoutes(e *echo.Echo, s svc.Service) *echo.Echo {
	handler := routes.NewHandler(s)

	e.Use(middleware.Logger())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apicontext := apicontext.NewContext(s, c)

			return next(apicontext)
		}
	})

	// Public routes for external access through API gateway
	publicAPI := e.Group("/api")

	// Internal routes only accessible by other services in the local container network
	internalAPI := e.Group("/internal")

	internalAPI.GET(routes.AuthRequestURL, apicontext.Handler(handler.AuthRequest), apicontext.Middleware(routes.AuthMiddleware))
	publicAPI.POST(routes.AuthDeviceURL, apicontext.Handler(handler.AuthDevice))
	publicAPI.POST(routes.AuthDeviceURLV2, apicontext.Handler(handler.AuthDevice))
	publicAPI.POST(routes.AuthUserURL, apicontext.Handler(handler.AuthUser))
	publicAPI.POST(routes.AuthUserURLV2, apicontext.Handler(handler.AuthUser))
	publicAPI.GET(routes.AuthUserURLV2, apicontext.Handler(handler.AuthUserInfo))
	internalAPI.GET(routes.AuthUserTokenURL, apicontext.Handler(handler.AuthGetToken))
	publicAPI.POST(routes.AuthPublicKeyURL, apicontext.Handler(handler.AuthPublicKey))
	publicAPI.GET(routes.AuthUserTokenURL, apicontext.Handler(handler.AuthSwapToken))

	publicAPI.PATCH(routes.UpdateUserDataURL, apicontext.Handler(handler.UpdateUserData))
	publicAPI.PATCH(routes.UpdateUserPasswordURL, apicontext.Handler(handler.UpdateUserPassword))
	publicAPI.PUT(routes.EditSessionRecordStatusURL, apicontext.Handler(handler.EditSessionRecordStatus))
	publicAPI.GET(routes.GetSessionRecordURL, apicontext.Handler(handler.GetSessionRecord))

	publicAPI.GET(routes.GetDeviceListURL,
		middlewares.Authorize(apicontext.Handler(handler.GetDeviceList)))
	publicAPI.GET(routes.GetDeviceURL,
		middlewares.Authorize(apicontext.Handler(handler.GetDevice)))
	publicAPI.DELETE(routes.DeleteDeviceURL, apicontext.Handler(handler.DeleteDevice))
	publicAPI.PATCH(routes.RenameDeviceURL, apicontext.Handler(handler.RenameDevice))
	internalAPI.POST(routes.OfflineDeviceURL, apicontext.Handler(handler.OfflineDevice))
	internalAPI.GET(routes.LookupDeviceURL, apicontext.Handler(handler.LookupDevice))
	publicAPI.PATCH(routes.UpdateStatusURL, apicontext.Handler(handler.UpdatePendingStatus))

	publicAPI.POST(routes.CreateTagURL, apicontext.Handler(handler.CreateTag))
	publicAPI.DELETE(routes.DeleteTagURL, apicontext.Handler(handler.DeleteTag))
	publicAPI.PUT(routes.RenameTagURL, apicontext.Handler(handler.RenameTag))
	publicAPI.GET(routes.ListTagURL, apicontext.Handler(handler.ListTag))
	publicAPI.PUT(routes.UpdateTagURL, apicontext.Handler(handler.UpdateTag))

	publicAPI.GET(routes.GetSessionsURL,
		middlewares.Authorize(apicontext.Handler(handler.GetSessionList)))
	publicAPI.GET(routes.GetSessionURL,
		middlewares.Authorize(apicontext.Handler(handler.GetSession)))
	internalAPI.PATCH(routes.SetSessionAuthenticatedURL, apicontext.Handler(handler.SetSessionAuthenticated))
	internalAPI.POST(routes.CreateSessionURL, apicontext.Handler(handler.CreateSession))
	internalAPI.POST(routes.FinishSessionURL, apicontext.Handler(handler.FinishSession))
	internalAPI.POST(routes.RecordSessionURL, apicontext.Handler(handler.RecordSession))
	publicAPI.GET(routes.PlaySessionURL, apicontext.Handler(handler.PlaySession))
	publicAPI.DELETE(routes.RecordSessionURL, apicontext.Handler(handler.DeleteRecordedSession))

	publicAPI.GET(routes.GetStatsURL,
		middlewares.Authorize(apicontext.Handler(handler.GetStats)))

	publicAPI.GET(routes.GetPublicKeysURL, apicontext.Handler(handler.GetPublicKeys))
	publicAPI.POST(routes.CreatePublicKeyURL, apicontext.Handler(handler.CreatePublicKey))
	publicAPI.PUT(routes.UpdatePublicKeyURL, apicontext.Handler(handler.UpdatePublicKey))
	publicAPI.DELETE(routes.DeletePublicKeyURL, apicontext.Handler(handler.DeletePublicKey))
	internalAPI.GET(routes.GetPublicKeyURL, apicontext.Handler(handler.GetPublicKey))
	internalAPI.POST(routes.CreatePrivateKeyURL, apicontext.Handler(handler.CreatePrivateKey))
	internalAPI.POST(routes.EvaluateKeyURL, apicontext.Handler(handler.EvaluateKey))

	publicAPI.GET(routes.ListNamespaceURL, apicontext.Handler(handler.GetNamespaceList))
	publicAPI.GET(routes.GetNamespaceURL, apicontext.Handler(handler.GetNamespace))
	publicAPI.POST(routes.CreateNamespaceURL, apicontext.Handler(handler.CreateNamespace))
	publicAPI.DELETE(routes.DeleteNamespaceURL, apicontext.Handler(handler.DeleteNamespace))
	publicAPI.PUT(routes.EditNamespaceURL, apicontext.Handler(handler.EditNamespace))
	publicAPI.PATCH(routes.AddNamespaceUserURL, apicontext.Handler(handler.AddNamespaceUser))
	publicAPI.PATCH(routes.RemoveNamespaceUserURL, apicontext.Handler(handler.RemoveNamespaceUser))

	return e
}
