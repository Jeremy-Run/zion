package main

import (
	"encoding/json"
	"fmt"
	"time"
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

type Target struct {
	DataSize      int
	Factor        int64
	CurrentFactor float64
	RehashTime    int
	LastEntry     *DictEntry
	MaxLoad       int
}

type Dict struct {
	t1, t2   []*DictEntry
	size     int64
	sizemask int64
	used     int64
}

type DictEntry struct {
	H    int64
	Key  string
	Val  string
	Next *DictEntry
}

func (d *Dict) migration() {
	d.t2 = make([]*DictEntry, d.size)
	for i := int64(0); i < d.size>>1; i++ {
		if entry := d.t1[i]; entry != nil {
			next := entry.Next
			for true {
				if entry2 := d.t2[entry.H&d.sizemask]; entry2 != nil {
					// 新队列里有值
					next2 := entry2.Next
					for true {
						if next2 == nil {
							entry2.Next = entry
							entry2.Next.Next = nil

							break
						}
						entry2 = entry2.Next
						next2 = entry2.Next
					}

				} else {
					// 新队列没值
					d.t2[entry.H&d.sizemask] = entry
					d.t2[entry.H&d.sizemask].Next = nil
				}

				if next == nil {
					break
				}

				entry = next
				next = entry.Next
			}
		}
	}
}

func (d *Dict) expandDict(t *Target) {
	d.size = d.size << 1
	d.sizemask = d.size - 1
	d.migration()
	d.t1 = d.t2
	d.t2 = nil
	t.RehashTime++
}

func (d *Dict) Set(key string, val string, t *Target) {
	t.CurrentFactor = float64(d.used) / float64(d.size)

	if d.used/d.size >= t.Factor {
		d.expandDict(t)
	}

	var tag int8
	h := MurmurHash64A([]byte(key))
	subscript := h & d.sizemask
	if entry := d.t1[subscript]; entry != nil {
		for true {
			if entry.Key == key {
				fmt.Printf("数据已存在 key: %s \n", key)
				break
			}
			if entry.Next == nil {
				entry.Next = &DictEntry{H: h, Key: key, Val: val}
				tag = 1
				break
			}
			entry = entry.Next
		}
	} else {
		d.t1[subscript] = &DictEntry{H: h, Key: key, Val: val}
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
			if entry.Key == key {
				return entry.Val
			}
			if entry.Next == nil {
				return ""
			}
			entry = entry.Next
		}
	}
	return ""
}

func (d *Dict) AllDB(t *Target) {
	slotMap := make(map[int64]int)

	maxIndex := int64(0)
	for i := int64(0); i < d.size; i++ {
		if entry := d.t1[i]; entry != nil {
			for true {
				slotMap[i]++
				if t.MaxLoad < slotMap[i] {
					t.MaxLoad = slotMap[i]
					maxIndex = i
				}

				if entry.Next == nil {
					break
				}
				entry = entry.Next
			}
		}
	}

	e := d.t1[maxIndex]
	for true {
		if e.Next == nil {
			t.LastEntry = e
			break
		}
		e = e.Next
	}

	//fmt.Printf("slotMap: %v \n", slotMap)
}

func InitDict() *Dict {
	return &Dict{
		t1:       make([]*DictEntry, 8),
		size:     8,
		sizemask: 7,
		used:     0,
	}
}

func InitTarget(factor int64, dataSize int) *Target {
	return &Target{
		DataSize:      dataSize,
		Factor:        factor,
		CurrentFactor: 0,
		RehashTime:    0,
		LastEntry:     nil,
		MaxLoad:       0,
	}
}

func main() {
	d := InitDict()
	t := InitTarget(5, 163840)

	for i := 0; i < t.DataSize; i++ {
		d.Set(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i), t)
	}

	d.AllDB(t)

	startTime := time.Now()
	for j := 0; j < 1000000; j++ {
		d.Get(t.LastEntry.Key)
	}
	elapsedTime := time.Since(startTime) / time.Millisecond
	fmt.Printf("Segment finished in %dms \n", elapsedTime)

	tj, _ := json.Marshal(t)
	fmt.Printf("%s \n", tj)
}
