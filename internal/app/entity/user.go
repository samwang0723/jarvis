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

import "time"

type User struct {
	EmailConfirmedAt *time.Time `gorm:"column:email_confirmed_at" mapstructure:"email_confirmed_at"`
	PhoneConfirmedAt *time.Time `gorm:"column:phone_confirmed_at" mapstructure:"phone_confirmed_at"`
	SessionExpiredAt *time.Time `gorm:"column:session_expired_at" mapstructure:"session_expired_at"`
	FirstName        string     `gorm:"column:first_name"                                           json:"firstName"`
	LastName         string     `gorm:"column:last_name"                                            json:"lastName"`
	Email            string     `gorm:"column:email"                                                json:"email"`
	Phone            string     `gorm:"column:phone"                                                json:"phone"`
	Password         string     `gorm:"column:password"                                             json:"password"`
	SessionID        string     `gorm:"column:session_id"                                           json:"sessionID"`
	Base
}

func (User) TableName() string {
	return "users"
}
