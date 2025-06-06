/*
 * Copyright 2023 iLogtail Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include <string>

#include "json/json.h"

namespace logtail {

struct ConfigContainerInfoUpdateCmd {
    std::string mConfigName; // config name
    bool mDeleteFlag = false; // if this flag is true, delete the container from this config's ContainerInfo array
    Json::Value mJsonParams; // jsonParams, json string.

    ConfigContainerInfoUpdateCmd(const std::string& configName, bool delFlag, const Json::Value& jsonParams)
        : mConfigName(configName), mDeleteFlag(delFlag), mJsonParams(jsonParams) {}
};

} // namespace logtail
