package deks

import (
	"encoding/binary"
	"time"

	"github.com/simia-tech/errx"
)

type container struct {
	key       []byte
	value     []byte
	revision  uint64
	deletedAt time.Time
}

func (c *container) delete() {
	c.deletedAt = time.Now()
}

func (c *container) undelete() {
	c.deletedAt = time.Time{}
}

func (c *container) isDeleted() bool {
	return !c.deletedAt.IsZero()
}

func (c *container) MarshalBinary() ([]byte, error) {
	keyLength := uint16(len(c.key))
	valueLength := len(c.value)
	buffer := make([]byte, int(keyLength)+valueLength+20)
	binary.BigEndian.PutUint64(buffer[:8], c.revision)
	binary.BigEndian.PutUint64(buffer[8:16], uint64(c.deletedAt.Unix()))
	binary.BigEndian.PutUint16(buffer[16:20], keyLength)
	copy(buffer[20:20+keyLength], c.key)
	copy(buffer[20+keyLength:], c.value)
	return buffer, nil
}

func (c *container) UnmarshalBinary(data []byte) error {
	if len(data) < 20 {
		return errx.Errorf("need at least 20 bytes")
	}
	c.revision = binary.BigEndian.Uint64(data[:8])
	c.deletedAt = time.Unix(int64(binary.BigEndian.Uint64(data[8:16])), 0)
	keyLength := binary.BigEndian.Uint16(data[16:20])
	c.key = data[20 : 20+keyLength]
	c.value = data[20+keyLength:]
	return nil
}
