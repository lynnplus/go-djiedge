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

#include "edge_cloud.h"

#include <cloud_api.h>

using namespace edge_sdk;

PUBLIC_API int Edge_Cloud_registerCustomMsgHandler(CEdgeCloudCustomMsgHandler handler) {
    return CloudAPI_RegisterCustomServicesMessageHandler(handler);
}

PUBLIC_API int Edge_Cloud_sendCustomEventsMessage(const uint8_t *data, uint32_t len) {
    return CloudAPI_SendCustomEventsMessage(data, len);
}
