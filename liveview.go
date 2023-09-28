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
	"io"
	"runtime"
	"unsafe"
)

//export esdkCgoStreamCallback
func esdkCgoStreamCallback(ctx unsafe.Pointer, buf *C.uint8_t, size C.uint32_t) {
	if ctx == nil {
		return
	}
	lv := (*LiveView[io.Writer])(ctx)
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
	lv := (*LiveView[io.Writer])(ctx)
	lv.onLiveStatusUpdate(status)
}

type CameraSource int

func (c CameraSource) IsValid() bool {
	return c >= CameraSourceWide && c <= CameraSourceIR
}

// The list of Stream source of payload camera
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

// CameraType type of Stream Camera
type CameraType int

func (c CameraType) IsValid() bool {
	return c == CameraTypeFpv || c == CameraTypePayload
}

const (
	CameraTypeFpv CameraType = iota
	CameraTypePayload
)

// StreamReceiver is an interface for receiving stream data, status and allocating stream data objects.
type StreamReceiver[T io.Writer] interface {
	OnStreamStatusUpdate(status *LiveStatus)
	OnReceiveStreamData(data T)
	// AllocateStreamData allocate a data object for method OnReceiveStreamData
	// memory data from cgo will be written to data using the method ‘Write’,
	// the bytes data passed into method 'Write' cannot modify its content because it points to the underlying cgo memory block,
	// data should be copied using methods such as copy or byte.Buffer
	//
	// Examples:
	// func (b *T) Write(p []byte) (n int, err error) {
	//     buf:=make([]byte,len(p))
	//     return copy(buf,p), nil
	// }
	// or use
	// buf:=new(bytes.Buffer)
	// return buf.Write(p)
	//
	// because a large amount of bytes data will be generated, considering the gc pressure,
	// recommended to use object caching technology optimization such as sync.Pool.
	AllocateStreamData() T
}

type LiveView[T io.Writer] struct {
	native         *C.CEdgeLiveView
	streamReceiver StreamReceiver[T]
}

// NewLiveView create a *LiveView object that receives stream state and data.
// the generic T represents the stream data type that needs to be returned in the interface StreamReceiver.
//
// when *LiveView is no longer used, *LiveView.Destroy() should be called promptly to destroy it,
// or wait for garbage collected
// If *LiveView is garbage collected, a finalizer may close the file descriptor,
// making it invalid; see runtime.SetFinalizer for more information on when
// a finalizer might be run.
func NewLiveView[T io.Writer]() *LiveView[T] {
	lv := &LiveView[T]{}
	p := C.Edge_LiveView_new()
	p.ctx = unsafe.Pointer(lv)
	lv.native = p
	runtime.SetFinalizer(lv, (*LiveView[T]).Destroy)
	return lv
}

func (lv *LiveView[T]) Destroy() {
	C.Edge_LiveView_delete(lv.native)
	lv.native = nil
	lv.streamReceiver = nil
	runtime.SetFinalizer(lv, nil)
}

// Init initialize live stream subscription.
// Note: For a specific camera, you can initialize only once
//
// handler implement data processing for received streams.
// allocator allocate data on received streams.
func (lv *LiveView[T]) Init(cameraType CameraType, quality StreamQuality, handler StreamReceiver[T]) error {
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
	ret := C.Edge_LiveView_init(lv.native, opts)
	if err := convertCCodeToError(int(ret)); err != nil {
		return err
	}
	if err := lv.setupStreamStatusCallback(); err != nil {
		return err
	}
	lv.streamReceiver = handler
	return nil
}

// DeInit de-initialize stream subscription
func (lv *LiveView[T]) DeInit() {
	C.Edge_LiveView_deInit(lv.native)
}

func (lv *LiveView[T]) onLiveStatusUpdate(status *LiveStatus) {
	if lv.streamReceiver != nil {
		lv.streamReceiver.OnStreamStatusUpdate(status)
	}
}
func (lv *LiveView[T]) onReceiveStream(buf *C.uint8_t, size C.uint32_t) {
	if lv.streamReceiver == nil {
		return
	}

	d := lv.streamReceiver.AllocateStreamData()
	if any(d) == nil {
		return
	}
	//Note: only reference the memory data from cgo, no memory copy occurs
	data := unsafe.Slice((*byte)(buf), int(size))
	//ignore error
	_, _ = d.Write(data)
	lv.streamReceiver.OnReceiveStreamData(d)
}

// SetCameraSource can switch the camera source used
func (lv *LiveView[T]) SetCameraSource(source CameraSource) error {
	if !source.IsValid() {
		return errors.New("invalid parameter for camera source")
	}
	ret := C.Edge_LiveView_setCameraSource(lv.native, C.int(source))
	return convertCCodeToError(int(ret))
}

func (lv *LiveView[T]) setupStreamStatusCallback() error {
	ret := C.Edge_LiveView_subscribeStreamStatus(lv.native, C.CEdgeLiveViewStreamStatusCallback(C.esdkCgoStreamStatusCallback))
	return convertCCodeToError(int(ret))
}

// StartH264Stream start receive live H264 stream,stream data can be received through StreamReceiver.OnStreamStatusUpdate
func (lv *LiveView[T]) StartH264Stream() error {
	ret := C.Edge_LiveView_startH264Stream(lv.native)
	return convertCCodeToError(int(ret))
}

// StopH264Stream stop receive live H264 stream
func (lv *LiveView[T]) StopH264Stream() error {
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
