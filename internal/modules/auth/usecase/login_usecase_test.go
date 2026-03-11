package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/auth/domain"
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginUsecase(t *testing.T) {
	os.Setenv("AUTH_BCRYPT_COST", "4")
	os.Setenv("AUTH_ACCESS_TOKEN_SECRET_KEY", "secret-key")
	os.Setenv("AUTH_ACCESS_TOKEN_EXPIRES", "24h")
	os.Setenv("AUTH_REFRESH_TOKEN_SECRET_KEY", "secret-key")
	os.Setenv("AUTH_REFRESH_TOKEN_EXPIRES", "48h")

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := utils.GenerateHashPassword(password)
	userAgent := "Mozilla/5.0"
	ipAddress := "127.0.0.1"

	user := domain.User{
		ID:       "user-123",
		Email:    email,
		Password: hashedPassword,
	}

	setupTest := func() (*domain.MockAuthRepository, sqlmock.Sqlmock, domain.AuthUsecase) {
		mockRepo := new(domain.MockAuthRepository)
		mockDB, mockSQL, _ := sqlmock.New()
		logger := zerolog.New(nil)
		u := NewAuthUsecase(mockRepo, &logger, mockDB)
		return mockRepo, mockSQL, u
	}

	t.Run("Success Login", func(t *testing.T) {
		mockRepo, mockSQL, u := setupTest()

		mockRepo.On("ReadDetailUserByEmailQuery", mock.Anything, email).
			Return(user, nil).Once()

		mockRepo.On("CreateSessionQuery", mock.Anything, mock.Anything).
			Return(nil).Maybe()

		mockRepo.On("CreateLoginRecord", mock.Anything, mock.Anything).
			Return(nil).Maybe()

		mockRepo.On("CreateLoginAttemp", mock.Anything, mock.AnythingOfType("domain.LoginAttemp")).
			Return(nil).Maybe()

		mockSQL.ExpectBegin()
		mockSQL.ExpectCommit()
		mockRepo.On("WithTx", mock.Anything).Return(mockRepo).Once()

		session, returnedUser, err := u.LoginUsecase(context.Background(), email, password, userAgent, ipAddress)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, returnedUser.ID)
		assert.NotEmpty(t, session.AccessToken)
		assert.NotEmpty(t, session.RefreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - User Not Found (Wrong Email)", func(t *testing.T) {
		mockRepo, _, u := setupTest()

		mockRepo.On("ReadDetailUserByEmailQuery", mock.Anything, email).
			Return(domain.User{}, sql.ErrNoRows).Maybe()

		mockRepo.On("CreateLoginAttemp", mock.Anything, mock.Anything).
			Return(nil).Maybe()

		_, _, err := u.LoginUsecase(context.Background(), email, password, userAgent, ipAddress)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrWrongEmailOrPassword, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Wrong Password", func(t *testing.T) {
		mockRepo, _, u := setupTest()

		mockRepo.On("ReadDetailUserByEmailQuery", mock.Anything, email).
			Return(user, nil).Maybe()

		mockRepo.On("CreateLoginAttemp", mock.Anything, mock.Anything).
			Return(nil).Maybe()

		_, _, err := u.LoginUsecase(context.Background(), email, "wrongpassword", userAgent, ipAddress)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrWrongEmailOrPassword, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Create Session Failed", func(t *testing.T) {
		mockRepo, mockSQL, u := setupTest()

		mockRepo.On("ReadDetailUserByEmailQuery", mock.Anything, email).
			Return(user, nil).Maybe()

		mockRepo.On("CreateSessionQuery", mock.Anything, mock.Anything).
			Return(errors.New("db error")).Maybe()

		mockRepo.On("CreateLoginAttemp", mock.Anything, mock.Anything).
			Return(nil).Maybe()

		mockSQL.ExpectBegin()
		mockSQL.ExpectRollback()
		mockRepo.On("WithTx", mock.Anything).Return(mockRepo).Once()

		_, _, err := u.LoginUsecase(context.Background(), email, password, userAgent, ipAddress)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrUnknownSource, errors.Unwrap(err))
		mockRepo.AssertExpectations(t)
	})
}
