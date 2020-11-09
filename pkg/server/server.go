package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

type HandlerFunc func(conn net.Conn)

type RequestFunc func(req *Request)

type Request struct {
	Conn        net.Conn
	QueryParams url.Values
	PathParams  map[string]string
}

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
	requests map[string]RequestFunc
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		handlers: make(map[string]HandlerFunc),
		requests: make(map[string]RequestFunc)}
}
func (s *Server) Register(path string, request RequestFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests[path] = request
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
		if err != nil {
			continue
		}
	}
}

func (s *Server) handle(conn net.Conn) (err error) {
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
	requestLineDelimiter := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelimiter)
	if requestLineEnd == -1 {
		log.Print("error line end")
		return
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		log.Print("error parts")
		return
	}

	method, pathWithQuery, version := parts[0], parts[1], parts[2]
	if method != "GET" {
		log.Print("error get")
	}

	uri, err := url.ParseRequestURI(pathWithQuery)
	if err != nil {
		log.Println("error parse request uri:", err)
	}

	if version != "HTTP/1.1" {
		log.Print(err)
		return
	}
	for pat, handler := range s.requests {
		if uri.Path == pat {
			handler(&Request{
				Conn:        conn,
				QueryParams: uri.Query(),
				PathParams: map[string]string{
					pat: pathWithQuery,
				},
			})
		}
	}
	return nil
}
