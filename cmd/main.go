package main

import (
	"github.com/bahrom656/http/pkg/server"
	"log"
	"net"
	"os"
)
//
//const crlf = "\r\n"

func main() {
	host := "0.0.0.0"
	port := "443"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}

}
func execute(host, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/payments/{id}", func(req *server.Request) {
		
		v := "1"
		req.PathParams["id"] = v
		id := req.PathParams["id"]
		log.Print(id)


	})
	return srv.Start()
}
