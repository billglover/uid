// Package uid provides functions to generate roughly time
// ordered unique identifiers. The implementation is based
// on the MongoDB ObjectId specification.
// See https://docs.mongodb.com/manual/reference/method/ObjectId/
package uid

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
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
func NextID(hostname string, pid int) ([]byte, error) {

	// 12-byte IDs
	uid := make([]byte, 12)

	// 4-byte timestamps
	now := time.Now().Unix()
	binary.BigEndian.PutUint32(uid, uint32(now))

	// 3-byte machine identifiers
	hid := make([]byte, 3)

	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(hid, hw.Sum(nil))

	uid[4] = hid[0]
	uid[5] = hid[1]
	uid[6] = hid[2]

	// 2-byte process identifier
	binary.BigEndian.PutUint16(uid[7:9], uint16(pid))

	// 3-byte counter starting at a random number
	atomic.AddUint64(&seq, 1)
	uid[9] = byte(seq >> 16)
	uid[10] = byte(seq >> 8)
	uid[11] = byte(seq)

	// time should never go backwards, for now
	if now < prev {
		return uid, fmt.Errorf("we went back in time, wait for %dms", prev-now)
	}

	return uid, nil
}

// NextStringID returns the next uid in the sequence as a hexadecimal
// string or an error if a valid uid could not be generated.
func NextStringID(hostname string, pid int) (string, error) {
	id, err := NextID(hostname, pid)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(id), err
}
