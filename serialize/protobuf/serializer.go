package protobuf

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

type SerializerProto struct {
}

func (s SerializerProto) Code() byte {
	return 2
}

func (s SerializerProto) Encode(val interface{}) ([]byte, error) {
	msg, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("必须使用protoc 编译的类型")
	}
	return proto.Marshal(msg)
}

func (s SerializerProto) Decode(data []byte, val interface{}) error {
	msg, ok := val.(proto.Message)
	if !ok {
		return errors.New("必须使用protoc 编译的类型")
	}
	return proto.Unmarshal(data, msg)
}
