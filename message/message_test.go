package message

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeReq(t *testing.T) {

	req := Request{
		Version:     byte(1),
		Compressor:  byte(2),
		Serializer:  byte(4),
		ServiceName: "",
		MethodName:  "fangfamingcheng",
		Meta: map[string]string{
			"mayDay":    "May",
			"Ja-Morant": "",
			"111":       "222",
		},
		Data: []byte(""),
	}
	req.CalcHeadLength()
	req.BodyLength = uint32(len(req.Data))

	bs := EncodeReq(&req)
	res := DecodeReq(bs)
	assert.Equal(t, req, res)
}

func TestEncodeResp(t *testing.T) {

	req := Response{
		MessageId:  0,
		Version:    byte(1),
		Compressor: byte(2),
		Serializer: byte(4),
		Error:      []byte("32432gerg"),
		Data:       []byte(""),
	}
	req.CalcHeadLength()
	req.BodyLength = uint32(len(req.Data))

	bs := EncodeResp(&req)
	res := DecodeResp(bs)
	assert.Equal(t, req, res)
}

func TestSth(t *testing.T) {
	a := make([]byte, 100)
	b := a[101]
	fmt.Println(b)
}
