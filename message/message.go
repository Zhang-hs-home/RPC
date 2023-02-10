package message

import (
	"bytes"
	"encoding/binary"
)

const (
	splitter     = '\n'
	pairSplitter = '\r'
)

// ------------------------- 这个就是自定义协议 --------------------------

type Request struct {
	// 头部长度
	HeadLength uint32
	// 消息体长度
	BodyLength uint32
	// 消息id 多路复用使用
	MessageId uint32
	// 版本
	Version byte
	// 压缩算法
	Compressor byte
	// 序列化协议
	Serializer byte
	// ping探活
	Ping byte
	// 服务名
	ServiceName string
	// 方法名
	MethodName string
	// 元数据 可扩展
	Meta map[string]string // 也可以用[]byte，但是处理起来麻烦
	// 消息体
	Data []byte // 不要用interface，interface不知道类型，所以序列化之后是一个map[string]interface类型
}

func (r *Request) CalcHeadLength() {
	// 固定16字节
	res := 16
	res += len(r.ServiceName)
	// 加一个换行符 \n 否则无法区分开serviceName 和methodName
	res += 1
	res += len(r.MethodName)
	res += 1

	for k, v := range r.Meta {
		// 加1是为了区分k和v
		res += len(k) + 1 + len(v) + 1
	}

	r.HeadLength = uint32(res)
}

/*
func EncodeReqV1(r Request) []byte {
	v,_ := json.Marshal(r)
	return v
}
   直接使用现有的序列化方法，不手写了
*/

func EncodeReq(r *Request) []byte {
	res := make([]byte, r.HeadLength+r.BodyLength)
	cur := res

	binary.BigEndian.PutUint32(cur, r.HeadLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur, r.BodyLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur, r.MessageId)
	cur = cur[4:]

	cur[0] = r.Version
	cur = cur[1:]

	cur[0] = r.Compressor
	cur = cur[1:]

	cur[0] = r.Serializer
	cur = cur[1:]

	cur[0] = r.Ping
	cur = cur[1:]

	copy(cur, r.ServiceName)
	cur[len(r.ServiceName)] = splitter
	cur = cur[len(r.ServiceName)+1:]

	copy(cur, r.MethodName)
	cur[len(r.MethodName)] = splitter
	cur = cur[len(r.MethodName)+1:]

	for k, v := range r.Meta {
		copy(cur, k)
		cur[len(k)] = pairSplitter
		cur = cur[len(k)+1:]

		copy(cur, v)
		cur[len(v)] = splitter
		cur = cur[len(v)+1:]
	}

	copy(cur, r.Data)
	return res
}

func DecodeReq(data []byte) *Request {
	req := &Request{}

	req.HeadLength = binary.BigEndian.Uint32(data[:4])
	req.BodyLength = binary.BigEndian.Uint32(data[4:8])
	req.MessageId = binary.BigEndian.Uint32(data[8:12])
	req.Version = data[12]
	req.Compressor = data[13]
	req.Serializer = data[14]
	req.Ping = data[15]

	// 取头部剩下所有的数据
	remainHeader := data[16:req.HeadLength]

	// 取出服务名
	split := bytes.IndexByte(remainHeader, splitter)
	req.ServiceName = string(remainHeader[:split])
	remainHeader = remainHeader[split+1:]

	// 取出方法名
	split = bytes.IndexByte(remainHeader, splitter)
	req.MethodName = string(remainHeader[:split])
	remainHeader = remainHeader[split+1:]

	// 取出元数据
	if len(remainHeader) > 0 {
		metaMap := make(map[string]string)
		split = bytes.IndexByte(remainHeader, splitter)
		for split != -1 {
			pair := remainHeader[:split]
			seg := bytes.IndexByte(remainHeader, pairSplitter)
			key := string(pair[:seg])
			value := string(pair[seg+1:])
			metaMap[key] = value
			remainHeader = remainHeader[split+1:]
			split = bytes.IndexByte(remainHeader, splitter)
		}
		req.Meta = metaMap
	}

	// 剩下的就是协议体了
	req.Data = data[req.HeadLength:]
	return req

}

type Response struct {
	// 头部长度
	HeadLength uint32
	// 消息体长度
	BodyLength uint32
	// 消息id 多路复用使用
	MessageId uint32
	// 版本
	Version byte
	// 压缩算法
	Compressor byte
	// 序列化协议
	Serializer byte
	// pong探活
	Pong byte
	// 错误信息 可以是业务error，也可以是框架error
	Error []byte
	// 协议体
	Data []byte
}

func (r *Response) CalcHeadLength() {
	res := 16
	res += len(r.Error)
	r.HeadLength = uint32(res)
}

func EncodeResp(r *Response) []byte {
	res := make([]byte, r.HeadLength+r.BodyLength)
	cur := res

	binary.BigEndian.PutUint32(cur, r.HeadLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur, r.BodyLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur, r.MessageId)
	cur = cur[4:]

	cur[0] = r.Version
	cur = cur[1:]

	cur[0] = r.Compressor
	cur = cur[1:]

	cur[0] = r.Serializer
	cur = cur[1:]

	cur[0] = r.Pong
	cur = cur[1:]

	copy(cur, r.Error)
	cur = cur[len(r.Error):]

	copy(cur, r.Data)
	return res
}

func DecodeResp(data []byte) *Response {
	res := &Response{}

	res.HeadLength = binary.BigEndian.Uint32(data[:4])
	res.BodyLength = binary.BigEndian.Uint32(data[4:8])
	res.MessageId = binary.BigEndian.Uint32(data[8:12])
	res.Version = data[12]
	res.Compressor = data[13]
	res.Serializer = data[14]
	res.Pong = data[15]

	// 取出error
	remainErr := data[16:res.HeadLength]
	res.Error = remainErr

	// 剩下的就是协议体了
	res.Data = data[res.HeadLength:]
	return res

}
