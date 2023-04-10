package cache

type ICache interface {
	Get(key string) (string, error)
	Set(key string, val string) error
	Delete(key string) error
	Has(key string) bool
}
