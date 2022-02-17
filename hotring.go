package main

import (
	"fmt"
	"unsafe"
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

const (
	c1_32 uint32 = 0xcc9e2d51
	c2_32 uint32 = 0x1b873593
)

// GetHash returns a murmur32 hash for the data slice.
func GetHash(data []byte) uint32 {
	// Seed is set to 37, same as C# version of emitter
	var h1 uint32 = 37

	nblocks := len(data) / 4
	var p uintptr
	if len(data) > 0 {
		p = uintptr(unsafe.Pointer(&data[0]))
	}

	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19) // rotl32(h1, 13)
		h1 = h1*5 + 0xe6546b64
	}

	tail := data[nblocks*4:]

	var k1 uint32
	switch len(tail) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return (h1 << 24) | (((h1 >> 8) << 16) & 0xFF0000) | (((h1 >> 16) << 8) & 0xFF00) | (h1 >> 24)
}

type Dict struct {
	t1, t2     []*DictEntry
	size       uint32
	sizemask   uint32
	used       uint32
	rehashTime int
}

type DictEntry struct {
	h    uint32
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
	for i := uint32(0); i < d.size>>1; i++ {
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

func set(key string, val string, d *Dict) {

	h := GetHash([]byte(key))
	subscript := h & d.sizemask
	if entry := d.t1[subscript]; entry != nil {
		for true {
			if entry.key == key {
				fmt.Printf("数据已存在 key: %s \n", key)
				break
			}
			if entry.next == nil {
				entry.next = &DictEntry{key: key, val: val}
				break
			}
			entry = entry.next
		}
	} else {
		d.t1[subscript] = &DictEntry{key: key, val: val}
	}
}

func (d *Dict) Set(key string, val string) {
	if d.used >= d.size {
		d.expandDict()
	}
	set(key, val, d)
	d.used++
}

func (d *Dict) Get(key string) string {
	h := GetHash([]byte(key))
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
	slotMap := make(map[uint32]int)
	maxLoad := 0
	for i := uint32(0); i < d.size; i++ {
		if entry := d.t1[i]; entry != nil {
			for true {
				slotMap[i]++
				if maxLoad < slotMap[i] {
					maxLoad = slotMap[i]
				}

				//fmt.Printf("slot: %d, entry.key: %s, entry.val: %s \n", i, entry.key, entry.val)

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
	for i := 0; i < 10000; i++ {
		d.Set(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))
	}
	d.AllDB()
	fmt.Printf("rehashTime: %d \n", d.rehashTime)
}
