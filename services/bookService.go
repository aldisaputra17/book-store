package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aldisaputra17/book-store/dto"
	"github.com/aldisaputra17/book-store/entities"
	"github.com/aldisaputra17/book-store/helper"
	"github.com/aldisaputra17/book-store/repositories"
	"github.com/google/uuid"
)

type BookService interface {
	Create(ctx context.Context, bookReq *dto.CreateBookRequest) (*dto.CreateBookResponse, error)
	Update(ctx context.Context, bookReq *dto.UpdateBookRequest) (*dto.UpdateBookResponse, error)
	Delete(ctx context.Context, book entities.Book) error
	FindByID(ctx context.Context, id string) (*dto.ReadBookResponse, error)
	GetBookByCondition(ctx context.Context, authorID string, name string, page int, PageSize int) ([]dto.ReadBookResponse, entities.Pagination, error)
	IsAllowedToEdit(ctx context.Context, bookID string) bool
}

type bookService struct {
	bookRepository repositories.BookRepository
	contextTimeOut time.Duration
}

func NewBookService(bookRepo repositories.BookRepository, timeout time.Duration) BookService {
	return &bookService{
		bookRepository: bookRepo,
		contextTimeOut: timeout,
	}
}

func (service *bookService) Create(ctx context.Context, bookReq *dto.CreateBookRequest) (*dto.CreateBookResponse, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	isbn := helper.GenerateRandomISBN()
	bookCreate := &entities.Book{
		ID:            id,
		Title:         bookReq.Title,
		PublishedYear: time.Now(),
		Isbn:          isbn,
		AuthorID:      bookReq.AuthorID,
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()

	res, err := service.bookRepository.Create(ctx, bookCreate)
	if err != nil {
		return nil, err
	}
	for _, authorID := range bookCreate.AuthorID {
		authorBook := new(entities.AuthorBook)
		authorBook.BookID = bookCreate.ID.String()
		log.Println(authorBook)
		authorBook.AuthorID = authorID

		err := service.bookRepository.AddAuthor(ctx, authorBook)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (service *bookService) Update(ctx context.Context, bookReq *dto.UpdateBookRequest) (*dto.UpdateBookResponse, error) {
	bookUpdate := &entities.Book{
		ID:    bookReq.ID,
		Title: bookReq.Title,
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()

	res, err := service.bookRepository.Update(ctx, bookUpdate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *bookService) Delete(ctx context.Context, book entities.Book) error {
	return service.bookRepository.Delete(ctx, book)
}

func (service *bookService) FindByID(ctx context.Context, id string) (*dto.ReadBookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()
	if id == "" {
		return nil, fmt.Errorf("id not found")
	}
	res, err := service.bookRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *bookService) GetBookByCondition(ctx context.Context, authorID string, name string, page int, PageSize int) ([]dto.ReadBookResponse, entities.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()
	res, pageInfo, err := service.bookRepository.GetBookByCondition(ctx, authorID, name, page, PageSize)
	if err != nil {
		return nil, entities.Pagination{}, err
	}
	return res, pageInfo, nil
}

func (service *bookService) IsAllowedToEdit(ctx context.Context, bookID string) bool {
	book, err := service.bookRepository.FindByID(ctx, bookID)
	if err != nil {
		return false
	}
	id := fmt.Sprintf("%v", book.ID)
	return bookID == id
}
