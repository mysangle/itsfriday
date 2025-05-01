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
	RegisterUserServiceHandler(group, apiv1Service)
	RegisterLibroServiceHandler(group, apiv1Service)
	RegisterDineroServiceHandler(group, apiv1Service)
	RegisterFitnessServiceHandler(group, apiv1Service)
	RegisterFediverseServiceHandler(group, apiv1Service)

	return apiv1Service
}

func RegisterAuthServiceHandler(group *echo.Group, srv AuthServiceServer) {
    group.POST("/signup", srv.SignUp)
	group.POST("/login", srv.Login)
	group.POST("/logout", srv.Logout)
}

func RegisterUserServiceHandler(group *echo.Group, srv UserServiceServer) {
	group.GET("/profile", srv.ProfileUser)
	group.PUT("/update-user", srv.UpdateUser)
	group.DELETE("/delete-user", srv.DeleteUser)
}

func RegisterLibroServiceHandler(group *echo.Group, srv LibroServiceServer) {
	group.POST("/libro", srv.CreateLibro)
	group.GET("/libro/:id", srv.GetLibro)
	group.PUT("/libro/:id", srv.UpdateLibro)
	group.DELETE("/libro/:id", srv.DeleteLibro)

	group.POST("/libro/review", srv.CreateLibroReview)
	group.GET("/libro/review/:id", srv.GetLibroReview)
	group.PUT("/libro/review/:id", srv.UpdateLibroReview)
	group.DELETE("/libro/review/:id", srv.DeleteLibroReview)

	group.GET("/libro/read", srv.ReadLibroReview)
	group.GET("/libro/report", srv.ReportLibroReview)
}

func RegisterDineroServiceHandler(group *echo.Group, srv DineroServiceServer) {

}

func RegisterFitnessServiceHandler(group *echo.Group, srv FitnessServiceServer) {

}

func RegisterFediverseServiceHandler(group *echo.Group, srv FediverseServiceServer) {

}
