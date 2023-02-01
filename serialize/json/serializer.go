package json

import "encoding/json"

type SerializerJson struct {
}

func (s SerializerJson) Code() byte {
	return 1
}

func (s SerializerJson) Encode(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

func (s SerializerJson) Decode(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}
