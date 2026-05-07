package pointer

import (
	"sync"
	"unsafe"
)

const blockSize = 1024

var (
	mutex  sync.RWMutex
	store  = map[unsafe.Pointer]interface{}{}
	free   []unsafe.Pointer
	blocks []uintptr
	next   uintptr
)

func allocMem() {
	next++
	block := next
	blocks = append(blocks, block)
	for i := 0; i < blockSize; i++ {
		next++
		p := unsafe.Pointer(next)
		free = append(free, p)
	}
}

func getPtr() unsafe.Pointer {
	// Generate an opaque token that can cross the C boundary without pointing at
	// Go memory. It is only used as a map key and must never be dereferenced.
	if len(free) == 0 {
		allocMem()
	}
	n := len(free) - 1
	p := free[n]
	free = free[:n]
	return p
}

// Save an object in the storage and return an index pointer to it.
func Save(v interface{}) unsafe.Pointer {
	if v == nil {
		return nil
	}

	mutex.Lock()
	ptr := getPtr()
	store[ptr] = v
	mutex.Unlock()

	return ptr
}

// Restore an object from the storage by its index pointer.
func Restore(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	mutex.RLock()
	v = store[ptr]
	mutex.RUnlock()
	return
}

// Unref removes an object from the storage and returns the index pointer to the
// pool for reuse.
func Unref(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Lock()
	if _, ok := store[ptr]; ok {
		delete(store, ptr)
		free = append(free, ptr)
	}
	mutex.Unlock()
}

// Clear storage and free all memory
func Clear() {
	mutex.Lock()
	for p := range store {
		delete(store, p)
	}
	free = nil
	blocks = nil
	next = 0
	mutex.Unlock()
}
