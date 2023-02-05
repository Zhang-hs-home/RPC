package RPC

import (
	"RPC/message"
	"errors"
	"net"
)

const PingPong = 1

func Ping(i interface{}) error {
	conn, ok := i.(net.Conn)
	if !ok {
		return errors.New("req must implements net.Conn")
	}

	req := &message.Request{
		Ping: PingPong,
	}
	req.CalcHeadLength()
	data := message.EncodeReq(req)
	wLen, err := conn.Write(data)
	if err != nil {
		return err
	}

	if wLen != len(data) {
		return errors.New("rpc: 未写入全部数据")
	}

	respMsg, err := ReadMsg(conn)
	if err != nil {
		return err
	}

	res := message.DecodeResp(respMsg)
	if res.Pong != PingPong {
		return errors.New("ping failed")
	}

	return nil
}
