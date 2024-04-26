// пакет, содержащий специфичные ошибки, используемые в проекте.
package errors

import (
	"errors"
)

// ErrAccessDenied ошибка, используемая при нарушении уникальности в БД.
var ErrAccessDenied = errors.New("unique_violation")
