package serialize

type Serializer interface { // 序列化协议只是用来序列化协议体的，不涉及头部
	Code() byte
	Encode(val interface{}) ([]byte, error)
	Decode(data []byte, val interface{}) error
}
