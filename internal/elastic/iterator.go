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
package elastic

type Iterator interface {
	// Close the iterator and release any allocated resources.
	Close() error

	// Next loads the next document matching the search query.
	// returns false if no more Documents are available.
	Next() bool

	// Error returns the last error encountered by the iterator.
	Error() error

	// Document returns the current document from the result set.
	Document() *Document

	// TotalCount returns the approximate number of search results.
	TotalCount() uint64
}
