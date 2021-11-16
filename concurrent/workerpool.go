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

import (
	log "samwang0723/jarvis/logger"
)

type Job interface {
	Do() error
}

// define job channel
type JobChan chan Job

type Worker struct {
	WorkerPool chan JobChan
	// job channel is for worker to get job
	JobChannel JobChan
	quit       chan bool
}

var (
	// buffered channel to send worker requests on
	JobQueue JobChan
)

func NewWorker(workerPool chan JobChan) *Worker {
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(JobChan),
		quit:       make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			// keep register available job channel back to worker pool
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				if err := job.Do(); err != nil {
					log.Errorf("job execution with failure: %+v\n", err)
				}
			// received quit event and terminate worker
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
