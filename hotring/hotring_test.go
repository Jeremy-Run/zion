package hotring

import (
	"fmt"
	"testing"
)

func Test_Set(t *testing.T) {
	db := InitDictHR()
	db.Set("hello", "world")
	fmt.Printf("%v", db)
}
