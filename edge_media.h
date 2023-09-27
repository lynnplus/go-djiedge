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

#ifndef CEDGE_EDGE_MEDIA_H
#define CEDGE_EDGE_MEDIA_H


#include "edge_common.h"

#ifdef __cplusplus

#include <media_manager.h>

using CEdgeMFReader = CWrapPtr<edge_sdk::MediaFilesReader>;

extern "C" {

#else
typedef void CEdgeMFReader;
#endif

#include <stdlib.h>

typedef struct {
    CCString file_name;
    CCString file_path;
    size_t file_size;
    int file_type;
    int camera_attr;
    double latitude;
    double longitude;
    double absolute_altitude;
    double relative_altitude;
    double gimbal_yaw_degree;
    int32_t image_width;
    int32_t image_height;
    uint32_t video_duration;
    time_t create_time;
} CEdgeMediaFile;

typedef void (*CEdgeMediaFilesObserver)(const CEdgeMediaFile *);


CEdgeMFReader *Edge_MediaMgr_createMediaFilesReader();

void Edge_MediaMgr_deleteMediaFilesReader(CEdgeMFReader *reader);

int Edge_MediaMgr_registerMediaFilesObserver(CEdgeMediaFilesObserver callback);

int Edge_MediaMgr_setDroneNestUploadCloud(bool enable);

int Edge_MediaMgr_setDroneNestAutoDelete(bool enable);


//***********MediaFilesReader start***********//

int Edge_MFReader_init(CEdgeMFReader *reader);

// need to manually use free to release the memory of the files on return num>0.
int32_t Edge_MFReader_fileList(CEdgeMFReader *reader, const CEdgeMediaFile **files);

int32_t Edge_MFReader_open(CEdgeMFReader *reader, const CCString *file_path);

size_t Edge_MFReader_read(CEdgeMFReader *reader, int32_t fileHandle, void *buf, size_t count);

int Edge_MFReader_close(CEdgeMFReader *reader, int32_t fileHandle);

int Edge_MFReader_deInit(CEdgeMFReader *reader);

//***********MediaFilesReader end***********//


#ifdef __cplusplus
}
#endif


#endif //CEDGE_EDGE_MEDIA_H
