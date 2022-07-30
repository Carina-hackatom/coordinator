package rpc

import (
	tmrpc "github.com/Carina-hackatom/coordinator/rpc/types"
	"github.com/Carina-hackatom/coordinator/utils"
	"net/url"
)

func MakeEventWatcher(remoteAddr url.URL) *tmrpc.WSClient {
	wsc, err := tmrpc.NewWS("//"+remoteAddr.Host, remoteAddr.Path)
	utils.CheckErr(err, "", 0)
	//wsc.SetLogger(tmlog.NewTMJSONLogger(tmlog.NewSyncWriter(os.Stdout)))
	return wsc
}
