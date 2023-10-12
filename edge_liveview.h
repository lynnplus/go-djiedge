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

#ifndef CEDGE_EDGE_LIVEVIEW_H
#define CEDGE_EDGE_LIVEVIEW_H

#include "edge_common.h"

#ifdef __cplusplus

#include <liveview.h>


typedef std::shared_ptr<edge_sdk::Liveview> EdgeLiveView;
extern "C" {
#else
typedef void* EdgeLiveView;
#endif

#include <stddef.h>
#include <stdint.h>

typedef struct {
    const void *ctx;
    EdgeLiveView instance;
} CEdgeLiveView;

typedef void (*CEdgeLiveViewStreamCallback)(const void *ctx, const uint8_t *buf, uint32_t len);
typedef void (*CEdgeLiveViewStreamStatusCallback)(const void *ctx, uint32_t status);

typedef struct {
    int camera;
    int quality;
    CEdgeLiveViewStreamCallback stream_callback;
} CEdgeLiveViewOptions;

PUBLIC_API CEdgeLiveView *Edge_LiveView_new(const void *ctx);
PUBLIC_API void Edge_LiveView_delete(CEdgeLiveView *obj);
PUBLIC_API int Edge_LiveView_init(CEdgeLiveView *obj, const CEdgeLiveViewOptions *opt);
PUBLIC_API int Edge_LiveView_deInit(CEdgeLiveView *obj);
PUBLIC_API int Edge_LiveView_setCameraSource(CEdgeLiveView *obj, int source);
PUBLIC_API int Edge_LiveView_subscribeStreamStatus(CEdgeLiveView *obj, CEdgeLiveViewStreamStatusCallback callback);
PUBLIC_API int Edge_LiveView_startH264Stream(CEdgeLiveView *obj);
PUBLIC_API int Edge_LiveView_stopH264Stream(CEdgeLiveView *obj);

#ifdef __cplusplus
}
#endif

// Liveview

#endif // CEDGE_EDGE_LIVEVIEW_H
