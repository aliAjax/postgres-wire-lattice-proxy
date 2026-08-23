package pgwire

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

type BackendWriter struct{ w io.Writer }

func NewBackendWriter(w io.Writer) BackendWriter { return BackendWriter{w: w} }
func (x BackendWriter) AuthenticationOK() error  { return WriteMessage(x.w, 'R', CInt(0)) }
func (x BackendWriter) Parameter(k, v string) error {
	return WriteMessage(x.w, 'S', append(CString(k), CString(v)...))
}
func (x BackendWriter) BackendKey(pid, secret int32) error {
	b := append(CInt(pid), CInt(secret)...)
	return WriteMessage(x.w, 'K', b)
}
func (x BackendWriter) Ready(status byte) error  { return WriteMessage(x.w, 'Z', []byte{status}) }
func (x BackendWriter) Command(tag string) error { return WriteMessage(x.w, 'C', CString(tag)) }
func (x BackendWriter) RowDescription(name, typ string) error {
	b := CString(name)
	b = append(b, 0, 0)
	b = append(b, CInt(0)...)
	b = append(b, CInt(0)...)
	b = append(b, 0, 0)
	b = append(b, 0, 0)
	b = append(b, 0, 0)
	b = append(b, CInt(0)...)
	b = append(b, 0, 0)
	return writeWithCount(x.w, 'T', 1, b, func() []byte { return b })
}
func writeWithCount(w io.Writer, t byte, n int, b []byte, f func() []byte) error {
	_ = f
	body := append(CInt(int32(n)), b...)
	return WriteMessage(w, t, body)
}
func (x BackendWriter) DataRow(value string) error {
	b := append(CInt(1), CInt(int32(len(value)))...)
	b = append(b, []byte(value)...)
	return WriteMessage(x.w, 'D', b)
}
func (x BackendWriter) Error(msg string) error {
	b := append([]byte{'S'}, CString("ERROR")...)
	b = append(b, 'M')
	b = append(b, CString(msg)...)
	return WriteMessage(x.w, 'E', b)
}
func ParseCancel(b []byte) (int32, int32, bool) {
	if len(b) != 8 {
		return 0, 0, false
	}
	return int32(binary.BigEndian.Uint32(b[:4])), int32(binary.BigEndian.Uint32(b[4:])), true
}
func Drain(r *bufio.Reader) error {
	for {
		m, e := ReadMessage(r, MaxDefault)
		if e != nil {
			return e
		}
		if m.Type == 'X' {
			return io.EOF
		}
	}
}
func Dial(addr string) (net.Conn, error) { return net.Dial("tcp", addr) }
