package proxy

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new con")

	if r.Method == http.MethodConnect {
		handleHttps(w, r)
	} else {
		handleHttp(w, r)
	}

}
func DialCustom(network, address string, timeout time.Duration, localIP net.IP) (net.Conn, error) {

	netAddr := &net.TCPAddr{}

	if localIP != nil {
		netAddr.IP = localIP
	}

	fmt.Println("netAddr:", netAddr)

	d := net.Dialer{Timeout: timeout, LocalAddr: netAddr}
	return d.Dial(network, address)
}

func GetRandomIp() net.IP {
	ipList := viper.GetStringSlice("app.ip_list")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ipStr := ipList[r.Intn(len(ipList))]

	ip := net.ParseIP(ipStr)
	return ip
}

func handleHttps(w http.ResponseWriter, r *http.Request) {

	localIP := GetRandomIp()

	destConn, err := DialCustom("tcp", r.Host, time.Second*10, localIP)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)

}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
