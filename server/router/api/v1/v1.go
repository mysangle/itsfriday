package v1

import (
    "github.com/labstack/echo/v4"

	"itsfriday/server/profile"
	"itsfriday/store"
)

type APIV1Service struct {
    Secret  string
	Profile *profile.Profile
	Store   *store.Store
}

func NewAPIV1Service(secret string, profile *profile.Profile, store *store.Store, echoServer *echo.Echo) *APIV1Service {
    apiv1Service := &APIV1Service{
		Secret:     secret,
		Profile:    profile,
		Store:      store,
	}

	group := echoServer.Group("/v1")
	RegisterAuthServiceHandler(group, apiv1Service)

	return apiv1Service
}

func RegisterAuthServiceHandler(group *echo.Group, srv AuthServiceServer) {
    group.POST("/signup", srv.SignUp)
	group.POST("/login", srv.Login)
	group.POST("/logout", srv.Logout)
}
