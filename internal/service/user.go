package service

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/model"
	"echo-demo/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	UserRepo *repository.UserRepository
	RoleRepo *repository.RoleRepository
}

func NewUserService(ur *repository.UserRepository, rr *repository.RoleRepository) *UserService {
	return &UserService{
		UserRepo: ur,
	}
}

func (us *UserService) Register(request *dto.SignupRequest) (*dto.SignUpResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hash),
		Phone:    request.PhoneNumber,
	}

	tx, err := us.UserRepo.DB.Beginx()
	defer tx.Rollback()

	err = us.UserRepo.InsertUser(tx, &user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = us.RoleRepo.InsertUserRole(tx, user.ID, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.SignUpResponse{
		Email:       request.Email,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
	}, nil
}

func (us *UserService) Login(request *dto.LoginRequest, secret string) (*dto.LoginResponse, error) {
	user, err := us.UserRepo.GetUserByEmail(nil, request.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenstr, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		JWT: tokenstr,
	}, nil
}

func (us *UserService) Profile(id int64) (*dto.SignUpResponse, error) {
	user, err := us.UserRepo.GetUserByID(nil, id)
	if err != nil {
		return nil, err
	}

	return &dto.SignUpResponse{
		Email:       user.Email,
		Name:        user.Name,
		PhoneNumber: user.Phone,
	}, nil
}
