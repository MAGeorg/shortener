// пакет, содержащий специфичные ошибки, используемые в проекте.
package errors

import (
	"errors"
)

// ErrAccessDenied ошибка используется при нарушении уникальности в БД.
var ErrAccessDenied = errors.New("unique_violation")

// ErrUnauthrozedID ошибка используется при нахождении в JWT UserID, которому не
// выдавался JWT в cookie.
var ErrUnauthrozedID = errors.New("unauthorized_id")
