// пакет для работы с токенами JWT в cookie.
package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	customerr "github.com/MAGeorg/shortener.git/internal/errors"
)

// структура токена JWT.
type claims struct {
	jwt.RegisteredClaims
	UserID int
}

// структура для хранения ID пользователей, которым выдан токен.
// TODO: добавить логгер
type TokensID struct {
	savedID   map[int]struct{}
	lastID    int
	secretKey string
	logger    *zap.SugaredLogger
	TokenEXP  time.Duration
}

// получение нового экземпляра структуры для работы с токенами.
func NewTokensID(k string, t time.Duration, l *zap.SugaredLogger) *TokensID {
	return &TokensID{
		savedID:   make(map[int]struct{}),
		lastID:    0,
		secretKey: k,
		logger:    l,
		TokenEXP:  t,
	}
}

// функция CheckToken проверяет токен, если его нет или он не актуален или его нет,
// то создаем новый.
func (t *TokensID) CheckToken(tokenString string) (string, error) {
	switch {
	// выполнить проверку если она пустая или на валидность куки.
	case tokenString == "" || !t.validToken(tokenString):
		res, err := t.createNewJWTString()
		if err != nil {
			t.logger.Errorln("createNewJWTString: error create jwt: ", err.Error())
			return "", err
		}
		return res, nil
	// проверка, что в куке валидный ID
	case !t.validID(tokenString):
		return "", customerr.ErrUnauthrozedID
	default:
		return tokenString, nil
	}
}

// функция createNewJWTString создает новый токен.
func (t *TokensID) createNewJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.TokenEXP)),
		},
		UserID: t.lastID,
	})

	// сохраняем выданный ID
	t.savedID[t.lastID] = struct{}{}
	// увеличиваем ID, которому последнему выдали JWT.
	t.lastID++
	tokenString, err := token.SignedString([]byte(t.secretKey))
	return tokenString, err
}

// функция parseToken выполняет парсинг токена.
func (t *TokensID) parseToken(tokenString string) (*jwt.Token, *claims, error) {
	c := &claims{}
	key := t.secretKey
	token, err := jwt.ParseWithClaims(tokenString, c,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(key), nil
		})
	return token, c, err
}

// функция validToken выполняет проверку валидности куки.
func (t *TokensID) validToken(tokenString string) bool {
	token, _, err := t.parseToken(tokenString)
	if err != nil {
		t.logger.Errorln("validToken: error parse jwt from cookie: ", err.Error())
	}

	if !token.Valid {
		return false
	}
	return true
}

// функция validID выполняет проверку валидности ID в куке.
func (t *TokensID) validID(tokenString string) bool {
	_, c, err := t.parseToken(tokenString)
	if err != nil {
		t.logger.Errorln("validID: error parse jwt from cookie: ", err.Error())
	}

	if _, ok := t.savedID[c.UserID]; ok {
		return true
	}
	return false
}
