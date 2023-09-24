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

/*

#include "edge_liveview.h"

void esdkCgoStreamCallback(void* ctx,uint8_t *data, uint32_t dataLen);
void esdkCgoStreamStatusCallback(void* ctx,uint32_t value);
*/
import "C"
import (
	"unsafe"
)

//export esdkCgoStreamCallback
func esdkCgoStreamCallback(ctx unsafe.Pointer, buf *C.uint8_t, size C.uint32_t) {
	if ctx == nil {
		return
	}
	lv := (*LiveView)(ctx)
	if lv.streamCallback == nil {
		return
	}
	lv.streamCallback(nil, 0)
}

//export esdkCgoStreamStatusCallback
func esdkCgoStreamStatusCallback(ctx unsafe.Pointer, value C.uint32_t) {

}

type CameraSource int

func (c CameraSource) IsValid() bool {
	return c >= CameraSourceWide && c <= CameraSourceIR
}

const (
	CameraSourceWide CameraSource = iota + 1
	CameraSourceZoom
	CameraSourceIR
)

type StreamQuality int

func (s StreamQuality) IsValid() bool {
	return s >= StreamQuality540p && s <= StreamQuality1080p
}

const (
	// StreamQuality540p 30fps, 960*540, bps 512*1024
	StreamQuality540p StreamQuality = iota + 1

	// StreamQuality720p 30fps, 1280*720, bps 1024*1024
	StreamQuality720p

	// StreamQuality720pHigh 30fps, 1280*720, bps 1024*1024 + 512*1024
	StreamQuality720pHigh

	// StreamQuality1080p 30fps, 1920*1080, bps 3*1024*1024
	StreamQuality1080p
)

type CameraType int

func (c CameraType) IsValid() bool {
	return c == CameraTypeFpv || c == CameraTypePayload
}

const (
	CameraTypeFpv CameraType = iota
	CameraTypePayload
)

type StreamDataHandler func(data any, size uint)
type StreamStatusHandler func(data any, size uint)

type LiveView struct {
	native               *C.CEdgeLiveView
	streamCallback       StreamDataHandler
	streamStatusCallback StreamStatusHandler
}

func NewLiveView() *LiveView {
	lv := &LiveView{}
	p := C.Edge_LiveView_new()
	p.ctx = unsafe.Pointer(lv)
	lv.native = p
	return lv
}

func (lv *LiveView) Destroy() {
	C.Edge_LiveView_delete(lv.native)
	lv.native = nil
	lv.streamCallback = nil
	lv.streamStatusCallback = nil
}

func (lv *LiveView) Init(cameraType CameraType, quality StreamQuality, handler StreamDataHandler) error {
	opts := &C.CEdgeLiveViewOptions{
		camera:          C.int(cameraType),
		quality:         C.int(quality),
		stream_callback: C.CEdgeLiveViewStreamCallback(C.esdkCgoStreamCallback),
	}
	ret := C.Edge_LiveView_init(lv.native, opts)
	err := convertCCodeToError(int(ret))
	if err != nil {
		return err
	}
	lv.streamCallback = handler
	return nil
}

func (lv *LiveView) DeInit() {
	C.Edge_LiveView_deInit(lv.native)
}

func (lv *LiveView) SetCameraSource(source CameraSource) error {
	ret := C.Edge_LiveView_setCameraSource(lv.native, C.int(source))
	return convertCCodeToError(int(ret))
}

func (lv *LiveView) SubscribeStreamStatus(handler StreamStatusHandler) error {
	ret := C.Edge_LiveView_subscribeStreamStatus(lv.native, C.CEdgeLiveViewStreamStatusCallback(C.esdkCgoStreamStatusCallback))
	err := convertCCodeToError(int(ret))
	if err != nil {
		return err
	}
	lv.streamStatusCallback = handler
	return nil
}

func (lv *LiveView) StartH264Stream() error {
	ret := C.Edge_LiveView_startH264Stream(lv.native)
	return convertCCodeToError(int(ret))
}

func (lv *LiveView) StopH264Stream() error {
	ret := C.Edge_LiveView_stopH264Stream(lv.native)
	return convertCCodeToError(int(ret))
}
