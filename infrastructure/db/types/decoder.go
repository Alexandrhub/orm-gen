package types

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

var defaultConfig = jsoniter.ConfigCompatibleWithStandardLibrary

var decoderSingleton *Decode

// Decoder интерфейс для декодирования
type Decoder interface {
	Decode(r io.Reader, val interface{}) error
	Encode(w io.Writer, value interface{}) error
}

// Decode структура для декодирования jsoniter.API
type Decode struct {
	api jsoniter.API
}

// NewDecoder конструктор
func NewDecoder(args ...jsoniter.Config) Decoder {
	conf := defaultConfig
	if len(args) == 0 && decoderSingleton == nil {
		decoderSingleton = &Decode{
			api: conf,
		}
		return decoderSingleton
	}
	if len(args) > 0 {
		conf = args[0].Froze()
	}

	return &Decode{
		api: conf,
	}
}

// Decode метод для декодирования
func (d *Decode) Decode(r io.Reader, val interface{}) error {
	var decoder = d.api.NewDecoder(r)
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}

// Encode метод для кодирования
func (d *Decode) Encode(w io.Writer, value interface{}) error {
	return d.api.NewEncoder(w).Encode(value)
}
