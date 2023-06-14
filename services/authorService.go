package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aldisaputra17/book-store/dto"
	"github.com/aldisaputra17/book-store/entities"
	"github.com/aldisaputra17/book-store/repositories"
	"github.com/google/uuid"
)

type AuthorService interface {
	Create(ctx context.Context, authorReq *dto.CreateAuthorRequest) (*dto.AuthorResponse, error)
	GetAuthorByCondition(ctx context.Context, bookID string, title string, page int, PageSize int) ([]dto.ReadAuthorResponse, entities.Pagination, error)
	Update(ctx context.Context, authorReq *dto.UpdateAuthorRequest) (*dto.UpdateAuthorResponse, error)
	Delete(ctx context.Context, author entities.Author) error
	FindByID(ctx context.Context, id string) (*dto.ReadAuthorResponse, error)
	IsAllowedToEdit(ctx context.Context, authorID string) bool
}

type authorService struct {
	authorRepository repositories.AuthorRepository
	contextTimeOut   time.Duration
}

func NewAuthorService(authorRepo repositories.AuthorRepository, timeout time.Duration) AuthorService {
	return &authorService{
		authorRepository: authorRepo,
		contextTimeOut:   timeout,
	}
}

func (service *authorService) Create(ctx context.Context, authorReq *dto.CreateAuthorRequest) (*dto.AuthorResponse, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	authorCreate := &entities.Author{
		ID:      id,
		Name:    authorReq.Name,
		Country: authorReq.Country,
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()

	res, err := service.authorRepository.Create(ctx, authorCreate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *authorService) Update(ctx context.Context, authorReq *dto.UpdateAuthorRequest) (*dto.UpdateAuthorResponse, error) {
	authorUpdate := &entities.Author{
		ID:      authorReq.ID,
		Name:    authorReq.Name,
		Country: authorReq.Country,
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()

	res, err := service.authorRepository.Update(ctx, authorUpdate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *authorService) Delete(ctx context.Context, author entities.Author) error {
	return service.authorRepository.Delete(ctx, author)
}

func (service *authorService) FindByID(ctx context.Context, id string) (*dto.ReadAuthorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()
	res, err := service.authorRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *authorService) GetAuthorByCondition(ctx context.Context, bookID string, title string, page int, PageSize int) ([]dto.ReadAuthorResponse, entities.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeOut)
	defer cancel()
	res, pageInfo, err := service.authorRepository.GetAuthorByCondition(ctx, bookID, title, page, PageSize)
	if err != nil {
		return nil, entities.Pagination{}, err
	}
	return res, pageInfo, nil
}

func (service *authorService) IsAllowedToEdit(ctx context.Context, authorID string) bool {
	author, err := service.authorRepository.FindByID(ctx, authorID)
	if err != nil {
		return false
	}
	id := fmt.Sprintf("%v", author.ID)
	return authorID == id
}
