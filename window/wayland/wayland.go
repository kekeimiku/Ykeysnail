package wayland

import (
	"encoding/json"
	"fmt"
	"net"
	"unsafe"
)

const (
	MAGICK        string = "i3-ipc"
	HEADERLEN     int    = 14
	IPC_SUBSCRIBE int    = 2
)

type SwayConnection struct {
	Conn net.Conn
}

func Conn(path string) (*SwayConnection, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &SwayConnection{Conn: conn}, nil
}

func (sc *SwayConnection) Raw(messageType int) ([]byte, error) {
	var (
		message  = []byte(MAGICK)
		payload  = []byte(`["window"]`)
		length   = int32(len(payload))
		bytelen  [4]byte
		bytetype [4]byte
	)

	bytelen = *(*[4]byte)(unsafe.Pointer(&length))
	bytetype = *(*[4]byte)(unsafe.Pointer(&messageType))

	for _, b := range bytelen {
		message = append(message, b)
	}
	for _, b := range bytetype {
		message = append(message, b)
	}

	message = append(message, payload...)

	_, err := sc.Conn.Write(message)
	if err != nil {
		return []byte(nil), err
	}

	msg, err := sc.Response()
	if err != nil {
		return []byte(nil), err
	}
	return msg, nil
}

func (sc *SwayConnection) Response() ([]byte, error) {
	header := make([]byte, HEADERLEN)
	n, err := sc.Conn.Read(header)

	// Check if this is a valid i3 message.
	if n != HEADERLEN || err != nil {
		return []byte(nil), err
	}

	magicString := string(header[:len(MAGICK)])
	if magicString != MAGICK {
		err = fmt.Errorf(
			"invalid magic string: got %q, expected %q",
			magicString, MAGICK)
		return []byte(nil), err
	}

	var bytelen [4]byte

	for i, b := range header[len(MAGICK) : len(MAGICK)+4] {
		bytelen[i] = b
	}
	length := *(*int32)(unsafe.Pointer(&bytelen))

	payload := make([]byte, length)
	n, err = sc.Conn.Read(payload)
	if n != int(length) || err != nil {
		return []byte(nil), err
	}

	var bytetype [4]byte
	for i, b := range header[len(MAGICK)+4 : len(MAGICK)+8] {
		bytetype[i] = b
	}

	return payload, err
}

type Subscription struct {
	Events chan Event
	Errors chan error
	quit   chan struct{}
}

type Event struct {
	Container struct {
		ID    int         `json:"id"`
		Name  string      `json:"name"`
		Pid   int         `json:"pid"`
		AppID interface{} `json:"app_id"`
	} `json:"container"`
}

func (sc *SwayConnection) Subscribe() *Subscription {
	subscription := &Subscription{make(chan Event), make(chan error), make(chan struct{})}

	go func() {
		for {
			select {
			case <-subscription.quit:
				return
			default:
				var event Event
				o, err := sc.Response()
				if err != nil {
					subscription.Errors <- err
				}
				err = json.Unmarshal(o, &event)
				if err != nil {
					subscription.Errors <- err
					continue
				}
				subscription.Events <- event
			}
		}
	}()

	return subscription
}
func (s *Subscription) Close() {
	s.quit <- <-make(chan struct{})
}

var SubscriptionSwayIPC *Subscription

func init() {
	conn, err := Conn("/run/user/1000/sway-ipc.1000.1153.sock")
	if err != nil {
		panic(err)
	}

	_, err = conn.Raw(IPC_SUBSCRIBE)
	if err != nil {
		panic(err)
	}

	SubscriptionSwayIPC = conn.Subscribe()

	go func() {
		defer SubscriptionSwayIPC.Close()
	}()

}
