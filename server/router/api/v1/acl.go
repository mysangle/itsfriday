package v1

import (
	"context"
	"log/slog"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"

	"itsfriday/store"
	"itsfriday/internal/util"
)

const (
	useridContextKey      string      = "userid"
	accessTokenContextKey string      = "access-token"

	InvalidUserID int32   = 0
)

type authHandler struct {
	Store      *store.Store
	secret     string
	ContextKey string
}

func NewAuthHandler(store *store.Store, secret string, contextKey string) *authHandler {
	return &authHandler{
		Store:  store,
		secret: secret,
		ContextKey: contextKey,
	}
}

func (ai *authHandler) ErrorHandler(c echo.Context, err error) error {
	req := c.Request()
	path := req.URL.Path
	if isUnauthorizeAllowedMethod(path) {
		return nil
	}

	if e, ok := err.(*echojwt.TokenParsingError); ok {
		return e
	}

	cookie, cookieErr := c.Cookie(AccessTokenCookieName)
    if cookieErr != nil {
		// show error from authorization header if both authorization header and cookie header do not exist
		slog.Error("failed to get access token", "error", err)
        return err
    }
	token, err := ai.ParseTokenFunc(c, cookie.Value)
	if err != nil {
		return err
	}
	c.Set(ai.ContextKey, token)
	return nil
}

func (ai *authHandler) ParseTokenFunc(c echo.Context, auth string) (interface{}, error) {
	ctx := c.Request().Context()
	token, userID, err := ai.authenticate(ctx, auth)
	if err != nil {
		return nil, &echojwt.TokenError{Token: token, Err: err}
	}

	if userID != InvalidUserID {
		c.Set(useridContextKey, userID)
	}
	c.Set(accessTokenContextKey, auth)
	return token, nil
}

func (ai *authHandler) authenticate(ctx context.Context, accessToken string) (*jwt.Token, int32, error) {
    if accessToken == "" {
		return nil, 0, errors.New("access token not found")
	}
	claims := &ClaimsMessage{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("unexpected access token signing method=%v, expect %v", t.Header["alg"], jwt.SigningMethodHS256)
		}
		if kid, ok := t.Header["kid"].(string); ok {
			if kid == "v1" {
				return []byte(ai.secret), nil
			}
		}
		return nil, fmt.Errorf("unexpected access token kid=%v", t.Header["kid"])
	})
	if err != nil {
		return nil, InvalidUserID, fmt.Errorf("invalid or expired access token: %v", err)
	}

	// We either have a valid access token or we will attempt to generate new access token.
	userID, err := util.ConvertStringToInt32(claims.Subject)
	if err != nil {
		return nil, InvalidUserID, fmt.Errorf("malformed ID in the token: %v", err)
	}
	user, err := ai.Store.GetUser(ctx, &store.FindUser{
		ID: &userID,
	})
	if err != nil {
		return nil, InvalidUserID, fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return nil, InvalidUserID, fmt.Errorf("user %q not exists", userID)
	}
	if user.RowStatus == store.Archived {
		return nil, InvalidUserID, fmt.Errorf("user %q is archived", userID)
	}

	accessTokens, err := ai.Store.GetUserAccessTokens(ctx, user.ID)
	if err != nil {
		return nil, InvalidUserID, fmt.Errorf("failed to get user access tokens: %v", err)
	}
	if !validateAccessToken(accessToken, accessTokens) {
		return nil, InvalidUserID, fmt.Errorf("invalid access token")
	}

	return token, userID, nil
}

func validateAccessToken(accessTokenString string, userAccessTokens []*store.UserSettingAccessToken) bool {
	for _, userAccessToken := range userAccessTokens {
		if accessTokenString == userAccessToken.AccessToken {
			return true
		}
	}
	return false
}
