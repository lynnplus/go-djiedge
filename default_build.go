//go:build linux && !custom_edge_env

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

// by default,just use the following cgo configuration,need to provide include dir and link search dir manually.
// stdc++ minimum support c++11

/*
#cgo LDFLAGS: -ledgesdk -lcrypto -lssh2
*/
import "C"
