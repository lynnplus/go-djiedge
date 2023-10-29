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

#include "edge_init.h"

void esdkCallGoLogger(uint8_t *data, uint32_t dataLen);

*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var sdkLogHandler func(string)

//export esdkCallGoLogger
func esdkCallGoLogger(data *C.uint8_t, dataLen C.uint32_t) {
	if sdkLogHandler == nil {
		return
	}
	//int(dataLen)-2 remove the last newline character(\r\n)
	tmp := unsafe.Slice((*byte)(data), int(dataLen)-2)
	str := string(tmp)
	sdkLogHandler(str)
}

func convertToCVersion(version *FirmwareVersion) C.CEdgeVersion {
	return C.CEdgeVersion{
		major_version:  C.uint8_t(version.MajorVersion),
		minor_version:  C.uint8_t(version.MinorVersion),
		modify_version: C.uint8_t(version.ModifyVersion),
		debug_version:  C.uint8_t(version.DebugVersion),
	}
}

var initState atomic.Int32 //0:none  1: initializing  2:initialized

func Initialized() bool {
	return initState.Load() == 2
}

// InitSDK initialize edge-sdk
// deInitOnFailed: de-initialize the sdk after failure,because DJI will have some threads continuing to run after initialization failure.
func InitSDK(device *DeviceInfo, auth *AuthInfo, key *RSA2048Key, logger *Logger, deInitOnFailed bool) (err error) {
	if !initState.CompareAndSwap(0, 1) {
		return errors.New("sdk is initializing or initialized")
	}
	defer func() {
		if err != nil {
			initState.Store(0)
		} else {
			initState.Store(2)
		}
	}()

	if device == nil || auth == nil || key == nil {
		return errors.New("parameter is nil")
	}
	// dji bug,an exception occurs when sn is empty
	if device.SerialNumber == "" {
		return errors.New("parameter is nil of device sn")
	}

	appInfo := C.CEdgeAppInfo{
		app_name:          convertToCString(auth.Name),
		app_id:            convertToCString(auth.Id),
		app_key:           convertToCString(auth.AppKey),
		app_license:       convertToCString(auth.License),
		developer_account: convertToCString(auth.Account),
	}
	ks := C.CEdgeKeyStore{
		private_key: convertToCString(key.PrivateKey),
		public_key:  convertToCString(key.PublicKey),
	}

	var logs C.CEdgeLogger
	if logger != nil {
		if !logger.Level.IsValid() {
			return fmt.Errorf("%v is not a valid log level", logger.Level)
		}
		logs = C.CEdgeLogger{
			level:            C.int(logger.Level),
			is_support_color: C.bool(logger.EnableColorful),
			output:           C.CEdgeLogOutput(C.esdkCallGoLogger),
		}
		sdkLogHandler = logger.Outputer
	}

	opts := C.CEdgeInitOptions{
		product_name:     convertToCString(device.ProductName),
		vendor_name:      convertToCString(device.VendorName),
		serial_number:    convertToCString(device.SerialNumber),
		firmware_version: convertToCVersion(&device.FirmwareVersion),
		app_info:         appInfo,
		key_store:        ks,
		logger:           logs,
	}

	ret := C.Edge_init(opts, C.bool(deInitOnFailed))
	runtime.KeepAlive(device)
	runtime.KeepAlive(auth)
	runtime.KeepAlive(key)
	return convertCCodeToError(int(ret))
}

// DeInitSDK will de-initialize SDK environment
func DeInitSDK() error {
	if !initState.CompareAndSwap(1, 0) {
		return errors.New("sdk is not initialized")
	}
	ret := C.Edge_deInit()
	return convertCCodeToError(int(ret))
}
