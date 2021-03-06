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

package main

import (
	"log"
	"time"

	"github.com/samwang0723/jarvis/internal/app/server"
	"github.com/samwang0723/jarvis/internal/helper"
)

func main() {
	// manually set time zone, docker image may not have preset timezone
	var err error
	time.Local, err = time.LoadLocation(helper.TimeZone)
	if err != nil {
		log.Printf("error loading location '%s': %v\n", helper.TimeZone, err)
	}

	server.Serve()
}
