package types

import (
	"net/url"
	"strings"
)

const (
	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
	protoUNIX  = "unix"
)

//-------------------------------------------------------------

type parsedURL struct {
	url.URL
	isUnixSocket bool
}

// Parse URL and set defaults
func newParsedURL(remoteAddr string) (*parsedURL, error) {
	u, err := url.Parse(remoteAddr)
	if err != nil {
		return nil, err
	}

	// default to tcp if nothing specified
	if u.Scheme == "" {
		u.Scheme = protoTCP
	}

	pu := &parsedURL{
		URL:          *u,
		isUnixSocket: false,
	}

	if u.Scheme == protoUNIX {
		pu.isUnixSocket = true
	}

	return pu, nil
}

// SetDefaultSchemeHTTP change protocol to HTTP for unknown protocols and TCP protocol - useful for RPC connections
func (u *parsedURL) SetDefaultSchemeHTTP() {
	switch u.Scheme {
	case protoHTTP, protoHTTPS, protoWS, protoWSS:
	default:
		u.Scheme = protoHTTP
	}
}

func (u *parsedURL) GetHostWithPath() string {
	return u.Host + u.EscapedPath()
}

func (u *parsedURL) GetTrimmedHostWithPath() string {
	if !u.isUnixSocket {
		return u.GetHostWithPath()
	}
	return strings.ReplaceAll(u.GetHostWithPath(), "/", ".")
}

// GetDialAddress returns the endpoint to dial for the parsed URL
func (u *parsedURL) GetDialAddress() string {
	if !u.isUnixSocket {
		return u.Host
	}
	return u.GetHostWithPath()
}

func (u *parsedURL) GetTrimmedURL() string {
	return u.Scheme + "://" + u.GetTrimmedHostWithPath()
}
