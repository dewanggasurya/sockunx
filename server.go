package sockunx

import (
	"fmt"
	"log"
	"net"
	"os"
	"sockunx/pkg/helper"
	"strings"
	"sync"
	"syscall"
)

// Handler struct
type Handler func(request string) (response interface{}, e error)

// Server struct
type Server struct {
	size      int
	path      string
	listener  net.Listener
	handler   Handler
	isStopped bool
	isRunning bool
	mutex     sync.Mutex
}

// NewServer socket
func NewServer(path string, size ...int) (*Server, error) {
	s := 2048
	if len(size) > 0 && size[0] > 0 {
		s = size[0]
	}

	syscall.Unlink(path)

	return &Server{
		size:  s,
		path:  path,
		mutex: sync.Mutex{},
		handler: func(request string) (response interface{}, e error) {
			return request, nil
		},
	}, nil
}

// RegisterHandler for a single handler
func (o *Server) RegisterHandler(handler Handler) {
	o.handler = handler
}

// Run server
func (o *Server) Run(await ...bool) error {
	needAwait := false
	if len(await) > 0 {
		needAwait = await[0]
	}

	if needAwait {
		hasRun := make(chan bool, 1)
		go o.run(hasRun)
		<-hasRun
		return nil
	}

	return o.run()
}

func (o *Server) run(hasRun ...chan bool) error {
	_hasRun := false
	hasRunChan := make(chan bool, 1)
	if len(hasRun) > 0 && hasRun[0] != nil {
		hasRunChan = hasRun[0]
	}

	listener, e := net.Listen("unix", o.path)
	if e != nil {
		if strings.Contains(e.Error(), "no such file or directory") {
			_, e := helper.EnsureFile(o.path)
			if e != nil {
				return fmt.Errorf("unable to create file: %s", e.Error())
			}
			close(hasRunChan)
			return o.run(hasRun...)
		}
		return fmt.Errorf("unable to listen %s: %s", o.path, e.Error())
	}

	e = os.Chmod(o.path, os.ModePerm)
	if e != nil {
		log.Printf("WARN unable to change permission %s", e.Error())
	}

	o.mutex.Lock()
	o.listener = listener
	o.isStopped = false
	o.mutex.Unlock()

	for {
		if o.isStopped {
			break
		}

		if !_hasRun {
			_hasRun = true
			hasRunChan <- _hasRun
		}

		o.mutex.Lock()
		o.isRunning = true
		o.mutex.Unlock()

		conn, e := o.listener.Accept()
		if e != nil {
			return e
		}

		if conn == nil {
			return fmt.Errorf("unable to create connection")
		}

		done := make(chan bool, 1)
		go func(conn net.Conn, size int) {
			defer func() {
				conn.Close()
				done <- true
			}()

			buf := make([]byte, size)
			nr, e := conn.Read(buf)
			if e != nil {
				log.Printf("ERR unable to read buffer: %s", e.Error())
				return
			}

			requestStr := blank
			requestBytes := buf[0:nr]
			if len(requestBytes) > size {
				log.Printf("WARN unable to parse request as socket.Message: %s", e.Error())
			} else {
				var request Message
				e = helper.FromJSON(requestBytes, &request)

				if e != nil {
					log.Printf("WARN unable to parse request as socket.Message: %s", e.Error())
				} else {
					requestStr = string(request.Data)
				}

			}

			if requestStr == blank {
				requestStr = ""
			}

			response, e := o.handler(requestStr)
			if e != nil {
				response = fmt.Errorf("ERR %s", e.Error())
			}

			if response == "" {
				response = blank
			}

			_, e = conn.Write([]byte(fmt.Sprintf("%v", response)))
			if e != nil {
				log.Printf("ERR Writing client error: %s", e)
			}

		}(conn, o.size)
		<-done
	}
	return nil
}

// IsRunning server
func (o *Server) IsRunning() bool {
	return o.isRunning
}

// Stop server
func (o *Server) Stop() error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.isRunning = false
	o.isStopped = true
	if o.listener != nil {
		return o.listener.Close()
	}

	os.Remove(o.path)

	return nil
}
