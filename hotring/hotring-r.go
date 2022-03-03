package hotring

type HT struct {
	Table    *Head
	Size     uint
	SizeMask uint
	R        uint
}

type Head struct {
	Active       byte
	TotalCounter uint
	Address      *HTEntry
}

type HTEntry struct {
	Key      string
	Value    string
	Tag      uint
	Occupied byte
	Rehash   byte
	Counter  uint
	Next     *HTEntry
}

func (e *HTEntry) Get(key string) string {

	return ""
}
