package hotring

import (
	"fmt"
	"testing"
)

func showAllDB(db *DictHR) {
	for slot, entry := range db.t {
		if entry == nil {
			fmt.Printf("slot number: %d, link: nil \n", slot)
			continue
		}
		link := make([]*DictEntryHR, 0, 5)
		t := 1
		tag := entry.Tag
		for true {
			if entry.Tag == tag {
				if t > 1 {
					break
				}
				t <<= 1
			}
			link = append(link, entry)
			entry = entry.Next
		}

		var linkStr string
		lenLink := len(link)
		for i, r := range link {
			if i == lenLink-1 {
				linkStr += fmt.Sprintf("{key:%s, value: %s, tag: %d} -> HEAD", r.Key, r.Val, r.Tag)
			} else {
				linkStr += fmt.Sprintf("{key:%s, value: %s, tag: %d} -> ", r.Key, r.Val, r.Tag)
			}
		}
		fmt.Printf("slot number: %d, link: %v \n", slot, linkStr)
	}
}

func Test_Set(t *testing.T) {
	db := InitDictHR()
	for i := 0; i <= 10; i++ {
		db.Set(fmt.Sprintf("hello%d", i), fmt.Sprintf("world%d", i))
	}
	showAllDB(db)
}
