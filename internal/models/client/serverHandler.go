package client

import (
	"encoding/json"
	"math/rand"
	"net"
	"os"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/netutil"
)

type serverList struct {
	client   *steam.Client
	listPath string
}

func newServerList(client *steam.Client, listPath string) *serverList {
	return &serverList{
		client:   client,
		listPath: listPath,
	}
}

func (s *serverList) HandleEvent(event interface{}) {
	switch e := event.(type) {
	case *steam.ClientCMListEvent:
		d, err := json.Marshal(e.Addresses)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(s.listPath, d, 0666)
		if err != nil {
			panic(err)
		}
	}
}

func (s *serverList) Connect() (bool, error) {
	return s.ConnectBind(nil)
}

func (s *serverList) ConnectBind(laddr *net.TCPAddr) (bool, error) {
	d, err := os.ReadFile(s.listPath)
	if err != nil {
		_, err := s.client.Connect()
		return err == nil, err
	}
	var addrs []*netutil.PortAddr
	err = json.Unmarshal(d, &addrs)
	if err != nil {
		return false, err
	}
	raddr := addrs[rand.Intn(len(addrs))]
	err = s.client.ConnectToBind(raddr, laddr)
	return err == nil, err
}
