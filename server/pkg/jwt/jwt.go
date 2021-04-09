package jwt

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/oauth"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

const Issuer string = "scrimpton"

type Claims struct {
	jwt.StandardClaims
	AuthorID string          `json:"author_id"`
	Approver bool            `json:"approver"`
	Identity *oauth.Identity `json:"identity"`
}

func (c *Claims) FromMap(claims jwt.MapClaims) {
	c.StandardClaims.Issuer, _ = claims["iss"].(string)
	c.StandardClaims.ExpiresAt, _ = claims["exp"].(int64)

	identityMap, ok := claims["identity"].(map[string]interface{})
	if !ok {
		return
	}
	c.Identity = &oauth.Identity{
		ID:               identityMap["id"].(string),
		Name:             identityMap["name"].(string),
		HasVerifiedEmail: identityMap["has_verified_email"].(bool),
		Icon:             identityMap["icon_img"].(string),
		IsSuspended:      identityMap["is_suspended"].(bool),
		Created:          identityMap["created"].(float64),
		CreatedUTC:       identityMap["created_utc"].(float64),
		TotalKarma:       int64(identityMap["total_karma"].(float64)),
		Over18:           identityMap["over_18"].(bool),
	}
}

type Config struct {
	SigningKey string
	ExpireTime int64
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.SigningKey, prefix, "jwt-key", "insecure", "Key used to sign JWTs")
	flag.Int64VarEnv(fs, &c.ExpireTime, prefix, "jwt-expire-time", 60*60*24*365, "Number of seconds token is valid for")
}

func NewAuth(cfg *Config) *Auth {
	return &Auth{cfg: cfg}
}

type Auth struct {
	cfg *Config
}

func (a *Auth) NewJWTForIdentity(author *models.Author, ident *oauth.Identity) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + a.cfg.ExpireTime,
			Issuer:    Issuer,
		},
		AuthorID: author.ID,
		Approver: author.Approver,
		Identity: ident,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.cfg.SigningKey))
}

func (a *Auth) VerifyToken(tokenString string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.SigningKey), nil
	})
	if token != nil && token.Valid {

		//claimMap, ok := token.Claims.(jwt.MapClaims)
		//if !ok {
		//	return nil, fmt.Errorf("token claims were malformed")
		//}
		//claims := &Claims{}
		//claims.FromMap(claimMap)
		return claims, nil
	}
	if err == nil {
		return nil, fmt.Errorf("no error occured, but token was not generated succesfully")
	}
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, fmt.Errorf("token was malformed")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, fmt.Errorf("token is either expired or not yet active")
		} else {
			return nil, fmt.Errorf("failed to verify token: %s", err)
		}
	}
	return nil, fmt.Errorf("failed to verify token: %s", err)
}

func ExtractTokenFromRequestContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	auth := md.Get("authorization")
	if len(auth) < 1 {
		return ""
	}

	token := strings.TrimSpace(strings.TrimPrefix(auth[0], "Bearer"))

	return token
}
