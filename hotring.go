package main

import (
	"fmt"
)

const (
	BIG_M = 0xc6a4a7935bd1e995
	BIG_R = 47
	SEED  = 0x1234ABCD
)

func MurmurHash64A(data []byte) (h int64) {
	var k int64
	h = SEED ^ int64(uint64(len(data))*BIG_M)

	var ubigm uint64 = BIG_M
	var ibigm = int64(ubigm)
	for l := len(data); l >= 8; l -= 8 {
		k = int64(data[0]) | int64(data[1])<<8 | int64(data[2])<<16 | int64(data[3])<<24 |
			int64(data[4])<<32 | int64(data[5])<<40 | int64(data[6])<<48 | int64(data[7])<<56

		k := k * ibigm
		k ^= int64(uint64(k) >> BIG_R)
		k = k * ibigm

		h = h ^ k
		h = h * ibigm
		data = data[8:]
	}

	switch len(data) {
	case 7:
		h ^= int64(data[6]) << 48
		fallthrough
	case 6:
		h ^= int64(data[5]) << 40
		fallthrough
	case 5:
		h ^= int64(data[4]) << 32
		fallthrough
	case 4:
		h ^= int64(data[3]) << 24
		fallthrough
	case 3:
		h ^= int64(data[2]) << 16
		fallthrough
	case 2:
		h ^= int64(data[1]) << 8
		fallthrough
	case 1:
		h ^= int64(data[0])
		h *= ibigm
	}

	h ^= int64(uint64(h) >> BIG_R)
	h *= ibigm
	h ^= int64(uint64(h) >> BIG_R)
	return
}

type Dict struct {
	t1, t2     []*DictEntry
	size       int64
	sizemask   int64
	used       int64
	rehashTime int
}

type DictEntry struct {
	h    int64
	key  string
	val  string
	next *DictEntry
}

func InitDict() *Dict {
	return &Dict{
		t1:       make([]*DictEntry, 8),
		size:     8,
		sizemask: 7,
		used:     0,
	}
}

func (d *Dict) migration() {
	d.t2 = make([]*DictEntry, d.size)
	for i := int64(0); i < d.size>>1; i++ {
		if entry := d.t1[i]; entry != nil {
			next := entry.next
			for true {
				if entry2 := d.t2[entry.h&d.sizemask]; entry2 != nil {
					// 新队列里有值
					next2 := entry2.next
					for true {
						if next2 == nil {
							entry2.next = entry
							entry2.next.next = nil

							break
						}
						entry2 = entry2.next
						next2 = entry2.next
					}

				} else {
					// 新队列没值
					d.t2[entry.h&d.sizemask] = entry
					d.t2[entry.h&d.sizemask].next = nil
				}

				if next == nil {
					break
				}

				entry = next
				next = entry.next
			}
		}
	}
}

func (d *Dict) expandDict() {
	d.size = d.size << 1
	d.sizemask = d.size - 1
	d.migration()
	d.t1 = d.t2
	d.t2 = nil
	d.rehashTime++
}

func (d *Dict) Set(key string, val string) {
	if d.used >= d.size {
		d.expandDict()
	}
	var tag int8
	h := MurmurHash64A([]byte(key))
	subscript := h & d.sizemask
	if entry := d.t1[subscript]; entry != nil {
		for true {
			if entry.key == key {
				fmt.Printf("数据已存在 key: %s \n", key)
				break
			}
			if entry.next == nil {
				entry.next = &DictEntry{h: h, key: key, val: val}
				tag = 1
				break
			}
			entry = entry.next
		}
	} else {
		d.t1[subscript] = &DictEntry{h: h, key: key, val: val}
		tag = 1
	}
	if tag == 1 {
		d.used++
	}
}

func (d *Dict) Get(key string) string {
	h := MurmurHash64A([]byte(key))
	if entry := d.t1[h&d.sizemask]; entry != nil {
		for true {
			if entry.key == key {
				return entry.val
			}
			if entry.next == nil {
				return ""
			}
			entry = entry.next
		}
	}
	return ""
}

func (d *Dict) AllDB() {
	slotMap := make(map[int64]int)
	maxLoad := 0
	for i := int64(0); i < d.size; i++ {
		if entry := d.t1[i]; entry != nil {
			for true {
				slotMap[i]++
				if maxLoad < slotMap[i] {
					maxLoad = slotMap[i]
				}

				if entry.next == nil {
					break
				}
				entry = entry.next
			}
		}
	}
	fmt.Printf("slotMap: %v \n", slotMap)
	fmt.Printf("maxLoad: %v \n", maxLoad)
}

func main() {
	d := InitDict()
	for i := 0; i < 100; i++ {
		d.Set(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))
	}
	d.AllDB()
	fmt.Printf("rehashTime: %d \n", d.rehashTime)
}
