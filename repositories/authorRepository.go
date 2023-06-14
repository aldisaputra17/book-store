package repositories

import (
	"context"

	"github.com/aldisaputra17/book-store/dto"
	"github.com/aldisaputra17/book-store/entities"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	Create(ctx context.Context, author *entities.Author) (*dto.AuthorResponse, error)
	GetAuthorByCondition(ctx context.Context, bookID string, title string, page int, PageSize int) ([]dto.ReadAuthorResponse, entities.Pagination, error)
	Update(ctx context.Context, author *entities.Author) (*dto.UpdateAuthorResponse, error)
	Delete(ctx context.Context, author entities.Author) error
	FindByID(ctx context.Context, id string) (*dto.ReadAuthorResponse, error)
}

type authorConnection struct {
	connection *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorConnection{
		connection: db,
	}
}

func (db *authorConnection) Create(ctx context.Context, author *entities.Author) (*dto.AuthorResponse, error) {
	res := db.connection.WithContext(ctx).Create(&author)
	if res.Error != nil {
		return nil, res.Error
	}
	authorRes := &dto.AuthorResponse{
		ID:      author.ID.String(),
		Name:    author.Name,
		Country: author.Country,
	}
	return authorRes, nil
}

func (db *authorConnection) Update(ctx context.Context, author *entities.Author) (*dto.UpdateAuthorResponse, error) {
	res := db.connection.WithContext(ctx).Model(&author).Where("id = ?", author.ID).Updates(&entities.Author{
		Name: author.Name, Country: author.Country,
	})
	if res.Error != nil {
		return nil, res.Error
	}
	authorRes := &dto.UpdateAuthorResponse{
		ID:      author.ID.String(),
		Name:    author.Name,
		Country: author.Country,
	}
	return authorRes, nil
}

func (db *authorConnection) Delete(ctx context.Context, author entities.Author) error {
	db.connection.Unscoped().Table("author_books").Where("author_id = ?", author.ID).Delete(nil)
	res := db.connection.WithContext(ctx).Delete(author)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (db *authorConnection) FindByID(ctx context.Context, id string) (*dto.ReadAuthorResponse, error) {
	var author *entities.Author
	res := db.connection.WithContext(ctx).Where("id = ?", id).First(&author)
	if res.Error != nil {
		return nil, res.Error
	}
	bookRes := make([]*dto.CreateBookResponse, len(author.Books))

	for i, book := range author.Books {
		bookRes[i] = &dto.CreateBookResponse{
			ID:            book.ID.String(),
			Title:         book.Title,
			PublishedYear: book.PublishedYear,
			Isbn:          book.Isbn,
		}
	}
	authorRes := dto.ReadAuthorResponse{
		ID:      author.ID.String(),
		Name:    author.Name,
		Country: author.Country,
		Book:    bookRes,
	}
	return &authorRes, nil
}

func (db *authorConnection) GetAuthorByCondition(ctx context.Context, bookID string, title string, page int, PageSize int) ([]dto.ReadAuthorResponse, entities.Pagination, error) {
	var authors []*entities.Author

	query := db.connection.Model(&authors).
		WithContext(ctx).
		Joins("JOIN author_books ON author_books.author_id = authors.id").
		Joins("JOIN books ON books.id = author_books.book_id")
	if bookID != "" {
		query = query.Where("author_books.book_id = ?", bookID)
	}
	if title != "" {
		query = query.Where("books.title LIKE ?", "%"+title+"%")
	}
	var total int64

	query.Count(&total)

	offset := entities.CalculateOffset(page, PageSize)

	query = query.Offset(offset).Limit(PageSize)
	err := query.Preload("Books").Find(&authors).Error
	if err != nil {
		return nil, entities.Pagination{}, err
	}
	pageInfo := entities.CalculatePagination(int(total), page, PageSize)

	authorRes := make([]dto.ReadAuthorResponse, len(authors))

	for i, author := range authors {
		bookRes := make([]*dto.CreateBookResponse, len(author.Books))

		for j, book := range author.Books {
			bookRes[j] = &dto.CreateBookResponse{
				ID:            book.ID.String(),
				Title:         book.Title,
				PublishedYear: book.PublishedYear,
				Isbn:          book.Isbn,
			}
		}
		authorRes[i] = dto.ReadAuthorResponse{
			ID:      author.ID.String(),
			Name:    author.Name,
			Country: author.Country,
			Book:    bookRes,
		}
	}
	return authorRes, pageInfo, nil
}
