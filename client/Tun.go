package client

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
	"tunProxy/utils"
)

type Tun struct {
	LocalConn  net.Conn
	RemoteConn net.Conn
	ServerAddr string
	key        []byte
	iv         []byte
	flag       bool
	errCh      chan error
	startTime  time.Time
	bytes      int64
}

func NewTun(localConn net.Conn, serverAddr string, iv []byte) *Tun {
	return &Tun{
		LocalConn:  localConn,
		ServerAddr: serverAddr,
		flag:       false,
		iv:         iv,
		errCh:      make(chan error),
		bytes:      0,
	}
}
func (t *Tun) Connect() error {
	dial, err := net.Dial("tcp", t.ServerAddr)
	if err != nil {
		return err
	}
	t.RemoteConn = dial
	ch := make(chan []byte, 1)
	var err1 error
	go func() {
		b := make([]byte, 32)
		_, err1 = dial.Read(b)
		ch <- b
	}()
	select {
	case <-time.After(100 * time.Second):
		return errors.New("connect time out")
	case k := <-ch:
		t.key = k
		return err1
	}
}
func (t *Tun) StartProxy() {
	t.flag = true
	t.startTime = time.Now()
	go t.recvfromBrowserEncodeAndSend()
	go t.recvfromServerDecodeAndSendToBrowser()
	go func() {
		_ = <-t.errCh
		t.flag = false
		t.ShutDown()
		endTime := time.Now()
		log.Printf("tunnel closed ,communication time : %.2f s, %d bytes transmitted ", endTime.Sub(t.startTime).Seconds(), t.bytes)
	}()
}
func (t *Tun) ShutDown() {
	t.LocalConn.Close()
	t.RemoteConn.Close()
}

func (t *Tun) recvfromBrowserEncodeAndSend() {
	for t.flag {
		buf := make([]byte, 1024*1024*5)
		n, err := t.LocalConn.Read(buf)
		if err != nil {
			t.errCh <- err
			return
		}
		message := EncodeMessage(buf[:n], t.key, t.iv)
		_, err1 := t.RemoteConn.Write(message)
		log.Printf("%d bytes >>> server ", len(message))
		t.bytes += int64(len(message))
		if err1 != nil {
			t.errCh <- err
			return
		}
	}
}
func (t *Tun) recvfromServerDecodeAndSendToBrowser() {
	frameDecoder := utils.NewFrameDecoder(2, binary.BigEndian.Uint32)
	go frameDecoder.Separate()
	go func() {
		for t.flag {
			buf := make([]byte, 1024*1024*5)
			n, err := t.RemoteConn.Read(buf)
			frameDecoder.In <- buf[:n]
			log.Printf("%d bytes <<< server", n)
			t.bytes += int64(n)
			if err != nil {
				log.Println(err)
				t.errCh <- err
				return
			}
		}
	}()
	go func() {
		for t.flag {
			bytes := <-frameDecoder.Out
			message := DecodeMessage(bytes, t.key, t.iv)
			_, err1 := t.LocalConn.Write(message)
			if err1 != nil {
				t.errCh <- err1
				return
			}
		}
	}()
}
