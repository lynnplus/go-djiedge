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
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"unsafe"
)

//export esdkCgoStreamCallback
func esdkCgoStreamCallback(ctx unsafe.Pointer, buf *C.uint8_t, size C.uint32_t) {
	if ctx == nil {
		return
	}
	lv := (*LiveView)(ctx)
	lv.onReceiveStream(buf, size)
}

//export esdkCgoStreamStatusCallback
func esdkCgoStreamStatusCallback(ctx unsafe.Pointer, value C.uint32_t) {
	if ctx == nil {
		return
	}
	s := int(value)
	status := &LiveStatus{
		Value:                 s,
		QualityAutoAvailable:  s&1 == 1,
		Quality540PAvailable:  s&2 == 2,
		Quality720PAvailable:  s&4 == 4,
		Quality720PHAvailable: s&8 == 8,
		Quality1080PAvailable: s&16 == 16,
	}
	lv := (*LiveView)(ctx)
	lv.onLiveStatusUpdate(status)
}

// The CameraSource type of stream source
type CameraSource int

func (c CameraSource) IsValid() bool {
	return c >= CameraSourceWide && c <= CameraSourceIR
}

const (
	CameraSourceWide CameraSource = iota + 1 //wide-angle lens camera
	CameraSourceZoom                         //zoom lens camera
	CameraSourceIR                           //infrared lens camera
)

// StreamQuality quality of the stream
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

// CameraType type of stream camera
type CameraType int

func (c CameraType) IsValid() bool {
	return c == CameraTypeFpv || c == CameraTypePayload
}

const (
	CameraTypeFpv CameraType = iota
	CameraTypePayload
)

// StreamReceiver is an interface for receiving stream data, status.
type StreamReceiver interface {
	// The OnStreamStatusUpdate is callback for LiveView stream-status.
	//
	//When the cloud is configured for live broadcast, the transmission channel is unbalanced, the aircraft is disconnected,
	//the video transmission has no signal and other factors will cause the code stream status to change.
	OnStreamStatusUpdate(status *LiveStatus)
	// The OnReceiveStreamData is a callback that receives stream data.
	//
	//Note: The 'data' parameter is a reference pointing to the C memory block.
	//you should use 'copy(dst,data)' or bytes.Buffer.Write() to copy it to the Go memory block,
	//directly using it as go bytes will cause panic.
	OnReceiveStreamData(data []byte)
}

type LiveView struct {
	native          *C.CEdgeLiveView
	streamReceiver  StreamReceiver
	cameraInitState atomic.Int32
}

// NewLiveView return a LiveView ptr that receives stream state and data.
//
// when LiveView is no longer used, LiveView.Destroy() should be called promptly to destroy it,
// or wait for garbage collected
// If LiveView is garbage collected, a finalizer may close the file descriptor,
// making it invalid; see runtime.SetFinalizer for more information on when
// a finalizer might be run.
func NewLiveView() *LiveView {
	lv := &LiveView{}
	p := C.Edge_LiveView_new(nil)
	lv.native = p
	lv.native.ctx = unsafe.Pointer(lv)
	runtime.SetFinalizer(lv, (*LiveView).Destroy)
	return lv
}

func (lv *LiveView) Destroy() {
	runtime.SetFinalizer(lv, nil)
	if lv.native != nil {
		lv.DeInit()
		C.Edge_LiveView_delete(lv.native)
		lv.native = nil
		lv.streamReceiver = nil
	}
}

// Init initialize live stream subscription.
// Note: For a specific camera, you can initialize only once
//
// handler implement data processing for received streams.
// allocator allocate data on received streams.
func (lv *LiveView) Init(cameraType CameraType, quality StreamQuality, handler StreamReceiver) error {
	if !cameraType.IsValid() || !quality.IsValid() {
		return errors.New("invalid parameter for camera or quality")
	}
	if handler == nil {
		return errors.New("parameter handler is nil")
	}

	opts := &C.CEdgeLiveViewOptions{
		camera:          C.int(cameraType),
		quality:         C.int(quality),
		stream_callback: C.CEdgeLiveViewStreamCallback(C.esdkCgoStreamCallback),
	}
	lv.streamReceiver = handler
	if lv.cameraInitState.CompareAndSwap(0, 1) {
		ret := C.Edge_LiveView_init(lv.native, opts)
		if err := convertCCodeToError(int(ret)); err != nil {
			lv.cameraInitState.Store(0)
			return err
		}
		lv.cameraInitState.Store(2)
	}

	if err := lv.setupStreamStatusCallback(); err != nil {
		return err
	}
	lv.streamReceiver = handler
	return nil
}

// DeInit de-initialize stream subscription
func (lv *LiveView) DeInit() {
	if lv.cameraInitState.CompareAndSwap(2, 0) {
		C.Edge_LiveView_deInit(lv.native)
	}
}

func (lv *LiveView) cameraInitialized() bool {
	return lv.cameraInitState.Load() == 2
}

func (lv *LiveView) onLiveStatusUpdate(status *LiveStatus) {
	if lv.streamReceiver != nil {
		lv.streamReceiver.OnStreamStatusUpdate(status)
	}
}

func (lv *LiveView) onReceiveStream(buf *C.uint8_t, size C.uint32_t) {
	if lv.streamReceiver == nil {
		return
	}

	//Note: only reference the memory data from cgo, no memory copy occurs
	data := unsafe.Slice((*byte)(buf), int(size))
	lv.streamReceiver.OnReceiveStreamData(data)
}

// SetCameraSource can switch the camera source used
func (lv *LiveView) SetCameraSource(source CameraSource) error {
	if !source.IsValid() {
		return errors.New("invalid parameter for camera source")
	}
	if !lv.cameraInitialized() {
		return errors.New(" live-view is not initialized")
	}

	ret := C.Edge_LiveView_setCameraSource(lv.native, C.int(source))
	return convertCCodeToError(int(ret))
}

func (lv *LiveView) setupStreamStatusCallback() error {
	ret := C.Edge_LiveView_subscribeStreamStatus(lv.native, C.CEdgeLiveViewStreamStatusCallback(C.esdkCgoStreamStatusCallback))
	return convertCCodeToError(int(ret))
}

// StartH264Stream start receive live H264 stream,stream data can be received through StreamReceiver.OnReceiveStreamData
func (lv *LiveView) StartH264Stream() error {
	if !lv.cameraInitialized() {
		return errors.New(" live-view is not initialized")
	}
	ret := C.Edge_LiveView_startH264Stream(lv.native)
	return convertCCodeToError(int(ret))
}

// StopH264Stream stop receive live H264 stream
func (lv *LiveView) StopH264Stream() error {
	ret := C.Edge_LiveView_stopH264Stream(lv.native)
	return convertCCodeToError(int(ret))
}

type LiveStatus struct {
	Value                 int
	QualityAutoAvailable  bool
	Quality540PAvailable  bool
	Quality720PAvailable  bool
	Quality720PHAvailable bool
	Quality1080PAvailable bool
}

func (l *LiveStatus) String() string {
	return fmt.Sprintf("value:%d auto:%v 540p:%v 720p:%v 720ph:%v 1080p:%v",
		l.Value,
		l.QualityAutoAvailable,
		l.Quality540PAvailable,
		l.Quality720PAvailable,
		l.Quality720PHAvailable,
		l.Quality1080PAvailable)
}
