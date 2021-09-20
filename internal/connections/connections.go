package connections

type Closer interface {
	Close() error
}
