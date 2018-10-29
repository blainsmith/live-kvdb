package db

type KVDB interface {
	Set(key string, value []byte)
	Get(key string) []byte
	Del(key string)
	Keys() []string
}
