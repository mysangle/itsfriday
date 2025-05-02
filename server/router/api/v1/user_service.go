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

// user service

type UserServiceServer interface {
	ProfileUser(echo.Context) error
	UpdateUser(echo.Context) error
	DeleteUser(echo.Context) error
}

type User struct {
	ID           int32             `json:"id"`
	State        store.RowStatus   `json:"state"`
	CreatedTime  int64             `json:"createdTime"`
	UpdatedTime  int64             `json:"updatedTime"`
	Username     string            `json:"username"`
	Role         store.Role        `json:"role"`
	Email        string            `json:"email"`
	Nickname     string            `json:"nickname"`
	AvatarURL    string            `json:"avatarUrl"`
	Description  string            `json:"description"`
}

type UpdateUserRequest struct {
	Email           string `json:"email"`
    Nickname        string `json:"nickname"`
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	Description     string `json:"description"`
}

type DeleteUserRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type DeleteUserAccessToken struct {
	ID           int32
	AccessToken  string
}

func (s *APIV1Service) ProfileUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &userID})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get user: %v", err),
		})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "user not found",
		})
	}

	userInfo := convertUserFromStore(user)
	return c.JSON(http.StatusOK, userInfo)
}

func (s *APIV1Service) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(UpdateUserRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid update user request: %v", err),
		})
	}

	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &userID})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get user: %v", err),
		})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "user not found",
		})
	}

	currentTs := time.Now().Unix()
	update := &store.UpdateUser{
		ID:        user.ID,
		UpdatedTs: &currentTs,
	}
	if request.Email != "" {
		update.Email = &request.Email
	}
	if request.Nickname != "" {
		update.Nickname = &request.Nickname
	}
	if request.OldPassword != "" && request.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.OldPassword)); err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrorResponse{
				Code:    InvalidRequest,
				Message: fmt.Sprintf("unmatched old password: %v", err),
			})
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	    if err != nil {
		    return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			    Code:    Internal,
		        Message: fmt.Sprintf("failed to generate new password hash: %v", err),
		    })
	    }
		passwordHashStr := string(passwordHash)
	    update.PasswordHash = &passwordHashStr
	}
	if request.Description != "" {
		update.Description = &request.Description
	}

	updatedUser, err := s.Store.UpdateUser(ctx, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to update user: %v", err),
		})
	}

	userInfo := convertUserFromStore(updatedUser)
	return c.JSON(http.StatusOK, userInfo)
}

func (s *APIV1Service) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(DeleteUserRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid delete request: %v", err),
		})
	}

	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &userID})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get user: %v", err),
		})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "user not found",
		})
	}
	if user.Username != request.Username {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "permission denied",
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("unmatched username and password: %v", err),
		})
	}

	if err := s.Store.DeleteUserSetting(ctx, &store.DeleteUserSetting{
		UserID:  &user.ID,
		Key:     store.UserSettingKey_ACCESS_TOKENS,
	}); err != nil {
		slog.Error("failed to delete access tokens before deleting user", "error", err)
	}

	if err := s.Store.DeleteUser(ctx, &store.DeleteUser{
		ID: user.ID,
	}); err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to delete user: %v", err),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *APIV1Service) DeleteUserAccessToken(ctx context.Context, request *DeleteUserAccessToken) (string, error) {
	userAccessTokens, err := s.Store.GetUserAccessTokens(ctx, request.ID)
	if err != nil {
		return "", fmt.Errorf("failed to list access tokens: %v", err)
	}
	updatedUserAccessTokens := []*store.UserSettingAccessToken{}
	for _, userAccessToken := range userAccessTokens {
		if userAccessToken.AccessToken == request.AccessToken {
			continue
		}
		updatedUserAccessTokens = append(updatedUserAccessTokens, userAccessToken)
	}
	println(updatedUserAccessTokens)
	value, err := store.ConvertUserSettingValueToString(&store.UserSettingAccessTokens{
		AccessTokens: updatedUserAccessTokens,
	})
	if err != nil {
		return "", fmt.Errorf("failed to convert user setting to string: %v", err)
	}
	if _, err := s.Store.UpsertUserSetting(ctx, &store.UserSetting{
		UserID: request.ID,
		Key:    store.UserSettingKey_ACCESS_TOKENS,
		Value: value,
	}); err != nil {
		return "", fmt.Errorf("failed to upsert user setting: %v", err)
	}

	return "", nil
}

func (s *APIV1Service) UpsertAccessTokenToStore(ctx context.Context, user *store.User, accessToken, description string) error {
    userAccessTokens, err := s.Store.GetUserAccessTokens(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user access tokens: %w", err)
	}
	userAccessToken := store.UserSettingAccessToken{
		AccessToken: accessToken,
		Description: description,
	}
	userAccessTokens = append(userAccessTokens, &userAccessToken)
	value, err := store.ConvertUserSettingValueToString(&store.UserSettingAccessTokens{
		AccessTokens: userAccessTokens,
	})
	if err != nil {
		return fmt.Errorf("failed to convert user setting to string: %v", err)
	}
	if _, err := s.Store.UpsertUserSetting(ctx, &store.UserSetting{
		UserID: user.ID,
		Key:    store.UserSettingKey_ACCESS_TOKENS,
		Value:  value,
	}); err != nil {
		return fmt.Errorf("failed to upsert user setting: %v", err)
	}
	return nil
}

func convertUserFromStore(user *store.User) *User {
    userInfo := &User{
        ID: user.ID,
		State: user.RowStatus,
		CreatedTime: user.CreatedTs,
		UpdatedTime: user.UpdatedTs,
		Username: user.Username,
		Role: user.Role,
		Email: user.Email,
		Nickname: user.Nickname,
		AvatarURL: user.AvatarURL,
		Description: user.Description,
	}
	return userInfo
}
