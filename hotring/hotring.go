package hotring

import (
	"fmt"
	"github.com/Jeremy-Run/zion/common"
)

type DictHR struct {
	t0, t1   []*DictEntryHR
	size     int64
	sizeMask int64
	used     int64
}

type DictEntryHR struct {
	H    int64
	Tag  int
	Key  string
	Val  string
	Next *DictEntryHR
}

func InitDictHR() *DictHR {
	return &DictHR{
		t0:       make([]*DictEntryHR, 8),
		size:     8,
		sizeMask: 7,
		used:     0,
	}
}

func (d *DictHR) migration() {
	d.t1 = make([]*DictEntryHR, d.size)
	for i := int64(0); i < d.size>>1; i++ {
		entry := d.t0[i]
		if entry == nil {
			continue
		}

		t := 1
		tag := entry.Tag

		for true {
			if entry.Tag == tag {
				if t > 1 {
					break
				}
				t <<= 1
			}

			subscript := entry.H & d.sizeMask
			entry1 := d.t1[subscript]
			if entry1 != nil {
				head := entry1
				pre := entry1
				t1 := 1
				tag1 := entry1.Tag
				for true {
					if entry1.Tag == tag1 {
						if t1 > 1 {
							pre.Next = &DictEntryHR{
								H: entry.H, Tag: pre.Tag + 1, Key: entry.Key, Val: entry.Val, Next: head}
							break
						}
						t1 <<= 1
					}
					pre = entry1
					entry1 = entry1.Next
				}
			} else {
				d.t1[subscript] = &DictEntryHR{H: entry.H, Tag: 1, Key: entry.Key, Val: entry.Val}
				d.t1[subscript].Next = d.t1[subscript]
			}

			entry = entry.Next
		}
	}
}

func (d *DictHR) expandDict() {
	d.size <<= 1
	d.sizeMask = d.size - 1
	d.migration()
	d.t0 = d.t1
	d.t1 = nil
}

func (d *DictHR) Set(key string, val string) {

	if float64(d.used)/float64(d.size) >= 1 {
		d.expandDict()
	}

	h := common.MurmurHash64A([]byte(key))
	subscript := h & d.sizeMask
	if entry := d.t0[subscript]; entry != nil {
		head := entry
		pre := entry
		t := 1
		tag := entry.Tag
		for true {
			if entry.Tag == tag {
				if t > 1 {
					pre.Next = &DictEntryHR{H: h, Tag: pre.Tag + 1, Key: key, Val: val, Next: head}
					break
				}
				t <<= 1
			}
			if entry.Key == key {
				entry.Val = val
				return
			}
			pre = entry
			entry = entry.Next
		}
	} else {
		d.t0[subscript] = &DictEntryHR{H: h, Tag: 1, Key: key, Val: val}
		d.t0[subscript].Next = d.t0[subscript]
	}
	d.used++
}

func (d *DictHR) Get(key string) string {
	h := common.MurmurHash64A([]byte(key))
	subscript := h & d.sizeMask
	if entry := d.t0[subscript]; entry != nil {
		t := 1
		tag := entry.Tag
		for true {
			if entry.Tag == tag {
				if t > 1 {
					return ""
				}
				t <<= 1
			}
			if entry.Key == key {
				if d.t0[subscript] != entry {
					d.t0[subscript] = entry
				}
				return entry.Val
			}

			entry = entry.Next
		}
	}
	return ""
}

func ShowTable(table []*DictEntryHR) {
	for slot, entry := range table {
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
