package main

import (
	"fmt"
)

type Config struct {
	DataSize int
	Factor   int64
	ReadTime int
}

type Metrics struct {
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

func InitConfig(factor int64, dataSize int, readTime int) Config {
	return Config{
		DataSize: dataSize,
		Factor:   factor,
		ReadTime: readTime,
	}
}

func InitDict() *Dict {
	return &Dict{
		t1:       make([]*DictEntry, 8),
		size:     8,
		sizemask: 7,
		used:     0,
	}
}

func InitMetrics() *Metrics {
	return &Metrics{
		CurrentFactor: 0,
		RehashTime:    0,
		LastEntry:     nil,
		MaxLoad:       0,
	}
}

func (d *Dict) migration() {
	d.t2 = make([]*DictEntry, d.size)
	for i := int64(0); i < d.size>>1; i++ {
		if entry := d.t1[i]; entry != nil {
			next := entry.Next
			for true {
				if entry2 := d.t2[entry.H&d.sizemask]; entry2 != nil {
					d.t2[entry.H&d.sizemask] = &DictEntry{H: entry.H, Key: entry.Key, Val: entry.Val, Next: entry2}
				} else {
					d.t2[entry.H&d.sizemask] = &DictEntry{H: entry.H, Key: entry.Key, Val: entry.Val}
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

func (d *Dict) expandDict() {
	d.size = d.size << 1
	d.sizemask = d.size - 1
	d.migration()
	d.t1 = d.t2
	d.t2 = nil
}

func (d *Dict) Set(key string, val string, m *Metrics, c Config) {
	m.CurrentFactor = float64(d.used) / float64(d.size)

	if d.used/d.size >= c.Factor {
		d.expandDict()
		m.RehashTime++
	}

	h := MurmurHash64A([]byte(key))
	subscript := h & d.sizemask
	if entry := d.t1[subscript]; entry != nil {
		for true {
			if entry.Key == key {
				fmt.Printf("Data already exists key: %s \n", key)
				return
			}
			if entry.Next == nil {
				entry.Next = &DictEntry{H: h, Key: key, Val: val}
				break
			}
			entry = entry.Next
		}
	} else {
		d.t1[subscript] = &DictEntry{H: h, Key: key, Val: val}
	}
	d.used++
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

func (d *Dict) AllDB(m *Metrics) {
	slotMap := make(map[int64]int)

	maxIndex := int64(0)
	for i := int64(0); i < d.size; i++ {
		if entry := d.t1[i]; entry != nil {
			for true {
				slotMap[i]++
				if m.MaxLoad < slotMap[i] {
					m.MaxLoad = slotMap[i]
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
			m.LastEntry = e
			break
		}
		e = e.Next
	}
}
