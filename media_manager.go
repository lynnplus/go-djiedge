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

/*

#include "edge_media.h"

void esdkCgoNewMediaFileCallback(CEdgeMediaFile*);
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

//export esdkCgoNewMediaFileCallback
func esdkCgoNewMediaFileCallback(f *C.CEdgeMediaFile) {
	if sdkNewMFObserver == nil {
		return
	}
	desc := convertToMFDesc(f)
	sdkNewMFObserver(desc)
}

func convertToMFDesc(f *C.CEdgeMediaFile) *MediaFileDesc {
	return &MediaFileDesc{
		FileName:         convertToGoString(f.file_name),
		FilePath:         convertToGoString(f.file_path), //C.GoStringN(f.file_path.data, C.int(f.file_path.len)),
		FileSize:         uint64(f.file_size),
		FileType:         MediaFileType(f.file_type),
		CameraAttr:       CameraAttr(f.camera_attr),
		Latitude:         float64(f.latitude),
		Longitude:        float64(f.longitude),
		AbsoluteAltitude: float64(f.absolute_altitude),
		RelativeAltitude: float64(f.relative_altitude),
		GimbalYawDegree:  float64(f.gimbal_yaw_degree),
		ImageWidth:       int(f.image_width),
		ImageHeight:      int(f.image_height),
		VideoDuration:    time.Duration(uint32(f.video_duration)) * time.Second,
		CreateTime:       time.Unix(int64(f.create_time), 0),
	}
}

var sdkNewMFObserver func(desc *MediaFileDesc)

// RegisterMediaFilesObserver
// register media file notification processing callback.
func RegisterMediaFilesObserver(observer func(desc *MediaFileDesc)) error {
	ret := C.Edge_MediaMgr_registerMediaFilesObserver(C.CEdgeMediaFilesObserver(C.esdkCgoNewMediaFileCallback))
	if err := convertCCodeToError(int(ret)); err != nil {
		return err
	}
	sdkNewMFObserver = observer
	return nil
}

// SetDroneNestUploadCloud
// Set media files for cloud upload.
// Note: Media files from flight waylines are uploaded to the cloud by default.
// If set to not upload, and edge computing goes offline for over 30s, the default method will be restored,
// restarting to upload media files to the cloud.
func SetDroneNestUploadCloud(enable bool) error {
	if !Initialized() {
		return ErrSDKNotInit
	}
	ret := C.Edge_MediaMgr_setDroneNestUploadCloud(C.bool(enable))
	return convertCCodeToError(int(ret))
}

// SetDroneNestAutoDelete
// set whether to delete local media files at the dock after uploading is complete.
// Note: When edge computing requires media file retrieval, it should be set not to delete.
func SetDroneNestAutoDelete(enable bool) error {
	if !Initialized() {
		return ErrSDKNotInit
	}
	ret := C.Edge_MediaMgr_setDroneNestAutoDelete(C.bool(enable))
	return convertCCodeToError(int(ret))
}

type MediaFileReader struct {
	native unsafe.Pointer
	status atomic.Int32 //0:closed  1:opening  2:opened 3:closing
}

// NewMediaFileReader return a media file reader
func NewMediaFileReader() *MediaFileReader {
	p := C.Edge_MediaMgr_createMediaFilesReader()
	r := &MediaFileReader{native: p}
	runtime.SetFinalizer(r, (*MediaFileReader).Destroy)
	return r
}

func (m *MediaFileReader) Destroy() {
	runtime.SetFinalizer(m, nil)
	m.status.Store(0)
	if m.native != nil {
		C.Edge_MediaMgr_deleteMediaFilesReader(m.native)
		m.native = nil
	}
}

// Open establish a media file transfer connection with the dock.
// can only access media files after successful initialization.
// note: This interface sets the dock's local media file strategy to not delete.
func (m *MediaFileReader) Open() error {
	if !Initialized() {
		return ErrSDKNotInit
	}

	if !m.status.CompareAndSwap(0, 1) {
		return errors.New("status abnormal")
	}
	ret := C.Edge_MFReader_init(m.native)
	if err := convertCCodeToError(int(ret)); err != nil {
		m.status.Store(0)
		return err
	}
	m.status.Store(2)
	return nil
}

// Close disconnect the media file transfer connection with the dock.
// when no longer need to pull media files, call this interface to disconnect.
// after disconnection, need to reinitialize to pull media files.
func (m *MediaFileReader) Close() error {
	if !m.status.CompareAndSwap(2, 3) {
		return errors.New("status abnormal")
	}

	ret := C.Edge_MFReader_deInit(m.native)
	if err := convertCCodeToError(int(ret)); err != nil {
		m.status.Store(2)
		return err
	}
	m.status.Store(0)
	return nil
}

func (m *MediaFileReader) IsOpened() bool {
	return m.status.Load() == 2
}

// GetFileList gets the media file list from the most recent wayline mission.
// Note: Only executed wayline missions that take photos and videos to produce media files,
// and have the delete strategy set to not delete,
// can retrieve the dock media file list through this interface.
func (m *MediaFileReader) GetFileList() ([]*MediaFileDesc, error) {
	if !m.IsOpened() {
		return nil, ErrFileReaderNotOpen
	}

	var arr *C.CEdgeMediaFile
	n := int(C.Edge_MFReader_fileList(m.native, &arr))
	if n <= 0 {
		return nil, nil
	}
	defer C.free(unsafe.Pointer(arr))
	list := unsafe.Slice(arr, n)
	ret := make([]*MediaFileDesc, n)

	for i, file := range list {
		ret[i] = convertToMFDesc(&file)
	}
	return ret, nil
}

// OpenFile returns an opened *MediaFile.
// The parameter 'path' from MediaFileDesc.FilePath
func (m *MediaFileReader) OpenFile(path string) (*MediaFile, error) {
	if !m.IsOpened() {
		return nil, ErrFileReaderNotOpen
	}

	p := convertToCString(path)
	fh := int(C.Edge_MFReader_open(m.native, &p))
	if fh < 0 {
		return nil, fmt.Errorf("open file %v fail,code: %v", path, fh)
	}

	mf := &MediaFile{
		handle: fileHandle(fh),
		reader: m,
		path:   path,
	}
	mf.setup()
	return mf, nil
}

func (m *MediaFileReader) readFile(fh fileHandle, buf []byte) (int, error) {
	if !m.IsOpened() {
		return 0, ErrFileReaderNotOpen
	}
	count := len(buf)
	n := C.Edge_MFReader_read(m.native, C.int32_t(fh), unsafe.Pointer(&buf[0]), C.size_t(count))
	return int(n), nil
}

func (m *MediaFileReader) closeFile(fh fileHandle) error {
	if !m.IsOpened() {
		return ErrFileReaderNotOpen
	}
	ret := C.Edge_MFReader_close(m.native, C.int32_t(fh))
	return convertCCodeToError(int(ret))
}
