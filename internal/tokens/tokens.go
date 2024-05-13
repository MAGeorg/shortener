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
//
//nolint:govet // FP
type TokenID struct {
	TokenEXP  time.Duration
	secretKey string
	lastID    int
	logger    *zap.SugaredLogger
	savedID   map[int]struct{}
}

// получение нового экземпляра структуры для работы с токенами.
func NewTokensID(k string, t time.Duration, l *zap.SugaredLogger) *TokenID {
	return &TokenID{
		savedID:   make(map[int]struct{}),
		lastID:    0,
		secretKey: k,
		logger:    l,
		TokenEXP:  t,
	}
}

// функция CheckToken проверяет токен, если его нет или он не актуален или его нет,
// то создаем новый, возвращает токен, UserID из токена, ошибку.
func (t *TokenID) CheckToken(tokenString string) (string, int, error) {
	// выполнить проверку если она пустая или на валидность куки.
	if tokenString == "" {
		res, err := t.createNewJWTString()
		if err != nil {
			t.logger.Errorln("createNewJWTString: error create jwt: ", err.Error())
			return "", t.lastID, err
		}
		return res, t.lastID, nil
	}

	// проверка jwt на наличие ID или актуального ID из уже зарегистрированных
	if userID, res := t.validID(tokenString); res {
		// проверям валидность jwt.
		if !t.validToken(tokenString) {
			res, err := t.createNewJWTString()
			if err != nil {
				t.logger.Errorln("createNewJWTString: error create jwt: ", err.Error())
				return "", t.lastID, err
			}
			return res, t.lastID, err
		}
		return tokenString, userID, nil
	}
	return "", 0, customerr.ErrUnauthrozedID
}

// функция createNewJWTString создает новый токен.
func (t *TokenID) createNewJWTString() (string, error) {
	// увеличиваем ID, которому последнему выдали JWT.
	t.lastID++

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.TokenEXP)),
		},
		UserID: t.lastID,
	})

	// сохраняем выданный ID.
	t.savedID[t.lastID] = struct{}{}

	tokenString, err := token.SignedString([]byte(t.secretKey))
	return tokenString, err
}

// функция parseToken выполняет парсинг токена.
func (t *TokenID) parseToken(tokenString string) (*jwt.Token, *claims, error) {
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
func (t *TokenID) validToken(tokenString string) bool {
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
func (t *TokenID) validID(tokenString string) (int, bool) {
	_, c, err := t.parseToken(tokenString)
	if err != nil {
		t.logger.Errorln("validID: error parse jwt from cookie: ", err.Error())
	}

	if _, ok := t.savedID[c.UserID]; ok {
		return c.UserID, true
	}
	return c.UserID, false
}
