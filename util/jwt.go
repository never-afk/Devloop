package util

import (
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

var jwtOptions = make(map[string]*JwtOption)

var jwtDefaultAuthName = `default`

var jwtRWLock = sync.RWMutex{}

//JwtOption
type JwtOption struct {
	SecretKey []byte
	Timeout   int64
	AuthName  string
}

type Claims struct {
	jwt.StandardClaims
	Id      string `json:"id"`
	Payload interface{}
}

func SetJwtOption(opt *JwtOption) {
	if opt.AuthName == `` {
		opt.AuthName = jwtDefaultAuthName
	}
	jwtRWLock.Lock()
	defer jwtRWLock.Unlock()
	jwtOptions[opt.AuthName] = opt
}

func GetJwtOption(AuthName string) *JwtOption {
	jwtRWLock.RLock()
	defer jwtRWLock.RUnlock()
	return jwtOptions[AuthName]
}

func getOpt(authName ...string) *JwtOption {
	key := jwtDefaultAuthName
	if len(authName) != 0 {
		key = authName[0]
	}
	return GetJwtOption(key)
}

// generate token
// pass auth name to use different option
func GenerateToken(id string, authName ...string) (string, error) {
	opt := getOpt(authName...)
	claims := Claims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(opt.Timeout)).Unix(),
		},
		id,
		nil,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(opt.SecretKey)
	return token, err
}

// generate token
// payload is an extension info
// pass auth name to use different option
func GenerateTokenPayload(id string, payload interface{}, authName ...string) (string, error) {
	opt := getOpt(authName...)
	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(opt.Timeout)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
		payload,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(opt.SecretKey)
	return token, err
}

// parse token string
func ParseToken(token string, authName ...string) (*Claims, error) {
	opt := getOpt(authName...)

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return opt.SecretKey, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
