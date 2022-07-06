package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func NewLBTransport() *LBTransport {
	return &LBTransport{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

type LBTransport struct {
	Transport http.RoundTripper
}

func (t *LBTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ips, err := net.LookupIP(req.URL.Host)
	if err != nil {
		return nil, err
	}
	req.URL.Host = ips[rand.Intn(len(ips))].String()
	return t.Transport.RoundTrip(req)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	client := http.Client{
		Transport: NewLBTransport(),
	}
	resp, err := client.Get("https://www.baidu.com")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))
	resp.Body.Close()

	resp, err = client.Get("https://www.baidu.com")
	if err != nil {
		panic(err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	resp, err = client.Get("https://www.baidu.com")
	if err != nil {
		panic(err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}
