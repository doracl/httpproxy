package main

import (
	// "golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type handler struct{}

func copyHeaders(dst, src http.Header) {
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request url: %s, method: %s", req.RequestURI, req.Method)
	client := &http.Client{}
	host := req.RequestURI

	if req.Method == "CONNECT" {
		host = "https://" + host

		// w.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "server doesn't support hijacking", http.StatusInternalServerError)
			return
		}

		clientConn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var targetSiteCon net.Conn
		// dialer := proxy.FromEnvironment()
		// log.Println(dialer)
		// if dialer != nil {
		// 	targetSiteCon, err = dialer.Dial("tcp", req.RequestURI)
		// } else {
		targetSiteCon, err = net.Dial("tcp", req.RequestURI)
		// }
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Accepting CONNECT to %s", host)
		clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
		go func() {
			if _, err := io.Copy(targetSiteCon, clientConn); err != nil {
				log.Printf("Error copying to client: %s", err)
			}
		}()
		go func() {
			if _, err := io.Copy(clientConn, targetSiteCon); err != nil {
				log.Printf("Error copying to client: %s", err)
			}
		}()
		return
	}
	request, err := http.NewRequest(req.Method, host, nil)

	if err != nil {
		log.Println(err)
	}
	// copyHeaders(request.Header, req.Header)
	for _, k := range []string{"Referer", "Cookie"} {
		request.Header[k] = req.Header[k]
	}
	resp, err := client.Do(request)
	log.Println(resp)
	if err == nil {
		defer resp.Body.Close()
		// copy header
		copyHeaders(w.Header(), resp.Header)

		body, _ := ioutil.ReadAll(resp.Body)
		w.Write(body)
	} else {
		log.Println(err)
	}
}

func main() {
	// log to file
	f, err := os.OpenFile("proxy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	// log.SetOutput(f)

	server := http.Server{
		Addr:    ":8888",
		Handler: &handler{},
	}
	log.Printf("Server is running on: %s", "8888")
	server.ListenAndServe()
}
