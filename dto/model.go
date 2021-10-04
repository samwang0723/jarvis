package dto

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

type ID uint64

const ZeroID = ID(0)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("failed to init sonyflake for distributed ID generation")
	}
}

func GenID() (ID, error) {
	rawID, err := sf.NextID()
	if err != nil {
		return ZeroID, fmt.Errorf("Cannot get sf.NextID(): %s\n", err)
	}

	id := ID(rawID)
	if id == ZeroID {
		return ZeroID, fmt.Errorf("ZeroID was generated")
	}

	return id, nil
}

func (id ID) Uint64() uint64 {
	return uint64(id)
}

type Model struct {
	ID        ID         `json:"ID"`
	CreatedAt *time.Time `json:"CreatedAt"`
	UpdatedAt *time.Time `json:"UpdatedAt"`
}
