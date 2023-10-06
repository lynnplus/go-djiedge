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

#include "edge_liveview.h"

using namespace edge_sdk;
using namespace std;


PUBLIC_API CEdgeLiveView *Edge_LiveView_new() {
    auto obj = CreateLiveview();
    return new CEdgeLiveView({obj, nullptr});
}

PUBLIC_API void Edge_LiveView_delete(CEdgeLiveView *obj) {
    if (obj == nullptr) {
        return;
    }
    obj->ctx = nullptr;
    obj->instance = nullptr;
    delete obj;
}

PUBLIC_API int Edge_LiveView_init(CEdgeLiveView *obj, const CEdgeLiveViewOptions *opt) {
    if (obj == nullptr || obj->instance == nullptr || opt == nullptr) {
        return kErrorInvalidArgument;
    }
    Liveview::H264Callback cb = nullptr;
    if (opt->stream_callback != nullptr) {
        auto streamCB = opt->stream_callback;
        cb = [obj, streamCB](const uint8_t *buf, uint32_t len) -> ErrorCode {
            streamCB(obj->ctx, buf, len);
            return kOk;
        };
    }
    Liveview::Options options = {static_cast<Liveview::CameraType>(opt->camera),
                                 static_cast<Liveview::StreamQuality>(opt->quality), cb};
    return obj->instance->Init(options);
}

PUBLIC_API int Edge_LiveView_deInit(CEdgeLiveView *obj) {
    if (obj == nullptr || obj->instance == nullptr) {
        return kErrorInvalidArgument;
    }
    return obj->instance->DeInit();
}

PUBLIC_API int Edge_LiveView_setCameraSource(CEdgeLiveView *obj, int source) {
    if (obj == nullptr || obj->instance == nullptr) {
        return kErrorInvalidArgument;
    }
    return obj->instance->SetCameraSource(static_cast<Liveview::CameraSource>(source));
}

PUBLIC_API int Edge_LiveView_subscribeStreamStatus(CEdgeLiveView *obj, CEdgeLiveViewStreamStatusCallback callback) {
    if (obj == nullptr || obj->instance == nullptr || callback == nullptr) {
        return kErrorInvalidArgument;
    }
    auto cb = [obj, callback](const Liveview::LiveviewStatus &status) -> void {
        callback(obj->ctx, status);
    };
    return obj->instance->SubscribeLiveviewStatus(cb);
}

PUBLIC_API int Edge_LiveView_startH264Stream(CEdgeLiveView *obj) {
    if (obj == nullptr || obj->instance == nullptr) {
        return kErrorInvalidArgument;
    }
    return obj->instance->StartH264Stream();
}

PUBLIC_API int Edge_LiveView_stopH264Stream(CEdgeLiveView *obj) {
    if (obj == nullptr || obj->instance == nullptr) {
        return kErrorInvalidArgument;
    }
    return obj->instance->StopH264Stream();
}
