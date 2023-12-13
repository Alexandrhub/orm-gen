package decoder

import (
	"io"
)

// Decoder интерфейс для декодирования
type Decoder interface {
	Decode(r io.Reader, val interface{}) error
	Encode(w io.Writer, value interface{}) error
}
