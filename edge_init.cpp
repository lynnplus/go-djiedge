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

#include "edge_init.h"

#include <init.h>

using namespace edge_sdk;
using namespace std;


class EdgeKeyStoreImpl : public KeyStore {
public:
    
    EdgeKeyStoreImpl(string prv_key, string pub_key) : pPrivateKey(std::move(prv_key)),
                                                       pPublicKey(std::move(pub_key)) {}
    
    ErrorCode RSA2048_GetDERPrivateKey(string &private_key) const override {
        private_key = pPrivateKey;
        return kOk;
    }
    
    ErrorCode RSA2048_GetDERPublicKey(string &public_key) const override {
        public_key = pPublicKey;
        return kOk;
    }

private:
    const string pPrivateKey;
    const string pPublicKey;
};

int Edge_init(const CEdgeInitOptions *opts, bool deInitOnFailed) {
    if (opts == nullptr) {
        return kErrorInvalidArgument;
    }
    Options option;
    option.product_name = copy_from_cstring(opts->product_name);
    option.vendor_name = copy_from_cstring(opts->vendor_name);
    option.serial_number = copy_from_cstring(opts->serial_number);
    option.firmware_version = {opts->firmware_version.major_version, opts->firmware_version.minor_version,
                               opts->firmware_version.modify_version, opts->firmware_version.debug_version};
    
    AppInfo app_info;
    app_info.app_name = copy_from_cstring(opts->app_info.app_name);
    app_info.app_id = copy_from_cstring(opts->app_info.app_id);
    app_info.app_key = copy_from_cstring(opts->app_info.app_key);
    app_info.app_license = copy_from_cstring(opts->app_info.app_license);
    app_info.developer_account = copy_from_cstring(opts->app_info.developer_account);
    
    option.app_info = app_info;
    
    auto logger = opts->logger;
    if (logger.level >= 0 && logger.output) {
        auto outputFunc = logger.output;
        auto ff = [outputFunc](const uint8_t *data, uint32_t dataLen) -> ErrorCode {
            outputFunc(data, dataLen);
            return kOk;
        };
        LoggerConsole console = {static_cast<LogLevel>(logger.level), ff, logger.is_support_color};
        option.logger_console_lists.push_back(console);
    }
    
    auto pubKey = copy_from_cstring(opts->key_store.public_key);
    auto prvKey = copy_from_cstring(opts->key_store.private_key);
    option.key_store = std::make_shared<EdgeKeyStoreImpl>(prvKey, pubKey);
    
    auto ret = ESDKInit::Instance()->Init(option);
    if (ret != kOk && deInitOnFailed) {
        Edge_deInit();
    }
    return ret;
}


int Edge_deInit() {
    return ESDKInit::Instance()->DeInit();
}
