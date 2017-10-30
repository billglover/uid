// Package uid provides functions to generate roughly time
// ordered unique identifiers. The implementation is based
// on Twitter's now abandoned Snowflake.
// See https://github.com/twitter/snowflake/tree/b3f6a3c6ca8e1b6847baa6ff42bf72201e2c2231
//
// TODO: the instance identifier currently defaults to 0
package uid

import (
	"fmt"
	"sync"
	"time"
)

const (
	seqBits  uint64 = 12
	instBits uint64 = 10
)

var (
	mu   sync.Mutex // we need a lock to ensure concurrency safe
	inst uint64     // an instance identifier, 10-bit precision
	seq  uint64     // a sequence number, 12-bit precision
	prev int64      // previous timestamp
)

// NextID returns the next uid in the sequence or an error if
// a valid uid could not be generated.
func NextID() (uint64, error) {

	mu.Lock()
	defer mu.Unlock()

	// get the current time in milliseconds
	now := time.Now().UnixNano() / 1e6

	// time should never go backwards, for now
	if now < prev {
		return 0, fmt.Errorf("we went back in time, wait for %dms", prev-now)
	}

	// increment the sequence number if the timestamp
	// hasn't changed since the last ID was generated
	if now == prev {
		seq++
	} else {
		seq = 0
	}
	prev = now

	// generate the uid
	uid := uint64(now)<<(seqBits+instBits) | inst<<seqBits | seq
	return uid, nil
}

// NextStringID returns the next uid in the sequence as a hexadecimal
// string or an error if a valid uid could not be generated.
func NextStringID() (string, error) {
	id, err := NextID()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%#x", id), err
}
