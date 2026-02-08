package storage

// Reader reads items from storage
type Reader struct{}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) Read(path string) ([]byte, error) {
	return nil, nil
}

func (r *Reader) Parse(data []byte) (interface{}, error) {
	return nil, nil
}
