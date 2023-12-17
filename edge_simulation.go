//go:build !linux || fake_edge

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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const fakeEdgeStreamFileName = "edge_stream.h264"

func InitSDK(device *DeviceInfo, auth *AuthInfo, key *RSA2048Key, logger *Logger, deInitOnFailed bool) (err error) {
	time.Sleep(5 * time.Second)
	return nil
}

// LiveView simulate edge device sending h264 data stream
type LiveView struct {
	handler StreamReceiver

	reading      atomic.Bool
	streamReader io.ReadCloser
	dataChan     chan []byte
	closeSig     chan bool

	wg sync.WaitGroup
}

func NewLiveView() *LiveView {
	return &LiveView{}
}

func (lv *LiveView) Init(cameraType CameraType, quality StreamQuality, handler StreamReceiver) error {
	if lv.handler != nil {
		return nil
	}
	lv.handler = handler
	go lv.updateState()
	return nil
}
func (lv *LiveView) SetCameraSource(source CameraSource) error {
	return nil
}

// StartH264Stream read local h264-stream file[edge_stream.h264] and push data to StreamReceiver
func (lv *LiveView) StartH264Stream() error {
	if !lv.reading.CompareAndSwap(false, true) {
		return nil
	}

	f, err := os.Open(fakeEdgeStreamFileName)
	if err != nil {
		lv.reading.Store(false)
		return err
	}
	lv.streamReader = f
	lv.closeSig = make(chan bool)
	lv.dataChan = make(chan []byte, 10)
	time.Sleep(1 * time.Second)

	lv.wg.Add(2)
	go func() {
		loopReadH264File(lv.streamReader, 30, lv.dataChan, lv.closeSig)
		lv.wg.Done()
	}()
	go func() {
		lv.pushStreamData()
		lv.wg.Done()
	}()
	return nil
}

func (lv *LiveView) StopH264Stream() error {
	if lv.streamReader == nil {
		return nil
	}
	if !lv.reading.CompareAndSwap(true, false) {
		return errors.New("state error")
	}
	lv.closeSig <- true
	lv.wg.Wait()

	close(lv.closeSig)
	close(lv.dataChan)
	_ = lv.streamReader.Close()
	lv.streamReader = nil
	return nil
}

func (lv *LiveView) updateState() {
	status := &LiveStatus{}
	c := time.Tick(2 * time.Second)
	i := 0
	for range c {
		if i == 3 {
			status.Value = 1
			status.QualityAutoAvailable = true
		}
		if i > 180 {
			i = 0
			status = &LiveStatus{}
		}
		i++
		if lv.handler != nil {
			lv.handler.OnStreamStatusUpdate(status)
		}
	}
}

func (lv *LiveView) pushStreamData() {
	for d := range lv.dataChan {
		if lv.handler != nil {
			lv.handler.OnReceiveStreamData(d)
		}
	}
}

func indexH264NaluStartCode(s []byte) (int, int) {
	start := 0
	for {
		i := bytes.Index(s, []byte{0x00, 0x00})
		if i < 0 {
			return -1, 0
		}

		if len(s) <= (i + 2) {
			return -1, 0
		}
		if s[i+2] == 0x01 {
			return start + i, 3
		}
		if len(s) <= (i + 3) {
			return -1, 0
		}
		if s[i+3] == 0x01 {
			return start + i, 4
		}
		s = s[1:]
		start++
	}
}

func loopReadH264File(reader io.Reader, frameInterval int, receiver chan []byte, closeSig <-chan bool) {
	r := bufio.NewScanner(reader)
	r.Buffer(nil, 1024*1024*2)
	r.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		start, w := indexH264NaluStartCode(data)
		if start < 0 {
			return 0, nil, nil
		}
		start += w
		end, w := indexH264NaluStartCode(data[start:])
		if end < 0 {
			return 0, nil, nil
		}
		return start + end, data[:start+end], nil
	})

	ticker := time.NewTicker(time.Duration(frameInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ok := r.Scan(); !ok {
				return
			}
			b := r.Bytes()
			if cap(receiver) <= len(receiver) {
				fmt.Println("fake stream block!!!")
				continue
			}
			receiver <- b
		case <-closeSig:
			return
		}
	}
}
