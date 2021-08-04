package config

type Config struct {
	Store
	Cache
}

type Store struct {
	URI string
}

type Cache struct {
	Host   string
	Passwd string
}
