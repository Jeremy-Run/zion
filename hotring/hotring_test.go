package hotring

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const Key = "hello%d"
const Value = "world%d"

func Test_Set(t *testing.T) {
	db, _ := InitDictHR(8)
	for i := 0; i <= 10; i++ {
		db.Set(fmt.Sprintf(Key, i), fmt.Sprintf(Value, i))
	}
	ShowTable(db.t0)
}

func Test_Get(t *testing.T) {
	db, _ := InitDictHR(8)

	max := 10
	for i := 0; i <= max; i++ {
		db.Set(fmt.Sprintf(Key, i), fmt.Sprintf(Value, i))
	}
	rand.Seed(time.Now().UnixNano())

	ShowTable(db.t0)

	k := fmt.Sprintf(Key, rand.Intn(max))
	fmt.Printf("======================================\n")
	fmt.Printf("key: %s, value: %s \n", k, db.Get(k))
	fmt.Printf("======================================\n")

	ShowTable(db.t0)
}
