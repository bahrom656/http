// package main

// import (
// 	"strconv"
// 	"strings"
// 	"bytes"
// 	"io"
// 	"log"
// 	"net"
// 	"os"
// )
// const crlf = "\r\n"
// func main() {
// 	host := "127.0.0.1"
// 	port := "443"

// 	if err := execute(host, port); err != nil {
// 		os.Exit(1)
// 	}
// }

// func execute(host string, port string) (err error) {
// 	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
// 	if err != nil {
// 		log.Print(err)
// 		return err
// 	}

// 	defer func() {
// 		if cerr := listener.Close(); cerr != nil {
// 			if err == nil {
// 				err = cerr
// 				return
// 			}
// 			log.Print(cerr)
// 		}
// 	}()

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Print(err)
// 			continue
// 		}

// 		err = handle(conn)
// 		if err != nil {
// 			log.Print(err)
// 			continue
// 		}
// 	}
// }

// func handle(conn net.Conn) (err error) {
// 	defer func() {
// 		if cerr := conn.Close(); cerr != nil {
// 			if err == nil {
// 				err = cerr
// 				return
// 			}
// 			log.Print(cerr)
// 		}
// 	}()


// 	buf := make([]byte, 4096)
// 	n, err := conn.Read(buf)
// 	if err == io.EOF{
// 		log.Print(err)
// 		return err
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	data := buf[:n]
// 	requestlinedelim := []byte{'\r', '\n'}
// 	requestlineend := bytes.Index(data, requestlinedelim)
// 	if requestlineend == -1 {
// 		return nil
// 	}

// 	requestline := string(data[:requestlineend])
// 	parts := strings.Split(requestline, " ")
// 	if len(parts) != 3 {
// 		log.Print("error parts")
// 		return
// 	}

// 	method, path, version := parts[0], parts[1], parts[2]

// 	if method != "GET"{
// 		log.Print("error GET")
// 		return
// 	}

// 	if version != "HTTP/1.1" {
// 		log.Print("error version")
// 		return
// 	}

// 	if path == "/" {
// 		body := "ok!"
// 		_, err = conn.Write([]byte(
// 						"HTTP/1.1 200 OK" + crlf +
// 							"Content-Length: " + strconv.Itoa(len(body)) + crlf +
// 							"Content-Type: text/html" + crlf +
// 							"Connection: close" + crlf +
// 							crlf + body,
// 						))

// 						if err != nil {
// 							log.Print(err)
// 							return err
// 						}
// 	}
// 	return nil
// }
package main

import (
	"log"
	"net"
	"os"
	"github.com/bahrom656/http/pkg/server"
	"strconv"
)

const crlf = "\r\n"

func main() {
	host := "0.0.0.0"
	port := "443"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}

}
func execute(host, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", func(conn net.Conn) {
		body := "Welcome to our web site"

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK" + crlf +
				"Content-Length: " + strconv.Itoa(len(body)) + crlf +
				"Content-Type: text/html" + crlf +
				"Connection: close" + crlf +
				crlf + body,
			))
		if err != nil {
			log.Print(err)
		}
	})
	srv.Register("/about", func(conn net.Conn) {
		body := "About Golang Academy"

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK" + crlf +
				"Content-Length: " + strconv.Itoa(len(body)) + crlf +
				"Content-Type: text/html" + crlf +
				"Connection: close" + crlf +
				crlf + body,
		))
		if err != nil {
			log.Print(err)
		}
	})
	return srv.Start()
}
