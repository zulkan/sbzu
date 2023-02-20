package domain

//go:generate mockery --name Cache
type Cache interface {
	Write(key, value string) error
	Read(key string) (string, error)
}
