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

type BookController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	GetBookByCondition(ctx *gin.Context)
}

type bookController struct {
	bookService services.BookService
	jwtService  services.JWTService
}

func NewBookController(bookService services.BookService, jwtService services.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) Create(ctx *gin.Context) {
	var (
		reqBook *dto.CreateBookRequest
		ctxt    = "bookHttpHandler-createBook"
	)

	err := ctx.ShouldBind(&reqBook)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err json.NewEncoder")
		res := helper.BuildErrorResponse("Failed get object post", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res, err := c.bookService.Create(ctx, reqBook)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err json.NewEncoder")
		res := helper.BuildErrorResponse("Failed get object post", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result := helper.BuildResponse(true, "Created", res)
	ctx.JSON(http.StatusCreated, result)
}

func (c *bookController) Update(ctx *gin.Context) {
	var (
		bookReq *dto.UpdateBookRequest
		ctxt    = "bookHttpHandler-updateBook"
	)

	errObj := ctx.BindJSON(&bookReq)
	if errObj != nil {
		helper.Log(ctx, log.ErrorLevel, errObj, ctxt, "err json.NewEncoder")
		res := helper.BuildErrorResponse("Failed get object post", errObj.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := c.bookService.Update(ctx, bookReq)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err json.NewEncoder")
		res := helper.BuildErrorResponse("Failed get object post", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result := helper.BuildResponse(true, "Ok", res)
	ctx.JSON(http.StatusOK, result)
}

func (c *bookController) Delete(ctx *gin.Context) {
	var (
		book entities.Book
		ctxt = "bookHttpHandler-deleteBook"
		err  error
	)
	id := ctx.Param("id")
	book.ID = uuid.Must(uuid.Parse(id))

	if c.bookService.IsAllowedToEdit(ctx, id) {
		err = c.bookService.Delete(ctx, book)
		if err != nil {
			helper.Log(ctx, log.ErrorLevel, err, ctxt, "err deleted book")
			res := helper.BuildErrorResponse("Failed deleted book", err.Error(), helper.EmptyObj{})
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

func (c *bookController) FindByID(ctx *gin.Context) {
	var (
		ctxt = "bookHttpHandler-findByIDBook"
	)
	id := ctx.Param("id")
	list, err := c.bookService.FindByID(ctx, id)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err fetch book")
		res := helper.BuildErrorResponse("Failed fetch book", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := helper.BuildResponse(true, "Ok", list)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) GetBookByCondition(ctx *gin.Context) {
	var (
		ctxt = "bookHttpHandler-getBookByAuthorBook"
	)
	authorID := ctx.Query("author_id")
	name := ctx.Query("name")
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

	books, total, err := c.bookService.GetBookByCondition(ctx, authorID, name, page, pageSize)
	if err != nil {
		helper.Log(ctx, log.ErrorLevel, err, ctxt, "err fetch book")
		res := helper.BuildErrorResponse("Failed fecth book", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := helper.BuildReadWithPagination(true, "Ok", books, total)
	ctx.JSON(http.StatusOK, res)
}
