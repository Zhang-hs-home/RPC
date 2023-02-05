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
		conn.Close()
		return nil, errors.New("read length data failed")
	}
	headLength := binary.BigEndian.Uint32(msgLenBytes[:4])
	bodyLength := binary.BigEndian.Uint32(msgLenBytes[4:lenBytes])
	bs = make([]byte, headLength+bodyLength)
	n, err := conn.Read(bs[lenBytes:])
	if n != int(headLength+bodyLength-lenBytes) {
		conn.Close()
		return nil, errors.New("tcp连接未读够全部数据")
	}
	copy(bs, msgLenBytes)

	return bs, err
}
