package storage

// Writer writes items to storage
type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(path string, data []byte) error {
	return nil
}
