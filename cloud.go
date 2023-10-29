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

#include "edge_cloud.h"

void esdkCgoCloudCustomMsgCallback(uint8_t *data, uint32_t len);
*/
import "C"
import (
	"errors"
	"unsafe"
)

//export esdkCgoCloudCustomMsgCallback
func esdkCgoCloudCustomMsgCallback(buf *C.uint8_t, size C.uint32_t) {
	if cloudCustomMsgHandler == nil {
		return
	}
	b := C.GoBytes(unsafe.Pointer(buf), C.int(size))
	cloudCustomMsgHandler(b)
}

var cloudCustomMsgHandler func([]byte)

// SendCustomMessageToCloud send custom event message to cloud,
// allows for sending data up to 256 bytes.
// when using this interface, data is internally encapsulated following the Cloud API protocol format.
func SendCustomMessageToCloud(data []byte) error {
	if !Initialized() {
		return ErrSDKNotInit
	}
	size := len(data)
	if size > 256 {
		return errors.New("data size exceeds 256 bytes")
	}
	ret := C.Edge_Cloud_sendCustomEventsMessage((*C.uint8_t)(unsafe.SliceData(data)), C.uint32_t(size))
	return convertCCodeToError(int(ret))
}

// RegisterCloudCustomMsgHandler register a callback function via this interface to manage incoming data from the cloud.
// Ensure that the data is handled asynchronously within the callback function.
// Synchronous handling might block message reception. Instead,
// move data to a buffer queue for asynchronous processing.
func RegisterCloudCustomMsgHandler(handler func([]byte)) error {
	ret := C.Edge_Cloud_registerCustomMsgHandler(C.CEdgeCloudCustomMsgHandler(C.esdkCgoCloudCustomMsgCallback))
	if err := convertCCodeToError(int(ret)); err != nil {
		return err
	}
	cloudCustomMsgHandler = handler
	return nil
}
