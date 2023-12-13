package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/decoder"
	
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
)

const (
	TypePost = iota + 1
	TypeAvatar
	TypeUserModerate
	TypeCommentReply
	TypeChatMessages
	TypeSubscribe
)

// NullString структура обёртки вокруг sql.NullString
type NullString struct {
	sql.NullString
}

// MarshalJSON метод, вызываемый json.Marshal,
// для экземпляр типа NullString
func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", x.String)), nil
}

// UnmarshalJSON метод, вызываемый json.Unmarshal
// для экземпляр типа NullString
func (x *NullString) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.String)
	if err != nil {
		return err
	}
	x.Valid = true
	if len(x.String) == 0 {
		x.Valid = false
	}
	
	return nil
}

// NullBool структура обёртки вокруг sql.NullBool
type NullBool struct {
	sql.NullBool
}

// MarshalJSON метод, вызываемый json.Marshal,
// для экземпляр типа NullBool
func (x *NullBool) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Bool)
}

// UnmarshalJSON метод, вызываемый json.Unmarshal
// для экземпляр типа NullBool
func (x *NullBool) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.Bool)
	if err != nil {
		return err
	}
	x.Valid = true
	
	return nil
}

// NullInt64 структура обёртки вокруг sql.NullInt64
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON метод, вызываемый json.Marshal,
// для экземпляр типа NullInt64
func (x *NullInt64) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Int64)
}

// UnmarshalJSON метод, вызываемый json.Unmarshal
// для экземпляр типа NullInt64
func (x *NullInt64) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.Int64)
	if err != nil {
		return err
	}
	x.Valid = true
	
	return nil
}

// NullFloat64 структура обёртки вокруг sql.NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON метод, вызываемый json.Marshal,
// для экземпляр типа NullFloat64
func (x *NullFloat64) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Float64)
}

// UnmarshalJSON метод, вызываемый json.Unmarshal
// для экземпляр типа NullFloat64
func (x *NullFloat64) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.Float64)
	if err != nil {
		return err
	}
	x.Valid = true
	
	return nil
}

// NullUUID структура обёртки вокруг null.UUID
type NullUUID struct {
	Binary []byte
	Valid  bool
	String string
}

// MarshalJSON метод, вызываемый json.Marshal,
// для экземпляр типа NullUUID
func (x *NullUUID) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", x.String)), nil
}

// UnmarshalJSON метод, вызываемый json.Unmarshal
// для экземпляр типа NullUUID
func (x *NullUUID) UnmarshalJSON(data []byte) error {
	if data[0] == 110 {
		*x = NullUUID{}
		return nil
	}
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.String)
	if err != nil {
		return err
	}
	x.Valid = true
	uuidRaw, err := uuid.Parse(x.String)
	if err != nil {
		return err
	}
	x.Binary, err = uuidRaw.MarshalBinary()
	if err != nil {
		return err
	}
	
	return nil
}

// Scan имплементирует sql.Scanner
func (x *NullUUID) Scan(value interface{}) error {
	if value == nil {
		*x = NullUUID{}
		return nil
	}
	
	var dest []byte
	switch source := value.(type) {
	case string:
		*x = NewNullUUID(source)
		return nil
	case []byte:
		if len(source) == 0 {
			dest = nil
		} else {
			dest = make([]byte, len(source))
			copy(dest, source)
		}
	case nil:
		*x = NullUUID{}
	default:
		return errors.New("incompatible type for NullUUID")
	}
	
	uuidRaw, err := uuid.FromBytes(dest)
	if err != nil {
		return err
	}
	x.Binary = dest
	x.Valid = true
	
	x.String = uuidRaw.String()
	
	return nil
}

// Value имплементирует driver.Valuer.
func (x NullUUID) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	b := make([]byte, len(x.Binary))
	copy(b, x.Binary)
	
	return b, nil
}

// NullUint64 структура обёртки вокруг sql.NullUint64
type NullUint64 struct {
	null.Uint64
}

// NullTime структура обёртки вокруг null.Time
type NullTime struct {
	null.Time
}

type User struct{}

type Chat struct{}
