package config

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"sync"
)

type HostChainInfo struct {
	GrpcAddr    string
	Validator   string
	HostAccount string
	Denom       string
	Decimal     uint32
	mu          sync.RWMutex
}

func (hci *HostChainInfo) Set(host string) {
	hci.mu.Lock()
	defer hci.mu.Unlock()

	ip := viper.GetString(fmt.Sprintf("net.ip.%s", host))
	hci.GrpcAddr = ip + ":" + viper.GetString("net.port.grpc")
	hci.Validator = viper.GetString(fmt.Sprintf("%s.val_addr", host))
	hci.HostAccount = viper.GetString(fmt.Sprintf("%s.host_addr", host))
	hci.Denom = viper.GetString(fmt.Sprintf("%s.denom", host))
	hci.Decimal = viper.GetUint32(fmt.Sprintf("%s.decimal", host))
}

type NovaInfo struct {
	Bot     *BotInfo
	ChainID string
	IP      string
	TmRPC   *url.URL
	TmWsRPC *url.URL
	mu      sync.RWMutex
}

func (ni *NovaInfo) Set(addrTarget string, keyname ...string) {
	ni.mu.Lock()
	defer ni.mu.Unlock()

	ni.ChainID = "nova"
	cid := ni.ChainID
	ni.Bot = &BotInfo{}
	if len(keyname) == 1 {
		ni.Bot.mnemonic = Sviper.GetString(keyname[0])
		ni.Bot.passphrase = GetPassphrase(Sviper)
	}
	ni.Bot.Addr = viper.GetString(fmt.Sprintf("%s.%s", cid, addrTarget))
	ni.IP = viper.GetString(fmt.Sprintf("net.ip.%s", cid))
	ni.TmRPC = &url.URL{Scheme: "tcp", Host: ni.IP + ":" + viper.GetString("net.port.tmrpc")}
	ni.TmWsRPC = &url.URL{Scheme: "ws", Host: ni.IP + ":" + viper.GetString("net.port.tmrpc"), Path: "/websocket"}
}

type BotInfo struct {
	mnemonic   string
	Addr       string
	passphrase string
}

func (b BotInfo) Passphrase() string {
	return b.passphrase
}

func (b BotInfo) Mnemonic() string {
	return b.mnemonic
}
