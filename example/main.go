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

func main() {

	live := edge.NewLiveView()

	err := live.Init(5, 0, func(data any, size uint) {

	})

	fmt.Println(err)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
	x := <-c

	fmt.Println("\nend", x)
}
