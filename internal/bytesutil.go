package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func uint64tob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btouint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func serialize(in interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserialize(in []byte, out interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(in)).Decode(out)
}
