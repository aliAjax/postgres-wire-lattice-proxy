package pgwire

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"errors"
	"io"
	"net"
	"time"
)

const MaxDefault = 1 << 20

var ErrTooLarge = errors.New("postgres message too large")

type Message struct {
	Type byte
	Body []byte
}

func ReadMessage(r *bufio.Reader, max int) (Message, error) {
	t, e := r.ReadByte()
	if e != nil {
		return Message{}, e
	}
	var h [4]byte
	if _, e = io.ReadFull(r, h[:]); e != nil {
		return Message{}, e
	}
	n := int(binary.BigEndian.Uint32(h[:])) - 4
	if n < 0 || n > max {
		return Message{}, fmt.Errorf("message rejected: %v", ErrTooLarge)
	}
	b := make([]byte, n)
	if _, e = io.ReadFull(r, b); e != nil {
		return Message{}, e
	}
	return Message{t, b}, nil
}
func WriteMessage(w io.Writer, t byte, b []byte) error {
	if len(b) > MaxDefault {
		return ErrTooLarge
	}
	h := []byte{t, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(h[1:], uint32(len(b)+4))
	if _, e := w.Write(h); e != nil {
		return e
	}
	_, e := w.Write(b)
	return e
}
func Startup(r *bufio.Reader, max int) (map[string]string, error) {
	var h [4]byte
	if _, e := io.ReadFull(r, h[:]); e != nil {
		return nil, e
	}
	n := int(binary.BigEndian.Uint32(h[:]))
	if n < 8 || n > max {
		return nil, ErrTooLarge
	}
	b := make([]byte, n-4)
	if _, e := io.ReadFull(r, b); e != nil {
		return nil, e
	}
	params := map[string]string{}
	if len(b) < 4 {
		return params, nil
	}
	for p := b[4:]; len(p) > 0 && p[0] != 0; {
		k, rest := splitCString(p)
		if rest == nil {
			return nil, errors.New("invalid startup key")
		}
		v, rest := splitCString(rest)
		if rest == nil {
			return nil, errors.New("invalid startup value")
		}
		params[string(k)] = string(v)
		p = rest
	}
	return params, nil
}
func splitCString(b []byte) ([]byte, []byte) {
	for i, c := range b {
		if c == 0 {
			return b[:i], b[i+1:]
		}
	}
	return nil, nil
}
func StartupPacket(params map[string]string) []byte {
	body := make([]byte, 4, 128)
	binary.BigEndian.PutUint32(body, 196608)
	for k, v := range params {
		body = append(body, []byte(k)...)
		body = append(body, 0)
		body = append(body, []byte(v)...)
		body = append(body, 0)
	}
	body = append(body, 0)
	out := make([]byte, 4+len(body))
	binary.BigEndian.PutUint32(out, uint32(len(out)))
	copy(out[4:], body)
	return out
}
func ReadWithDeadline(conn net.Conn, max int, d time.Duration) (Message, error) {
	if e := conn.SetReadDeadline(time.Now().Add(d)); e != nil {
		return Message{}, e
	}
	return ReadMessage(bufio.NewReader(conn), max)
}
func I32(b []byte) int32 {
	if len(b) < 4 {
		return 0
	}
	return int32(binary.BigEndian.Uint32(b))
}
func CInt(v int32) []byte     { b := make([]byte, 4); binary.BigEndian.PutUint32(b, uint32(v)); return b }
func CString(s string) []byte { return append([]byte(s), 0) }
func ReadCString(b []byte) (string, []byte, error) {
	x, y := splitCString(b)
	if x == nil {
		return "", nil, fmt.Errorf("unterminated cstring")
	}
	return string(x), y, nil
}
