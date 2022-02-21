package main

type MetricsHR struct {
	CurrentFactor float64
	RehashTime    int
	LastEntry     *DictEntryHR
	MaxLoad       int
}

type DictHR struct {
	t        []*DictEntryHR
	size     int64
	sizemask int64
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
		t:        make([]*DictEntryHR, 8),
		size:     8,
		sizemask: 7,
		used:     0,
	}
}

func (d *DictHR) Set(key string, val string) {

	h := MurmurHash64A([]byte(key))
	subscript := h & d.sizemask
	if entry := d.t[subscript]; entry != nil {
		pre := entry
		t := 0
		tag := entry.Tag
		for true {
			if entry.Tag == tag {
				if t >= 1 {
					pre.Next = &DictEntryHR{H: h, Tag: pre.Tag + 1, Key: key, Val: val, Next: entry}
					break
				}
				t++
			}
			if entry.Key == key {
				entry.Val = val
				return
			}
			entry = entry.Next
			pre = entry
		}
	} else {
		d.t[subscript] = &DictEntryHR{H: h, Tag: 1, Key: key, Val: val}
		d.t[subscript].Next = d.t[subscript]
	}
	d.used++
}

func (d *DictHR) Get(key string) string {
	h := MurmurHash64A([]byte(key))
	subscript := h & d.sizemask
	if entry := d.t[subscript]; entry != nil {
		t := 0
		tag := entry.Tag
		for true {
			if entry.Tag == tag {
				if t >= 1 {
					return ""
				}
				t++
			}
			if entry.Key == key {
				if d.t[subscript] != entry {
					d.t[subscript] = entry
				}
				return entry.Val
			}

			entry = entry.Next
		}
	}
	return ""
}

func (d *DictHR) AllDB(m *MetricsHR) {
	slotMap := make(map[int64]int)

	maxIndex := int64(0)
	for i := int64(0); i < d.size; i++ {
		if entry := d.t[i]; entry != nil {
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

	e := d.t[maxIndex]
	for true {
		if e.Next == nil {
			m.LastEntry = e
			break
		}
		e = e.Next
	}
}
