// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package entity

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/samwang0723/jarvis/internal/helper"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

type ID uint64

const (
	ZeroID = ID(0)
	MyIP   = "MY_IP"
)

//nolint:nolintlint, gochecknoglobals
var generator *sonyflake.Sonyflake

//nolint:nolintlint, gochecknoinits
func init() {
	var st sonyflake.Settings
	if helper.GetCurrentEnv() == "prod" {
		st.MachineID = machineID
	}

	generator = sonyflake.NewSonyflake(st)
	if generator == nil {
		panic("failed to init sonyflake for distributed ID generation")
	}
}

func GenID() (ID, error) {
	rawID, err := generator.NextID()
	if err != nil {
		return ZeroID, fmt.Errorf("cannot get sf.NextID(): %w", err)
	}

	id := ID(rawID)
	if id == ZeroID {
		return ZeroID, fmt.Errorf("zeroID was generated")
	}

	return id, nil
}

func (id ID) Uint64() uint64 {
	return uint64(id)
}

func machineID() (uint16, error) {
	ipStr := os.Getenv(MyIP)
	if ipStr == "" {
		return 0, errors.New("'MY_IP' environment variable not set")
	}

	ip := net.ParseIP(ipStr)
	if ip == nil || len(ip) < 16 {
		return 0, errors.New("invalid IP")
	}

	return uint16(ip[8])<<7 + uint16(ip[9])<<6 +
			uint16(ip[10])<<5 + uint16(ip[11])<<4 +
			uint16(ip[12])<<3 + uint16(ip[13])<<2 +
			uint16(ip[14])<<1 + uint16(ip[15]),
		nil
}

type Model struct {
	ID        ID         `gorm:"primaryKey" mapstructure:"id"`
	CreatedAt *time.Time `gorm:"column:created_at" mapstructure:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" mapstructure:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" mapstructure:"deleted_at"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == ZeroID {
		m.ID, err = GenID()
	}

	if err != nil {
		return err
	}

	return nil
}
