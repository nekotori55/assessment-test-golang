package account

import (
	"fmt"
)

var newID float64 = 0

func GetNewID() string {
	newID += 1
	return fmt.Sprint(newID)
}
