package account

import (
	"fmt"
)

var newID float64 = 0

// Как альтернативу можно было бы
// использовать github.com/google/uuid
func GetNewID() string {
	newID += 1
	return fmt.Sprint(newID)
}
