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
package concurrent

type Dispatcher struct {
	// a pool of worker channels that registered with dispatcher
	workerPool chan JobChan
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan JobChan, maxWorkers)
	return &Dispatcher{
		workerPool: pool,
		maxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		// job request received
		case job := <-JobQueue:
			go func(job Job) {
				// try to obtain a available worker job channel
				// this will block until a worker is idle
				jobChannel := <-d.workerPool

				// dispatch the job to worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
