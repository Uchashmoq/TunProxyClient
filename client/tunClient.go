package client

import (
	"log"
	"net"
)

type TunClient struct {
	ServerAddr string
	ClientAddr string
	listener   net.Listener
	iv         []byte
}

func NewTunClient(serverAddr, clientAddr, ivstr string) *TunClient {
	return &TunClient{
		ServerAddr: serverAddr,
		ClientAddr: clientAddr,
		listener:   nil,
		iv:         []byte(ivstr),
	}
}
func (t *TunClient) Launch() {
	listen, err := net.Listen("tcp", t.ClientAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Listening %s\n", t.ClientAddr)
	t.listener = listen
}
func (t *TunClient) Accepting() {
	for {
		accept, err := t.listener.Accept()
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%s connected\n", accept.RemoteAddr().String())
			tun := NewTun(accept, t.ServerAddr, t.iv)
			err := tun.Connect()
			if err != nil {
				log.Println(err)
				accept.Close()
				continue
			}
			tun.StartProxy()
		}
	}
}
