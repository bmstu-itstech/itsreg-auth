package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/itsreg-auth/internal/domain/entity"
	"github.com/itsreg-auth/internal/domain/value"
	"github.com/itsreg-auth/internal/infrastructure/interfaces"
)

type Service struct {
	userRepo     interfaces.UserRepository
	jwtSecretKey string
}

const JwtAliveTime = 24

func NewAuthService(userRepo interfaces.UserRepository, jwtSecretKey string) *Service {
	return &Service{
		userRepo:     userRepo,
		jwtSecretKey: jwtSecretKey,
	}
}

// TODO: обсудить
// Здесь обязательно нужно сделать так, что те ошибки, которые случились по вине вызывающей стороны
// (например, невалидная почта или ненайденный пользователь) возвращать как Err<что-то там>, а не fmt.Errorf.
// Но как через errors.New(...) дать понять пользователю, что именно в валидации не так?
// Пароль может быть короткий, в нем может не хватать цифры или заглавной буквы, и пользователю это нужно
// конкретно сообщить, что возможно только через вызов fmt.Error, который мы дополнительно обогатили из нашего приложения
// что это ошибка валидации и из errors.Is мы сможем это понять

var (
	ErrValidationUser          = errors.New("the user failed validation")
	ErrUserNotFound            = errors.New("the user does not exist")
	ErrEmailAlreadyUsed        = errors.New("a user with this email already exists")
	ErrFailedComparedPasswords = errors.New("failure to compare passwords")
	ErrInvalidCredentials      = errors.New("invalid credentials")
)

// Register registers a new user
func (s *Service) Register(ctx context.Context, email string, password string) (value.UserId, error) {
	const op = "AuthService.Register"

	user, err := entity.NewUserRegistration(value.Email(email), value.Password(password))
	if err != nil {
		return value.UnknownUserId, fmt.Errorf("%w: %w", ErrValidationUser, err)
	}

	_, err = s.userRepo.Find(ctx, &user)
	if err == nil {
		return value.UnknownUserId, ErrEmailAlreadyUsed
	} else if !errors.Is(err, interfaces.ErrUserNotFound) {
		return value.UnknownUserId, fmt.Errorf("%s: failed to check user in database: %w", op, err)
	}

	userId, err := s.userRepo.Save(ctx, &user)
	if err != nil {
		return value.UnknownUserId, fmt.Errorf("%s: failed to save user: %w", op, err)
	}

	return userId, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "AuthService.Login"

	loginUser, err := entity.NewUserRegistration(value.Email(email), value.Password(password))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrValidationUser, err)
	}

	user, err := s.userRepo.Find(ctx, &loginUser)
	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("%s: failed to find user: %w", op, err)
	}

	match, err := value.CompareHashAndPassword(user.PasswordHash, value.Password(password))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedComparedPasswords, err)
	}
	if !match {
		return "", ErrInvalidCredentials
	}

	token, err := s.generateJWT(user.Id)
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate JWT: %v", op, err)
	}

	return token, nil
}

func (s *Service) generateJWT(userId value.UserId) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(JwtAliveTime * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecretKey))
}
