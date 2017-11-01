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
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	seqBits  uint64 = 12
	instBits uint64 = 10
)

type UID struct {
	hid []byte
	pid uint16
}

var (
	mu   sync.Mutex // we need a lock to ensure concurrency safe
	inst uint64     // an instance identifier, 10-bit precision
	seq  uint64     // a sequence number, 12-bit precision
	prev int64      // previous timestamp
)

// NewGenerator returns an instance of a UID generator. During
// initialisation it generates the machine and process identifiers.
func NewGenerator() *UID {

	// generate a 3-byte machine identifier from the hostname
	// or random bytes
	hid := make([]byte, 3)

	hostname, err := os.Hostname()
	if err != nil {
		rand.Read(hid)
	} else {
		hw := md5.New()
		hw.Write([]byte(hostname))
		copy(hid, hw.Sum(nil))
	}

	// generate a process identifier
	pid := os.Getpid()

	// seed the sequence identifier
	seq = uint64(rand.Int63())

	uid := UID{
		hid: hid,
		pid: uint16(pid),
	}
	return &uid
}

// NextID returns the next uid in the sequence or an error if
// a valid uid could not be generated.
func (uid *UID) NextID() ([]byte, error) {

	// 12-byte IDs
	id := make([]byte, 12)

	// 4-byte timestamps
	now := time.Now().Unix()
	binary.BigEndian.PutUint32(id, uint32(now))

	// 3-byte machine identifier
	id[4] = uid.hid[0]
	id[5] = uid.hid[1]
	id[6] = uid.hid[2]

	// time should never go backwards, for now
	if now < prev {
		return id, fmt.Errorf("stop the clock, we went back in time, wait for %dms", prev-now)
	}

	// 2-byte process identifier
	binary.BigEndian.PutUint16(id[7:9], uid.pid)

	// 3-byte counter starting at a random number
	atomic.AddUint64(&seq, 1)
	id[9] = byte(seq >> 16)
	id[10] = byte(seq >> 8)
	id[11] = byte(seq)

	return id, nil
}

// NextStringID returns the next uid in the sequence as a hexadecimal
// string or an error if a valid uid could not be generated.
func (uid *UID) NextStringID() (string, error) {
	id, err := uid.NextID()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(id), err
}
