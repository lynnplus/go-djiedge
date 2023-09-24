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
	"fmt"
	"unsafe"
)

type LogLevel int

func (l LogLevel) IsValid() bool {
	return l >= LogLevelError && l <= LogLevelDebug
}

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

func cstring(str string) C.CCString {
	p := unsafe.Pointer(unsafe.StringData(str))
	return C.CCString{
		data: (*C.char)(p),
		len:  C.size_t(len(str)),
	}
}

//export esdkCallGoLogger
func esdkCallGoLogger(data *C.uint8_t, dataLen C.uint32_t) {
	//int(dataLen)-2 remove the last newline character(\r\n)
	tmp := unsafe.Slice(data, int(dataLen)-2)
	str := string(tmp)
	fmt.Println(str)
}

func getVersion() C.CEdgeVersion {
	return C.CEdgeVersion{
		major_version:  C.uint8_t(0),
		minor_version:  C.uint8_t(1),
		modify_version: C.uint8_t(0),
		debug_version:  C.uint8_t(0),
	}
}

func InitSDK() error {
	device := &C.CEdgeDevice{
		product_name:     cstring("product_name"),
		vendor_name:      cstring("vendor_name"),
		serial_number:    cstring("serial_number"),
		firmware_version: getVersion(),
	}

	app := &C.CEdgeAppInfo{
		app_name:          cstring("app_name"),
		app_id:            cstring("app_id"),
		app_key:           cstring("app_key"),
		app_license:       cstring("app_license"),
		developer_account: cstring("developer_account"),
	}
	key := &C.CEdgeKeyStore{
		private_key: cstring("private_key"),
		public_key:  cstring("public_key"),
	}

	logger := &C.CEdgeLogger{
		level:            3,
		is_support_color: true,
		output:           C.CEdgeLogOutput(C.esdkCallGoLogger),
	}

	ret := C.Edge_init(device, app, key, logger, true)
	return convertCCodeToError(int(ret))
}
