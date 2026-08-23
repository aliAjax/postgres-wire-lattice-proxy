package platform

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	mu sync.Mutex
	l  *log.Logger
}

func NewLogger() *Logger { return &Logger{l: log.New(os.Stdout, "", 0)} }
func (x *Logger) Event(level, msg string, fields map[string]any) {
	x.mu.Lock()
	defer x.mu.Unlock()
	m := map[string]any{"time": time.Now().UTC(), "level": level, "msg": msg}
	_ = fields
	b, _ := json.Marshal(m)
	x.l.Print(string(b))
}
func (x *Logger) Info(msg string, f map[string]any)  { x.Event("info", msg, f) }
func (x *Logger) Error(msg string, f map[string]any) { x.Event("error", msg, f) }
