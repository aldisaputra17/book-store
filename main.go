package main

import (
	"fmt"
	"time"

	"github.com/aldisaputra17/book-store/controllers"
	"github.com/aldisaputra17/book-store/database"
	"github.com/aldisaputra17/book-store/middleware"
	"github.com/aldisaputra17/book-store/repositories"
	"github.com/aldisaputra17/book-store/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	contextTimeOut   time.Duration                 = 10 * time.Second
	db               *gorm.DB                      = database.ConnectionDB()
	bookRepository   repositories.BookRepository   = repositories.NewBookRepository(db)
	authorRepository repositories.AuthorRepository = repositories.NewAuthorRepository(db)
	userRepository   repositories.UserRepository   = repositories.NewUserRepository(db)
	bookService      services.BookService          = services.NewBookService(bookRepository, contextTimeOut)
	authorService    services.AuthorService        = services.NewAuthorService(authorRepository, contextTimeOut)
	authService      services.AuthService          = services.NewAuthService(userRepository, contextTimeOut)
	jwtService       services.JWTService           = services.NewJWTService()
	authController   controllers.AuthController    = controllers.NewAuthController(authService, jwtService)
	bookController   controllers.BookController    = controllers.NewBookController(bookService, jwtService)
	authorController controllers.AuthorController  = controllers.NewAuthorController(authorService, jwtService)
)

func main() {
	fmt.Println("Starting Server")
	defer database.CloseDatabaseConnection(db)

	r := gin.Default()

	api := r.Group("api")

	authRoutes := api.Group("/user")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
	}

	bookRoutes := api.Group("/book")
	{
		bookRoutes.POST("", bookController.Create, middleware.AuthorizeJWT(jwtService))
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.GET("", bookController.GetBookByCondition)
		bookRoutes.PUT("", bookController.Update, middleware.AuthorizeJWT(jwtService))
		bookRoutes.DELETE("/:id", bookController.Delete, middleware.AuthorizeJWT(jwtService))
	}
	authorRoutes := api.Group("/author")
	{
		authorRoutes.POST("", authorController.Create, middleware.AuthorizeJWT(jwtService))
		authorRoutes.GET("", authorController.GetAuthorByCondition)
		authorRoutes.GET("/:id", authorController.FindByID)
		authorRoutes.PUT("", authorController.Update, middleware.AuthorizeJWT(jwtService))
		authorRoutes.DELETE("/:id", authorController.Delete, middleware.AuthorizeJWT(jwtService))
	}
	r.Run()
}
