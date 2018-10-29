package db

type KVDB interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Del(string) error
	Keys() ([]string, error)
}
