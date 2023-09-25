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

#include "edge_media.h"

using namespace edge_sdk;

CEdgeMFReader *Edge_MediaMgr_createMediaFilesReader() {
    auto p = MediaManager::Instance()->CreateMediaFilesReader();
    return new CEdgeMFReader(p);
}

void Edge_MediaMgr_deleteMediaFilesReader(CEdgeMFReader *reader) {
    delete reader;
}

int Edge_MediaMgr_registerMediaFilesObserver(CEdgeMediaFilesObserver callback) {
    auto observer = [callback](const MediaFile &) -> ErrorCode {
        callback(nullptr);
        return kOk;
    };
    return MediaManager::Instance()->RegisterMediaFilesObserver(observer);
}

int Edge_MediaMgr_setDroneNestUploadCloud(bool enable) {
    return MediaManager::Instance()->SetDroneNestUploadCloud(enable);
}

int Edge_MediaMgr_setDroneNestAutoDelete(bool enable) {
    return MediaManager::Instance()->SetDroneNestAutoDelete(enable);
}

int Edge_MFReader_init(CEdgeMFReader *reader) {
    return (*reader)->Init();
}

int32_t Edge_MFReader_fileList(CEdgeMFReader *reader, const CEdgeMediaFile **files) {
    MediaFilesReader::MediaFileList list;
    auto num = (*reader)->FileList(list);
    if (num <= 0) {
        return num;
    }
    
    auto out = (CEdgeMediaFile *) malloc(sizeof(CEdgeMediaFile) * num);
    *files = out;
    for (const auto &item: list) {
        out++;
    }
    return num;
}

int32_t Edge_MFReader_open(CEdgeMFReader *reader, const CCString file_path) {
    return (*reader)->Open(to_cstring(file_path));
}

size_t Edge_MFReader_read(CEdgeMFReader *reader, int32_t fileHandle, void *buf, size_t count) {
    return (*reader)->Read(fileHandle, buf, count);
}

int Edge_MFReader_close(CEdgeMFReader *reader, int32_t fileHandle) {
    return (*reader)->Close(fileHandle);
}

int Edge_MFReader_deInit(CEdgeMFReader *reader) {
    return (*reader)->DeInit();
}