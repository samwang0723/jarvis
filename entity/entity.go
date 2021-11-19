// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package entity

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

type ID uint64

const ZeroID = ID(0)

// use Sonyflake to support distributed unique IDs
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
	ID        ID         `gorm:"primaryKey" mapstructure:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" mapstructure:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" mapstructure:"updated_at"`
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
