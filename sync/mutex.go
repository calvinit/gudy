// Copyright (c) 2024 Calvin. All Rights Reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Mutex struct {
	sync.Mutex
}

// Count gets the metric of the number of mutex waiters and holders.
func (m *Mutex) Count() int {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	// The sum of the goroutines waiting to acquire and holding the mutex (0 or 1).
	count := state>>mutexWaiterShift + state&mutexLocked
	return int(count)
}

// IsLocked whether the mutex is locked.
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// IsWoken whether any waiters have been woken.
func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}
