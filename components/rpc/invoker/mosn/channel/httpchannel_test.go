package channel

import (
	"bufio"
	"context"
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"mosn.io/layotto/components/rpc"
)

type testhttpServer struct {
}

func (ts *testhttpServer) accept(conn net.Conn, listener string) error {
	go ts.readLoop(conn)
	return nil
}

func (ts *testhttpServer) readLoop(conn net.Conn) {
	defer conn.Close()

	for {
		req := fasthttp.AcquireRequest()
		if err := req.Read(bufio.NewReader(conn)); err != nil {
			break
		}

		content := string(req.Body())
		switch content {
		case "close":
			return
		default:
		}

		resp := fasthttp.AcquireResponse()
		resp.SetBody(req.Body())

		if _, err := resp.WriteTo(conn); err != nil {
			break
		}
	}
}

func startTestHttpServer() {
	ts := &testhttpServer{}
	acceptFunc = ts.accept
}

func TestHttpChannel(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello"), Timeout: 1000}
	resp, err := channel.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(resp.Data))
}

func TestRenewHttpConn(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.Error(t, err)

	req = &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello"), Timeout: 1000}
	resp, err := channel.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(resp.Data))
}

func TestConcurrent(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello" + strconv.Itoa(i)), Timeout: 1000}
			resp, err := channel.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, "hello"+strconv.Itoa(i), string(resp.Data))
		}(i)
	}
	wg.Wait()
}
