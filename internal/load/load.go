package load

type Loader interface {
	Start()

	Store(jsonData []byte, total *int, limit int)

	Finish()
}
