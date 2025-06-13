package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"itsfriday/internal/util"
	"itsfriday/store"
)

// authentication service

type AuthServiceServer interface {
    SignUp(echo.Context) error
	Login(echo.Context) error
	Logout(echo.Context) error
}

type SignUpRequest struct {
	Email        string `json:"email"`
    Username     string `json:"username"`
	Nickname     string `json:"nickname"`
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

	// check if nickname exists or not
	nickname := request.Nickname
	if nickname == "" {
		nickname = request.Username
	}
	create := &store.User{
		Email:        request.Email,
		Username:     request.Username,
		Nickname:     nickname,
		PasswordHash: string(passwordHash),
		Role:         store.RoleUser,
	}
	if !util.UIDMatcher.MatchString(strings.ToLower(create.Username)) {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid username: %s", create.Username),
		})
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

	expireTime := time.Now().Add(AccessTokenDuration)
	accessToken, err := s.doSignIn(ctx, user, expireTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Unauthenticated,
		    Message: fmt.Sprintf("failed to log in: %v", err),
		})
	}
	
	origin := c.Request().Header.Get("Origin")
	cookie, err := s.buildAccessTokenCookie(accessToken, origin, expireTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("failed to build access token cookie, error: %v", err),
		})
	}
	c.Response().Header().Add("Set-Cookie", cookie)
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

	if err := s.clearAccessTokenCookie(c); err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to set cookie: %v", err),
		})
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

func (s *APIV1Service) clearAccessTokenCookie(c echo.Context) error {
	cookie, err := s.buildAccessTokenCookie("", "", time.Time{})
	if err != nil {
		return fmt.Errorf("failed to build access token cookie: %v", err)
	}

	c.Response().Header().Add("Set-Cookie", cookie)
	return nil
}

func (*APIV1Service) buildAccessTokenCookie(accessToken string, origin string, expireTime time.Time) (string, error) {
	attrs := []string{
		fmt.Sprintf("%s=%s", AccessTokenCookieName, accessToken),
		"Path=/",
		"HttpOnly",
	}
	if expireTime.IsZero() {
		attrs = append(attrs, "Expires=Thu, 01 Jan 1970 00:00:00 GMT")
	} else {
		attrs = append(attrs, "Expires="+expireTime.Format(time.RFC1123))
	}

	isHTTPS := strings.HasPrefix(origin, "https://")
	if isHTTPS {
		attrs = append(attrs, "SameSite=None")
		attrs = append(attrs, "Secure")
	} else {
		if strings.HasPrefix(origin, "http://localhost") {
			attrs = append(attrs, "SameSite=None")
		    attrs = append(attrs, "Secure")
		} else {
            attrs = append(attrs, "SameSite=Strict")
		}	
	}
	return strings.Join(attrs, "; "), nil
}
