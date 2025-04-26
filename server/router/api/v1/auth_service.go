package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"itsfriday/store"
)

type AuthServiceServer interface {
    SignUp(echo.Context) error
	Login(echo.Context) error
	Logout(echo.Context) error
}

type SignUpRequest struct {
	Email        string `json:"email"`
    Username     string `json:"username"`
	Password     string `json:"password"`
	AvatarUrl    string `json:"avatar_url"`
}

type LoginRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type AccessTokenInfo struct {
	AccessToken  string `json:"accessToken"`
}

func (s *APIV1Service) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
    request := new(SignUpRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid signup request: %v", err),
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
		    Message: fmt.Sprintf("failed to generate password hash: %v", err),
		})
	}

	create := &store.User{
		Email:        request.Email,
		Username:     request.Username,
		Nickname:     request.Username,
		PasswordHash: string(passwordHash),
		Role:         store.RoleUser,
	}
	user, err := s.Store.CreateUser(ctx, create)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
		    Message: fmt.Sprintf("failed to create user: %v", err),
		})
	}

	userInfo := convertUserFromStore(user)
	return c.JSON(http.StatusOK, userInfo)
}

func (s *APIV1Service) Login(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(LoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid login request: %v", err),
		})
	}

	var findUser store.FindUser
	if request.Email != "" {
		findUser = store.FindUser{
			Email: &request.Email,
		}
	} else if request.Username != "" {
		findUser = store.FindUser{
			Username: &request.Username,
		}
	} else {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: "email or username should be provided",
		})
	}
	user, err := s.Store.GetUser(ctx, &findUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
		    Message: fmt.Sprintf("failed to get user: %v", err),
		})
	}
	if user == nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("unmatched username and password: %v", err),
		})
	}

	// Compare the stored hashed password, with the hashed version of the password that was received.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("unmatched username and password: %v", err),
		})
	}

	accessToken, err := s.doSignIn(ctx, user, time.Now().Add(AccessTokenDuration))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Unauthenticated,
		    Message: fmt.Sprintf("failed to log in: %v", err),
		})
	}

	return c.JSON(http.StatusOK, &AccessTokenInfo{
		AccessToken: accessToken,
	})
}

func (s *APIV1Service) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	
	accessToken, ok := c.Get(accessTokenContextKey).(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get access token",
		})
	}

	// try to delete the access token from the store.
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	user, err := s.Store.GetUser(ctx, &store.FindUser{
		ID: &userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get user: %v", err),
		})
	}
	if user != nil {
		if _, err := s.DeleteUserAccessToken(ctx, &DeleteUserAccessToken{
			ID:        user.ID,
			AccessToken: accessToken,
		}); err != nil {
			slog.Error("failed to delete access token", "error", err)
			return c.JSON(http.StatusInternalServerError, &ErrorResponse{
				Code:    Internal,
				Message: fmt.Sprintf("failed to delete access token: %v", err),
			})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *APIV1Service) doSignIn(ctx context.Context, user *store.User, expireTime time.Time) (string, error) {
	accessToken, err := GenerateAccessToken(user.Username, user.ID, expireTime, []byte(s.Secret))
	if err != nil {
		slog.Error("failed to generate access token: ", "error", err)
		return "", err
	}
	if err := s.UpsertAccessTokenToStore(ctx, user, accessToken, "user login"); err != nil {
		return "", fmt.Errorf("failed to upsert access token to store, error: %v", err)
	}

	return accessToken, nil
}
