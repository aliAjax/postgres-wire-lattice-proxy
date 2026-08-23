package pgwire

import (
	"errors"
	"fmt"
	"sync"
)

type Phase uint8

const (
	PhaseStartup Phase = iota
	PhaseAuth
	PhaseReady
	PhaseQuery
	PhaseExtended
	PhaseClosed
)

type StateMachine struct {
	mu     sync.Mutex
	phase  Phase
	failed error
	seen   uint64
}

func NewStateMachine() *StateMachine { return &StateMachine{phase: PhaseStartup} }
func (m *StateMachine) Phase() Phase { m.mu.Lock(); defer m.mu.Unlock(); return m.phase }
func (m *StateMachine) Accept(t byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failed != nil {
		return m.failed
	}
	m.seen++
	switch m.phase {
	case PhaseStartup:
		if t != 'Q' && t != 'P' && t != 'X' {
			m.failed = errors.New("message not allowed after startup")
		} else if t == 'Q' {
			m.phase = PhaseQuery
		} else {
			m.phase = PhaseExtended
		}
	case PhaseReady:
		if t == 'Q' {
			m.phase = PhaseQuery
		} else if t == 'P' || t == 'B' || t == 'D' || t == 'E' || t == 'S' {
			m.phase = PhaseExtended
		} else if t == 'X' {
			m.phase = PhaseClosed
		} else {
			m.failed = fmt.Errorf("message %q not allowed", t)
		}
	case PhaseQuery:
		if t != 'X' && t != 'Q' {
			m.failed = errors.New("query must be completed before next message")
		}
	case PhaseExtended:
		if t == 'S' {
			m.phase = PhaseReady
		} else if t == 'X' {
			m.phase = PhaseClosed
		}
	case PhaseClosed:
		return ioClosed()
	}
	return m.failed
}
func ioClosed() error { return errors.New("connection closed") }
func (m *StateMachine) Complete() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// A finished simple query or extended-query sequence returns the
	// connection to the idle/ready state so the next message is accepted.
	if m.phase == PhaseQuery || m.phase == PhaseExtended {
		m.phase = PhaseReady
	}
}
func (m *StateMachine) Seen() uint64 { m.mu.Lock(); defer m.mu.Unlock(); return m.seen }
