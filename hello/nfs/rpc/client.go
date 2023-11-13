package rpc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"work/golearn/hello/nfs/util"
	"work/golearn/hello/nfs/xdr"
)

const (
	MsgAccepted = iota
	MsgDenied
)

const (
	Success = iota
	ProgUnavail
	ProgMismatch
	ProcUnavail
	GarbageArgs
	SystemErr
)

const (
	RpcMismatch = iota
)

var xid uint32

func init() {
	// seed the XID (which is set by the client)
	xid = rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
}

type Client struct {
	*tcpTransport
}

func DialTCP(network string, ldr *net.TCPAddr, addr string) (*Client, error) {
	a, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP(a.Network(), ldr, a)
	if err != nil {
		return nil, err
	}

	t := &tcpTransport{
		r:  bufio.NewReader(conn),
		wc: conn,
	}

	return &Client{t}, nil
}

type message struct {
	Xid     uint32
	Msgtype uint32
	Body    interface{}
}

func (c *Client) Call(call interface{}) (io.ReadSeeker, error) {
	retries := 1

	msg := &message{
		Xid:  atomic.AddUint32(&xid, 1),
		Body: call,
	}

retry:
	fmt.Printf("%v \n", msg)
	fmt.Printf("%x \n", msg)
	w := new(bytes.Buffer)
	if err := xdr.Write(w, msg); err != nil {
		fmt.Printf("%v \n", err)
		return nil, err
	}

	if _, err := c.Write(w.Bytes()); err != nil {
		fmt.Printf("%v \n", err)
		return nil, err
	}

	res, err := c.recv()
	if err != nil {
		fmt.Printf("%v \n", err)
		return nil, err
	}
	fmt.Printf("%v \n", res)
	fmt.Printf("%x \n", res)

	xid, err := xdr.ReadUint32(res)
	if err != nil {
		fmt.Printf("%v \n", err)
		return nil, err
	}
	fmt.Printf("xid: %x \n", xid)
	fmt.Printf("msg.Xid: %x \n", msg.Xid)
	if xid != msg.Xid {
		return nil, fmt.Errorf("xid did not match, expected: %x, received: %x", msg.Xid, xid)
	}

	mtype, err := xdr.ReadUint32(res)
	if err != nil {
		fmt.Printf("%v \n", err)
		return nil, err
	}

	if mtype != 1 {
		return nil, fmt.Errorf("message as not a reply: %d", mtype)
	}
	fmt.Printf("mtype: %x \n", mtype)

	status, err := xdr.ReadUint32(res)
	if err != nil {
		fmt.Printf("%v \n", err)
		return nil, err
	}
	fmt.Printf("status: %x \n", status)
	switch status {
	case MsgAccepted: //0

		// padding
		_, err = xdr.ReadUint32(res)
		if err != nil {
			panic(err.Error())
		}

		opaque_len, err := xdr.ReadUint32(res)
		if err != nil {
			panic(err.Error())
		}

		_, err = res.Seek(int64(opaque_len), io.SeekCurrent)
		if err != nil {
			panic(err.Error())
		}

		acceptStatus, _ := xdr.ReadUint32(res)

		switch acceptStatus {
		case Success:
			return res, nil
		case ProgUnavail:
			return nil, fmt.Errorf("rpc: PROG_UNAVAIL - server does not recognize the program number")
		case ProgMismatch:
			return nil, fmt.Errorf("rpc: PROG_MISMATCH - program version does not exist on the server")
		case ProcUnavail:
			return nil, fmt.Errorf("rpc: PROC_UNAVAIL - unrecognized procedure number")
		case GarbageArgs:
			// emulate Linux behaviour for GARBAGE_ARGS
			if retries > 0 {
				util.Debugf("Retrying on GARBAGE_ARGS per linux semantics")
				retries--
				goto retry
			}

			return nil, fmt.Errorf("rpc: GARBAGE_ARGS - rpc arguments cannot be XDR decoded")
		case SystemErr:
			return nil, fmt.Errorf("rpc: SYSTEM_ERR - unknown error on server")
		default:
			return nil, fmt.Errorf("rpc: unknown accepted status error: %d", acceptStatus)
		}

	case MsgDenied: // 1
		rejectStatus, _ := xdr.ReadUint32(res)
		fmt.Printf("rejectStatus: %x \n", rejectStatus)
		switch rejectStatus {
		//		case RpcMismatch:

		default:
			return nil, fmt.Errorf("rejected, MsgDenied, Status was not valid: %d", rejectStatus)
		}

	default:
		return nil, fmt.Errorf("rejected,Other, Status was not valid: %d", status)
	}

	//	panic("unreachable")
}