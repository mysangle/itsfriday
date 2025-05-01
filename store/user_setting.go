package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type UserSetting struct {
	UserID int32
	Key    UserSettingKey
	Value  string
}

func (us *UserSetting) GetAccessTokens() (*UserSettingAccessTokens, error) {
	var tokens UserSettingAccessTokens
    err := json.Unmarshal([]byte(us.Value), &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}
	return &tokens, nil
}

func ConvertUserSettingValueToString(accessTokens *UserSettingAccessTokens) (string, error) {
    value, err := json.Marshal(accessTokens)
	if err != nil {
		return "",  fmt.Errorf("failed to marshal accessTokens: %v", err)
	}
	return string(value), nil
}

type FindUserSetting struct {
	UserID *int32
	Key    UserSettingKey
}

type DeleteUserSetting struct {
	UserID *int32
	Key    UserSettingKey
}

type UserSettingKey int32

const (
	UserSettingKey_USER_SETTING_KEY_UNSPECIFIED UserSettingKey = 0
	// Access tokens for the user.
	UserSettingKey_ACCESS_TOKENS UserSettingKey = 1
	// The locale of the user.
	UserSettingKey_LOCALE UserSettingKey = 2
)

var (
	UserSettingKey_name = map[int32]string{
		0: "USER_SETTING_KEY_UNSPECIFIED",
		1: "ACCESS_TOKENS",
		2: "LOCALE",
	}
	UserSettingKey_value = map[string]int32{
		"USER_SETTING_KEY_UNSPECIFIED": 0,
		"ACCESS_TOKENS":                1,
		"LOCALE":                       2,
	}
)

func (x UserSettingKey) String() string {
	if name, ok := UserSettingKey_name[int32(x)]; ok {
        return name
    }
    return fmt.Sprintf("UNKNOWN_USER_SETTING_KEY(%d)", x)
}

type UserSettingAccessTokens struct {
	AccessTokens []*UserSettingAccessToken `json:"accessTokens"`
}

type UserSettingAccessToken struct {
    AccessToken string `json:"accessToken"`
	Description string `json:"description"`
}

func (s *Store) UpsertUserSetting(ctx context.Context, upsert *UserSetting) (*UserSetting, error) {
	userSetting, err := s.driver.UpsertUserSetting(ctx, upsert)
	if err != nil {
		return nil, err
	}
	if userSetting == nil {
		return nil, errors.New("unexpected nil user setting")
	}
	s.userSettingCache.Store(getUserSettingCacheKey(userSetting.UserID, userSetting.Key.String()), userSetting)
	return userSetting, nil
}

func (s *Store) DeleteUserSetting(ctx context.Context, delete *DeleteUserSetting) (error) {
	err := s.driver.DeleteUserSetting(ctx, delete)
	if err != nil {
		return err
	}

	s.userSettingCache.Delete(getUserSettingCacheKey(*delete.UserID, delete.Key.String()))
	return nil
}

func (s *Store) ListUserSettings(ctx context.Context, find *FindUserSetting) ([]*UserSetting, error) {
    userSettingList, err := s.driver.ListUserSettings(ctx, find)
	if err != nil {
		return nil, err
	}

	for _, userSetting := range userSettingList {
        s.userSettingCache.Store(getUserSettingCacheKey(userSetting.UserID, userSetting.Key.String()), userSetting)
	}
	return userSettingList, nil
}

func (s *Store) GetUserSetting(ctx context.Context, find *FindUserSetting) (*UserSetting, error) {
    if find.UserID != nil {
		if cache, ok := s.userSettingCache.Load(getUserSettingCacheKey(*find.UserID, find.Key.String())); ok {
			userSetting, ok := cache.(*UserSetting)
			if ok {
				return userSetting, nil
			}
		}
	}

	list, err := s.ListUserSettings(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	if len(list) > 1 {
		return nil, fmt.Errorf("expected 1 user setting, but got %d", len(list))
	}

	userSetting := list[0]
	s.userSettingCache.Store(getUserSettingCacheKey(userSetting.UserID, userSetting.Key.String()), userSetting)
	return userSetting, nil
}

func (s *Store) GetUserAccessTokens(ctx context.Context, userID int32) ([]*UserSettingAccessToken, error) {
    userSetting, err := s.GetUserSetting(ctx, &FindUserSetting{
		UserID: &userID,
		Key:    UserSettingKey_ACCESS_TOKENS,
	})
	if err != nil {
		return nil, err
	}
	if userSetting == nil {
		return []*UserSettingAccessToken{}, nil
	}

	accessTokensUserSetting, err := userSetting.GetAccessTokens()
	if err != nil {
		return nil, err
	}
	return accessTokensUserSetting.AccessTokens, nil
}
