//go:build linux

/*
 * Copyright (c) 2023 Lynn <lynnplus90@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package djiedge

import (
	"runtime"
)

type fileHandle int32

type MediaFile struct {
	path   string
	handle fileHandle
	reader *MediaFileReader
}

func (m *MediaFile) setup() {
	runtime.SetFinalizer(m, (*MediaFile).Close)
}

func (m *MediaFile) Path() string {
	return m.path
}

func (m *MediaFile) Close() error {
	runtime.SetFinalizer(m, nil)
	return m.reader.closeFile(m.handle)
}

func (m *MediaFile) Read(p []byte) (n int, err error) {
	return m.reader.readFile(m.handle, p)
}
