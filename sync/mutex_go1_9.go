// Copyright (c) 2024 Calvin. All Rights Reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//go:build go1.9 && !go1.18
// +build go1.9,!go1.18

package sync

import (
	"sync/atomic"
	"unsafe"
)

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
)

// TryLock tries to lock m and reports whether it succeeded.
func (m *Mutex) TryLock() bool {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// If it is locked, starving, or woken, we are not contended and simply return false.
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	// Attempt to acquire the mutex in a race state.
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, old|mutexLocked)
}

// IsStarving whether the mutex is starving.
func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}
