// Copyright (c) 2024 Calvin. All Rights Reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/calvinit/gudy/sync"
)

func TestMutex_TryLock(t *testing.T) {
	var mu sync.Mutex
	// Starting a goroutine holds the lock for a random amount of time (< 2 secs) before releasing it.
	go func() {
		mu.Lock()
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		mu.Unlock()
	}()

	time.Sleep(time.Second)

	ok := mu.TryLock()
	if ok {
		fmt.Println("got the lock.")
		// do something...
		mu.Unlock()
		return
	}
	fmt.Println("can't get the lock.")
}

func TestMutex_Count_IsLocked_IsWoken_IsStarving(t *testing.T) {
	var mu sync.Mutex
	for i := 0; i < 1000; i++ {
		go func() {
			mu.Lock()
			time.Sleep(time.Second)
			mu.Unlock()
		}()
	}

	time.Sleep(time.Second)

	// mutex legacy
	fmt.Printf("waiters(and holders): %d, locked: %t, woken: %t.\n", mu.Count(), mu.IsLocked(), mu.IsWoken())

	// fmt.Printf("waiters(and holders): %d, locked: %t, woken: %t, starving: %t.\n",
	//     mu.Count(), mu.IsLocked(), mu.IsWoken(), mu.IsStarving())
}
