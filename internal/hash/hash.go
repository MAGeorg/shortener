// пакет для работы с хэшем.
package hash

import (
	"hash/fnv"
)

// GetHash - функция для генерации хэша uint32 от строки.
func GetHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
