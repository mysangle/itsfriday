package v1

import (
	"context"
	"fmt"

	"itsfriday/store"
)

type User struct {
	ID           int32
	State        store.RowStatus
	CreatedTime  int64
	UpdatedTime  int64
	Username     string
	Role         store.Role
	Email        string
	Nickname     string
	AvatarURL    string
	Description  string
}

type DeleteUserAccessToken struct {
	ID           int32
	AccessToken  string
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
