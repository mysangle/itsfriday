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
    group.POST("/user/signup", srv.SignUp)
	group.POST("/user/login", srv.Login)
	group.POST("/user/logout", srv.Logout)
}

func RegisterUserServiceHandler(group *echo.Group, srv UserServiceServer) {
	group.GET("/user/profile", srv.ProfileUser)
	group.PUT("/user/update-user", srv.UpdateUser)
	group.DELETE("/user/delete-user", srv.DeleteUser)
}

func RegisterLibroServiceHandler(group *echo.Group, srv LibroServiceServer) {
	group.POST("/libro/books", srv.CreateBook)
	group.GET("/libro/books/:id", srv.GetBook)
	group.PUT("/libro/books/:id", srv.UpdateBook)
	group.DELETE("/libro/books/:id", srv.DeleteBook)

	group.POST("/libro/reviews", srv.CreateBookReview)
	group.GET("/libro/reviews/:id", srv.GetBookReview)
	group.PUT("/libro/reviews/:id", srv.UpdateBookReview)
	group.DELETE("/libro/reviews/:id", srv.DeleteBookReview)

	group.GET("/libro/dashboard", srv.Dashboard)
	group.GET("/libro/reads", srv.ReadBook)
	group.GET("/libro/report", srv.ReportBook)
	group.GET("/libro/books/:id/reviews", srv.BookReviews)
}

func RegisterDineroServiceHandler(group *echo.Group, srv DineroServiceServer) {
	group.POST("/dinero/categories", srv.CreateDineroCaterory)
	group.PUT("/dinero/categories/:id", srv.UpdateDineroCaterory)
	group.GET("/dinero/categories", srv.ListDineroCaterories)
	group.DELETE("/dinero/categories/:id", srv.DeleteDineroCaterory)
}

func RegisterFitnessServiceHandler(group *echo.Group, srv FitnessServiceServer) {

}

func RegisterFediverseServiceHandler(group *echo.Group, srv FediverseServiceServer) {

}
