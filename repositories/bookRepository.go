package repositories

import (
	"context"
	"fmt"

	"github.com/aldisaputra17/book-store/dto"
	"github.com/aldisaputra17/book-store/entities"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(ctx context.Context, book *entities.Book) (*dto.CreateBookResponse, error)
	Update(ctx context.Context, book *entities.Book) (*dto.UpdateBookResponse, error)
	Delete(ctx context.Context, book entities.Book) error
	FindByID(ctx context.Context, id string) (*dto.ReadBookResponse, error)
	AddAuthor(ctx context.Context, authorbook *entities.AuthorBook) error
	GetBookByCondition(ctx context.Context, authorID string, name string, page int, PageSize int) ([]dto.ReadBookResponse, entities.Pagination, error)
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookConnection{
		connection: db,
	}
}

func (db *bookConnection) Create(ctx context.Context, book *entities.Book) (*dto.CreateBookResponse, error) {
	res := db.connection.WithContext(ctx).Create(&book)
	if res.Error != nil {
		return nil, res.Error
	}
	bookRes := &dto.CreateBookResponse{
		ID:            book.ID.String(),
		Title:         book.Title,
		PublishedYear: book.PublishedYear,
		Isbn:          book.Isbn,
	}
	return bookRes, nil
}

func (db *bookConnection) Update(ctx context.Context, book *entities.Book) (*dto.UpdateBookResponse, error) {
	if len(book.ID) == 0 {
		return nil, fmt.Errorf("id not found")
	}
	res := db.connection.WithContext(ctx).Model(&book).Where("id = ?", book.ID).Updates(entities.Book{
		Title: book.Title,
	})
	if res.Error != nil {
		return nil, res.Error
	}
	bookRes := &dto.UpdateBookResponse{
		ID:    book.ID.String(),
		Title: book.Title,
	}
	return bookRes, nil
}

func (db *bookConnection) Delete(ctx context.Context, book entities.Book) error {

	db.connection.Unscoped().Table("author_books").Where("book_id = ?", book.ID).Delete(nil)

	res := db.connection.WithContext(ctx).Delete(book)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (db *bookConnection) AddAuthor(ctx context.Context, authorbook *entities.AuthorBook) error {
	res := db.connection.WithContext(ctx).Create(&authorbook)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (db *bookConnection) FindByID(ctx context.Context, id string) (*dto.ReadBookResponse, error) {
	var book *entities.Book
	res := db.connection.WithContext(ctx).Where("id = ?", id).First(&book)
	if res.Error != nil {
		return nil, res.Error
	}
	authorRes := make([]*dto.AuthorResponse, len(book.Authors))

	for i, author := range book.Authors {
		authorRes[i] = &dto.AuthorResponse{
			ID:      author.ID.String(),
			Name:    author.Name,
			Country: author.Country,
		}
	}
	bookRes := dto.ReadBookResponse{
		ID:            book.ID.String(),
		Title:         book.Title,
		PublishedYear: book.PublishedYear,
		Isbn:          book.Isbn,
		Author:        authorRes,
	}
	return &bookRes, nil
}

func (db *bookConnection) GetBookByCondition(ctx context.Context, authorID string, name string, page int, PageSize int) ([]dto.ReadBookResponse, entities.Pagination, error) {
	var books []*entities.Book
	query := db.connection.Model(&books).
		WithContext(ctx).
		Joins("JOIN author_books ON author_books.book_id = books.id").
		Joins("JOIN authors ON authors.id = author_books.author_id")

	if authorID != "" {
		query = query.Where("author_books.author_id = ?", authorID)
	}

	if name != "" {
		query = query.Where("authors.name LIKE ?", "%"+name+"%")
	}

	var total int64

	query.Count(&total)

	offset := entities.CalculateOffset(page, PageSize)

	query = query.Offset(offset).Limit(PageSize)
	err := query.Preload("Authors").Find(&books).Error
	if err != nil {
		return nil, entities.Pagination{}, err
	}
	pageInfo := entities.CalculatePagination(int(total), page, PageSize)

	bookRes := make([]dto.ReadBookResponse, len(books))

	for i, book := range books {
		authorRes := make([]*dto.AuthorResponse, len(book.Authors))

		for j, author := range book.Authors {
			authorRes[j] = &dto.AuthorResponse{
				ID:      author.ID.String(),
				Name:    author.Name,
				Country: author.Country,
			}
		}
		bookRes[i] = dto.ReadBookResponse{
			ID:            book.ID.String(),
			Title:         book.Title,
			PublishedYear: book.PublishedYear,
			Isbn:          book.Isbn,
			Author:        authorRes,
		}
	}
	return bookRes, pageInfo, nil
}
