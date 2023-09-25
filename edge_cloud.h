// Copyright (c) 2023 Lynn <lynnplus90@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef CEDGE_EDGE_CLOUD_H
#define CEDGE_EDGE_CLOUD_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

typedef void (*CEdgeCloudCustomMsgHandler)(const uint8_t *data, uint32_t len);

int Edge_Cloud_registerCustomMsgHandler(CEdgeCloudCustomMsgHandler handler);

int Edge_Cloud_sendCustomEventsMessage(const uint8_t *data, uint32_t len);


#ifdef __cplusplus
}
#endif


#endif //CEDGE_EDGE_CLOUD_H
