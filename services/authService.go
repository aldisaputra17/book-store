package services

import (
	"context"
	"log"
	"time"

	"github.com/aldisaputra17/book-store/dto"
	"github.com/aldisaputra17/book-store/entities"
	"github.com/aldisaputra17/book-store/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	Register(ctx context.Context, registerReq *dto.AuthRequest) (*entities.User, error)
	IsDuplicateEmail(email string) bool
	FindByEmail(email string) *entities.User
}

type authService struct {
	userRepository repositories.UserRepository
	contextTimeout time.Duration
}

func NewAuthService(userRepo repositories.UserRepository, time time.Duration) AuthService {
	return &authService{
		userRepository: userRepo,
		contextTimeout: time,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	if v, ok := res.(entities.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false

}

func (service *authService) Register(ctx context.Context, registerReq *dto.AuthRequest) (*entities.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	userCreate := &entities.User{
		ID:        id,
		Email:     registerReq.Email,
		Password:  registerReq.Password,
		CreatedAt: time.Now(),
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()

	res, err := service.userRepository.Create(ctx, userCreate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (service *authService) FindByEmail(email string) *entities.User {
	return service.userRepository.FindByEmail(email)
}
