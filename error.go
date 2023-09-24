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
	"errors"
	"fmt"
)

var (
	ErrInvalidArgument   = errors.New("djiedge: invalid argument")
	ErrSystemError       = errors.New("djiedge: invalid use of UnreadRune")
	ErrInvalidOperation  = errors.New("djiedge: invalid operation")
	ErrRepeatOperation   = errors.New("djiedge: repeated operation")
	ErrNullPointer       = errors.New("djiedge: null pointer")
	ErrParamOutOfRange   = errors.New("djiedge: parameter has exceeded the expected range")
	ErrParamGetFailure   = errors.New("djiedge: failed to get a parameter")
	ErrParamSetFailure   = errors.New("djiedge: failed to set or modify a parameter")
	ErrSendPackFailure   = errors.New("djiedge: failed to send pack")
	ErrRequestTimeout    = errors.New("djiedge: request has timed out")
	ErrAuthVerifyFailure = errors.New("djiedge: a failure in verifying the authorization information")
	ErrEncryptFailure    = errors.New("djiedge: failed to encrypt data")
	ErrDecryptFailure    = errors.New("djiedge: failed to decrypt data")
	ErrInvalidRespond    = errors.New("djiedge: invalid respond")
	ErrRemoteFailure     = errors.New("djiedge: a failure on the remote server or remote process")
	ErrNoVideoID         = errors.New("djiedge: failed to get a valid video ID while starting a live stream")
	ErrConnectFailure    = errors.New("djiedge: a failure in establishing a connection")
)

func convertCCodeToError(code int) error {
	if code == 0 {
		return nil
	}
	switch code {
	case 1:
		return ErrInvalidArgument
	case 2:
		return ErrSystemError
	case 3:
		return ErrInvalidOperation
	case 4:
		return ErrRepeatOperation
	case 5:
		return ErrNullPointer
	case 6:
		return ErrParamOutOfRange
	case 7:
		return ErrParamGetFailure
	case 8:
		return ErrParamSetFailure
	case 9:
		return ErrSendPackFailure
	case 10:
		return ErrRequestTimeout
	case 11:
		return ErrAuthVerifyFailure
	case 12:
		return ErrEncryptFailure
	case 13:
		return ErrDecryptFailure
	case 14:
		return ErrInvalidRespond
	case 15:
		return ErrRemoteFailure
	case 16:
		return ErrNoVideoID
	case 17:
		return ErrConnectFailure
	}
	return fmt.Errorf("djiedge: unknown error,code %v", code)
}
