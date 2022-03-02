package config

type Config struct {
	DataSize int
	Factor   int64
	ReadTime int
}

func InitConfig(factor int64, dataSize int, readTime int) Config {
	return Config{
		DataSize: dataSize,
		Factor:   factor,
		ReadTime: readTime,
	}
}
