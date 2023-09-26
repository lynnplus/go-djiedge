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

#ifndef CEDGE_EDGE_INIT_H
#define CEDGE_EDGE_INIT_H

#include "edge_common.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    uint8_t major_version;
    uint8_t minor_version;
    uint8_t modify_version;
    uint8_t debug_version;
} CEdgeVersion;

typedef struct {
    CCString app_name;
    CCString app_id;
    CCString app_key;
    CCString app_license;
    CCString developer_account;
} CEdgeAppInfo;

typedef struct {
    CCString private_key;
    CCString public_key;
} CEdgeKeyStore;

typedef void (*CEdgeLogOutput)(const uint8_t *data, uint32_t dataLen);

typedef struct {
    int level;
    bool is_support_color;
    CEdgeLogOutput output;
} CEdgeLogger;

typedef struct {
    CCString product_name;
    CCString vendor_name;
    CCString serial_number;
    CEdgeVersion firmware_version;
    CEdgeAppInfo app_info;
    CEdgeKeyStore key_store;
    CEdgeLogger logger;
} CEdgeInitOptions;

int Edge_init(const CEdgeInitOptions *opts, bool deInitOnFailed);

int Edge_deInit();

#ifdef __cplusplus
}
#endif

#endif //CEDGE_EDGE_INIT_H
