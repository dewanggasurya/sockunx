package sockunx

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	blank = "<!blank>"
)

// Client struct
type Client struct {
	size     int
	path     string
	handler  Handler
	isClosed bool
	mutex    sync.Mutex
}

// NewClient socket
func NewClient(path string, size ...int) (*Client, error) {
	s := 4096
	if len(size) > 0 && size[0] > 0 {
		s = size[0]
	}

	return &Client{
		size:  s,
		path:  path,
		mutex: sync.Mutex{},
	}, nil
}

// Send message
func (o *Client) Send(message string, size ...int) (string, error) {
	s := o.size
	if len(size) > 0 && size[0] > 0 {
		s = size[0]
	}

	conn, e := net.Dial("unix", o.path)
	if e != nil {
		return "", fmt.Errorf("unable to dial %s : %s", o.path, e.Error())
	}
	// conn.SetDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()

	response := make(chan []byte, 1)
	go func(r io.Reader) {
		buf := make([]byte, s)
		for {
			n, e := r.Read(buf[:])
			if e != nil {
				response <- []byte(blank)
				return
			}
			response <- buf[0:n]
		}
	}(conn)

	if message == "" {
		message = blank
	}

	messageBytes := []byte(message)
	_, e = conn.Write(Message{
		Data:   messageBytes,
		Length: len(messageBytes),
	}.Bytes())
	if e != nil {
		return "", fmt.Errorf("unable to send message : %s", e.Error())
	}

	responseStr := string(<-response)
	if responseStr == blank {
		responseStr = ""
	}

	return responseStr, nil
}
