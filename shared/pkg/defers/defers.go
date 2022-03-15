package defers

import (
	"sync"
)

var (
	globalDefers = NewStack()
)

// Registers a defer function to be called at the end of the main function.
// It is recommended to call `Clean()` at the end of the main function.
func Register(fns ...func() error) {
	globalDefers.Push(fns...)
}

// Runs (Cleans) all defer functions.
// It is recommended to call this function at the end of the main function.
func Clean() {
	globalDefers.Clean()
}

type DeferStack struct {
	fns []func() error
	mu  sync.RWMutex
}

func NewStack() *DeferStack {
	return &DeferStack{
		fns: make([]func() error, 0),
		mu:  sync.RWMutex{},
	}
}

func (ds *DeferStack) Push(fns ...func() error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.fns = append(ds.fns, fns...)
}

func (ds *DeferStack) Clean() {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for i := len(ds.fns) - 1; i >= 0; i-- {
		if ds.fns[i] != nil {
			_ = ds.fns[i]()
		}
	}
}
