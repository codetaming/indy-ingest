package load

type Loader interface {
	Load(data []byte)
}

type Closer interface {
	Close() error
}
