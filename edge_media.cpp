//go:build linux

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

PUBLIC_API CEdgeMFReader *Edge_MediaMgr_createMediaFilesReader() {
    auto p = MediaManager::Instance()->CreateMediaFilesReader();
    return new CEdgeMFReader(p);
}

PUBLIC_API void Edge_MediaMgr_deleteMediaFilesReader(CEdgeMFReader *reader) {
    delete reader;
}


void convert_to_c_media_file(const MediaFile &src, CEdgeMediaFile &dst) {
    dst.file_name = {src.file_name.size(), src.file_name.c_str()};
    dst.file_path = {src.file_path.size(), src.file_path.c_str()};
    dst.file_size = src.file_size;
    dst.file_type = src.file_type;
    dst.camera_attr = src.camera_attr;
    dst.latitude = src.latitude;
    dst.longitude = src.longitude;
    dst.absolute_altitude = src.absolute_altitude;
    dst.relative_altitude = src.relative_altitude;
    dst.gimbal_yaw_degree = src.gimbal_yaw_degree;
    dst.image_width = src.image_width;
    dst.image_height = src.image_height;
    dst.video_duration = src.video_duration;
    dst.create_time = src.create_time;
}


PUBLIC_API int Edge_MediaMgr_registerMediaFilesObserver(CEdgeMediaFilesObserver callback) {
    auto observer = [callback](const MediaFile &file) -> ErrorCode {
        CEdgeMediaFile data;
        convert_to_c_media_file(file, data);
        callback(&data);
        return kOk;
    };
    return MediaManager::Instance()->RegisterMediaFilesObserver(observer);
}

PUBLIC_API int Edge_MediaMgr_setDroneNestUploadCloud(bool enable) {
    return MediaManager::Instance()->SetDroneNestUploadCloud(enable);
}

PUBLIC_API int Edge_MediaMgr_setDroneNestAutoDelete(bool enable) {
    return MediaManager::Instance()->SetDroneNestAutoDelete(enable);
}

PUBLIC_API int Edge_MFReader_init(CEdgeMFReader *reader) {
    return (*reader)->Init();
}

PUBLIC_API int32_t Edge_MFReader_fileList(CEdgeMFReader *reader, const CEdgeMediaFile **files) {
    MediaFilesReader::MediaFileList list;
    auto num = (*reader)->FileList(list);
    if (num <= 0) {
        return num;
    }
    size_t n = 0;
    for (const auto &item: list) {
        n += item->file_name.size() + item->file_path.size();
    }
    auto out = (CEdgeMediaFile *) malloc(sizeof(CEdgeMediaFile) * num + sizeof(char) * n);
    *files = out;
    char *buf = (char *) (out + sizeof(CEdgeMediaFile) * num);
    for (const auto &item: list) {
        convert_to_c_media_file(*(item), (*out));
        
        out->file_name = copy_from_string(item->file_name, buf);
        buf += item->file_name.size();
        out->file_path = copy_from_string(item->file_path, buf);
        buf += item->file_path.size();
        
        out++;
    }
    return num;
}

PUBLIC_API int32_t Edge_MFReader_open(CEdgeMFReader *reader, const CCString *file_path) {
    return (*reader)->Open(copy_from_cstring(*file_path));
}

PUBLIC_API size_t Edge_MFReader_read(CEdgeMFReader *reader, int32_t fileHandle, void *buf, size_t count) {
    return (*reader)->Read(fileHandle, buf, count);
}

PUBLIC_API int Edge_MFReader_close(CEdgeMFReader *reader, int32_t fileHandle) {
    return (*reader)->Close(fileHandle);
}

PUBLIC_API int Edge_MFReader_deInit(CEdgeMFReader *reader) {
    return (*reader)->DeInit();
}
