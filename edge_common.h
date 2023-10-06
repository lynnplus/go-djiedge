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

#ifndef CEDGE_EDGE_COMMON_H
#define CEDGE_EDGE_COMMON_H


#define PUBLIC_API __attribute__((unused))


#ifdef __cplusplus
extern "C" {
#endif
#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

typedef struct {
    size_t len;
    const char *data;
} CCString;

#ifdef __cplusplus
}
#endif


#ifdef __cplusplus

#include <string>
#include <memory>
#include <cstring>

inline std::string copy_from_cstring(const CCString &src) {
    return {src.data, src.len};
}

inline CCString copy_from_string(const std::string &src, char *buf) {
    if (src.empty()) {
        return {0, nullptr};
    }
    strncpy(buf, src.c_str(), src.size());
    return {src.size(), buf};
}


//wrap cxx shared_ptr
template<typename T>
class CWrapPtr {
public:
    __attribute__((unused))
    explicit CWrapPtr(const std::shared_ptr<T> _ptr) : ptr(_ptr) {}
    
    CWrapPtr(const CWrapPtr &p) : ptr(p.ptr) {}
    
    ~CWrapPtr() {
        this->ptr = nullptr;
    }
    
    T *operator->() const {
        return ptr.get();
    }

private:
    std::shared_ptr<T> ptr;
};

#endif //__cplusplus

#endif //CEDGE_EDGE_COMMON_H
