package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/acme/pg-lattice-proxy/internal/pgwire"
	"net"
	"os"
	"time"
)

func main() {
	addr := "127.0.0.1:6432"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	c, e := net.DialTimeout("tcp", addr, time.Second)
	if e != nil {
		panic(e)
	}
	defer c.Close()
	if _, e = c.Write(pgwire.StartupPacket(map[string]string{"user": "smoke", "database": "postgres", "application_name": "pgwire-smoke"})); e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)
	for i := 0; i < 5; i++ {
		m, e := pgwire.ReadMessage(r, pgwire.MaxDefault)
		if e != nil {
			panic(e)
		}
		fmt.Printf("startup response %c len=%d\n", m.Type, len(m.Body))
		if m.Type == 'Z' {
			break
		}
	}
	q := append([]byte("SELECT 1"), 0)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(len(q)+4))
	c.Write(append(append([]byte{'Q'}, b...), q...))
	for i := 0; i < 5; i++ {
		m, e := pgwire.ReadMessage(r, pgwire.MaxDefault)
		if e != nil {
			panic(e)
		}
		fmt.Printf("query response %c len=%d\n", m.Type, len(m.Body))
		if m.Type == 'Z' {
			break
		}
	}
}
