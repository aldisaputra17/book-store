package controllers

import (
	"net/http"
	"strconv"

	"github.com/aldisaputra17/book-store/dto"
	"github.com/aldisaputra17/book-store/entities"
	"github.com/aldisaputra17/book-store/helper"
	"github.com/aldisaputra17/book-store/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AuthorController interface {
	Create(ctx *gin.Context)
	GetAuthorByCondition(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByID(ctx *gin.Context)
}

type authorController struct {
	authorService services.AuthorService
	jwtService    services.JWTService
}

func NewAuthorController(authorService services.AuthorService, jwtService services.JWTService) AuthorController {
	return &authorController{
		authorService: authorService,
		jwtService:    jwtService,
	}
}

func (c *authorController) Create(ctx *gin.Context) {
	var (
		reqAuthor *dto.CreateAuthorRequest
		ctxt      = "authorHttpHandler-createAuthor"
	)

	err := ctx.ShouldBind(&reqAuthor)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err json.NewEncoder")
		res := helper.BuildErrorResponse("Failed get object post", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res, err := c.authorService.Create(ctx, reqAuthor)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err json.NewEncoder")
		res := helper.BuildErrorResponse("Failed create author", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result := helper.BuildResponse(true, "success", res)
	ctx.JSON(http.StatusCreated, result)
}

func (c *authorController) Update(ctx *gin.Context) {
	var (
		authorReq *dto.UpdateAuthorRequest
		ctxt      = "authorHttpHandler-updateAuthor"
	)

	errObj := ctx.BindJSON(&authorReq)
	if errObj != nil {
		helper.Log(ctx, log.ErrorLevel, errObj, ctxt, "err json.NewEncoder")
		res := helper.BuildErrorResponse("Failed get object author", errObj.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := c.authorService.Update(ctx, authorReq)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err update author")
		res := helper.BuildErrorResponse("Failed updated author", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result := helper.BuildResponse(true, "Ok", res)
	ctx.JSON(http.StatusOK, result)
}

func (c *authorController) Delete(ctx *gin.Context) {
	var (
		book entities.Author
		ctxt = "authorHttpHandler-deleteAuthor"
		err  error
	)
	id := ctx.Param("id")
	book.ID = uuid.Must(uuid.Parse(id))

	if c.authorService.IsAllowedToEdit(ctx, id) {
		err = c.authorService.Delete(ctx, book)
		if err != nil {
			helper.Log(ctx, log.ErrorLevel, err, ctxt, "err deleted author")
			res := helper.BuildErrorResponse("Failed deleted author", err.Error(), helper.EmptyObj{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		res := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
	}
}

func (c *authorController) FindByID(ctx *gin.Context) {
	var (
		ctxt = "authorHttpHandler-findByIDAuthor"
	)
	id := ctx.Param("id")
	list, err := c.authorService.FindByID(ctx, id)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err fetch author")
		res := helper.BuildErrorResponse("Failed fetch author", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := helper.BuildResponse(true, "Ok", list)
	ctx.JSON(http.StatusOK, res)
}

func (c *authorController) GetAuthorByCondition(ctx *gin.Context) {
	var (
		ctxt = "authorHttpHandler-getByConditionAuthor"
	)
	bookID := ctx.Query("book_id")
	title := ctx.Query("title")
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	lists, total, err := c.authorService.GetAuthorByCondition(ctx, bookID, title, page, pageSize)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err fetch author")
		res := helper.BuildErrorResponse("Failed fetch author", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := helper.BuildReadWithPagination(true, "Ok", lists, total)
	ctx.JSON(http.StatusOK, res)
}
