// Copyright (c) 2024 Calvin. All Rights Reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//go:build !go1.9
// +build !go1.9

package sync

import (
	"sync/atomic"
	"unsafe"
)

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexWaiterShift = iota
)

// TryLock tries to lock m and reports whether it succeeded.
func (m *Mutex) TryLock() bool {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// If it is locked, or woken, we are not contended and simply return false.
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexWoken) != 0 {
		return false
	}

	// Attempt to acquire the mutex in a race state.
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, old|mutexLocked)
}
