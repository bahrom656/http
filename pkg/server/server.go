package server

import (
	"strings"
	"bytes"
	"io"
	"log"
	"net"
	"sync"
)

type Handlerfunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]Handlerfunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]Handlerfunc)}
}
func (s *Server) Register(path string, handler Handlerfunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}
func (s *Server) Start() (err error) {

		
		listener, err := net.Listen("tcp", s.addr)
		
		if err != nil {
			return err
		}
	
		defer func() {
			if cerr := listener.Close(); cerr != nil {
				if err == nil {
					err = cerr
					return
				}
				log.Print(cerr)
			}
		}()	
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			err = s.handle(conn)
			if err != nil{
				continue
			}
		}
}


		// for {
		// 	conn, err := listener.Accept()
		// 	if err != nil {
		// 		log.Print(err)
		// 		continue
		// 	}
		// 	go handle(conn)
			
		// }
	



func(s *Server) handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Print(err)
		return err
	}

	if err != nil {
		return err
	}
	data := buf[:n]
	requestlinedelim := []byte{'\r', '\n'}
	requestlineend := bytes.Index(data, requestlinedelim)
	if requestlineend == -1 {
		log.Print("error line end")
		return
	}

	requestline := string(data[:requestlineend])
	parts := strings.Split(requestline, " ")

	if len(parts) != 3 {
		log.Print("error parts")
		return
	}

	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET"{
		log.Print("error get")
		return
	}

	if version != "HTTP/1.1"{
		log.Print(err)
		return
	}
	for pat, handler := range s.handlers{
		if path == pat{
			handler(conn)
		}
	}
	return nil
}

