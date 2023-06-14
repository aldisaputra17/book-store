package controllers

import (
	"fmt"
	"net/http"

	"github.com/aldisaputra17/book-store/dto"
	"github.com/aldisaputra17/book-store/entities"
	"github.com/aldisaputra17/book-store/helper"
	"github.com/aldisaputra17/book-store/services"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var reqLogin *dto.AuthRequest
	err := ctx.ShouldBind(&reqLogin)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(reqLogin.Email, reqLogin.Password)
	if v, ok := authResult.(entities.User); ok {
		generatedToken := c.jwtService.GenerateToken(v.ID)
		v.Token = generatedToken
		response := helper.BuildResponse(true, "Ok!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var reqRegister *dto.AuthRequest
	errObj := ctx.ShouldBind(&reqRegister)
	if errObj != nil {
		response := helper.BuildErrorResponse("Failed to process request", errObj.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(reqRegister.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate username", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	}
	createdUser, err := c.authService.Register(ctx, reqRegister)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to created", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		fmt.Println("erorr", err)
		return
	} else {
		token := c.jwtService.GenerateToken(createdUser.ID)
		createdUser.Token = token
		response := helper.BuildResponse(true, "Created!", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}
