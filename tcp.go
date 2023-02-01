package RPC

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// 前4个字节是协议头长度，后4个字节是协议体长度，一共8个字节用于描述协议长度
const lenBytes = 8

func ReadMsg(conn net.Conn) (bs []byte, err error) {
	msgLenBytes := make([]byte, lenBytes)
	length, err := conn.Read(msgLenBytes)
	defer func() {
		if msg := recover(); msg != nil {
			err = errors.New(fmt.Sprintf("%v", msg))
		}
	}()
	if err != nil {
		return nil, err
	}

	if length != lenBytes {
		return nil, errors.New("read length data failed")
	}
	headLength := binary.BigEndian.Uint32(msgLenBytes[:4])
	bodyLength := binary.BigEndian.Uint32(msgLenBytes[4:lenBytes])
	bs = make([]byte, headLength+bodyLength)
	n, err := conn.Read(bs[lenBytes:])
	if n != int(headLength+bodyLength-lenBytes) {
		conn.Close()
		return nil, errors.New("tcp连接未都够全部数据")
		// 你没有读够，你根本不知道该怎么处理，例如你接着等待后续数据，如果后续数据不是本次请求的呢？你读进来了，不仅本次请求的数据不对了，你还
		// 破坏了下次请求的数据
	}
	copy(bs, msgLenBytes)

	return bs, err
}

func EncodeMsg(data []byte) []byte {
	l := len(data)
	resp := make([]byte, l+lenBytes)
	// 先放入长度
	binary.BigEndian.PutUint64(resp, uint64(l))
	// 再放入内容
	copy(resp[lenBytes:], data)
	return resp

}
