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

package main

import (
	"fmt"
	edge "github.com/lynnplus/go-djiedge"
	"os"
	"os/signal"
	"syscall"
)

type BytesData struct {
}

func (b *BytesData) Write(p []byte) (n int, err error) {
	return 0, err
}

func (b *BytesData) Test(p []byte) (n int, err error) {
	return 0, err
}

type streamHandler struct {
}

func (s *streamHandler) OnStreamStatusUpdate(status *edge.LiveStatus) {
	//TODO implement me
	panic("implement me")
}

func (s *streamHandler) OnReceiveStreamData(data []byte) {
	//TODO implement me
	panic("implement me")
}

func (s *streamHandler) AllocateStreamData() *BytesData {
	//TODO implement me
	panic("implement me")
}

func main() {

	s := &streamHandler{}

	lv := edge.NewLiveView[*BytesData]()
	err := lv.Init(6, 6, s)

	fmt.Println(err)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
	<-c
}

func initSDK() error {
	app := &edge.AuthInfo{
		Name:    "",
		Id:      "",
		AppKey:  "",
		License: "",
		Account: "",
	}

	dev := &edge.DeviceInfo{
		ProductName:     "",
		VendorName:      "",
		SerialNumber:    "sn",
		FirmwareVersion: edge.FirmwareVersion{MinorVersion: 1},
	}

	key := &edge.RSA2048Key{
		PrivateKey: "",
		PublicKey:  "",
	}

	logger := &edge.Logger{
		Level:          edge.LogLevelDebug,
		EnableColorful: true,
		Outputer: func(msg string) {
			fmt.Println(msg)
		},
	}
	return edge.InitSDK(dev, app, key, logger, true)
}
