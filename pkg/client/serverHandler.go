package client

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	"github.com/Philipp15b/go-steam/v3/netutil"
)

var serverManager = newServerList()

type serverList struct {
	serverList []*net.TCPAddr
	fastest    *net.TCPAddr
}

func newServerList() *serverList {
	s := &serverList{}
	s.init()
	return s
}

func (s *serverList) Connect(client *steam.Client) error {
	if s.fastest == nil {
		_, err := client.Connect()
		return err
	}

	return client.ConnectTo(&netutil.PortAddr{IP: s.fastest.IP, Port: uint16(s.fastest.Port)})
}

func (s *serverList) init() {
	s.serverList = parserServers(fetchCMList())
	s.fastest = getFastestCM(s.serverList)
}

type cmListResponse struct {
	Serverlist struct {
		EndpointList []struct {
			EndpointIP string `json:"endpoint"`
		} `json:"serverlist"`
	} `json:"response"`
}

const cmListFetchURL string = "http://api.steampowered.com/ISteamDirectory/GetCMListForConnect/v1/?cmtype=netfilter"

func fetchCMList() cmListResponse {
	var result cmListResponse
	apiResponse, err := http.Get(cmListFetchURL)
	if err != nil {
		return cmListResponse{}
	}

	apiResponseBody, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return cmListResponse{}
	}

	if err = json.Unmarshal(apiResponseBody, &result); err != nil {
		return cmListResponse{}
	}

	return result
}

func parserServers(response cmListResponse) []*net.TCPAddr {
	addrList := make([]*net.TCPAddr, 0, len(response.Serverlist.EndpointList))

	for _, a := range response.Serverlist.EndpointList {
		if tcpAddr, err := net.ResolveTCPAddr("tcp", a.EndpointIP); err == nil {
			addrList = append(addrList, tcpAddr)
		}
	}

	return addrList
}

func getFastestCM(addrList []*net.TCPAddr) *net.TCPAddr {
	var smallestLatencyAddr *net.TCPAddr
	var smallestLatency int64 = 1000 // random
	var largestLatency int64 = 0
	for _, a := range addrList {
		curTime := time.Now()
		conn, err := net.DialTCP("tcp", nil, a)
		if err != nil {
			if conn != nil {
				conn.Close()
			}
			continue
		}
		latency := time.Since(curTime).Milliseconds()
		conn.Close()
		if latency < smallestLatency {
			smallestLatency = latency
			smallestLatencyAddr = a
		}
		if latency > largestLatency {
			largestLatency = latency
		}
	}
	return smallestLatencyAddr
}
