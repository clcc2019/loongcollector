/*
 * Copyright 2023 iLogtail Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "plugin/processor/CommonParserOptions.h"

#include "common/ParamExtractor.h"
#include "constants/Constants.h"
#include "plugin/processor/inner/ProcessorParseContainerLogNative.h"

using namespace std;

namespace logtail {

const string CommonParserOptions::legacyUnmatchedRawLogKey = "__raw_log__";

bool CommonParserOptions::Init(const Json::Value& config,
                               const CollectionPipelineContext& ctx,
                               const string& pluginType) {
    string errorMsg;

    // KeepingSourceWhenParseFail
    if (!GetOptionalBoolParam(config, "KeepingSourceWhenParseFail", mKeepingSourceWhenParseFail, errorMsg)) {
        PARAM_WARNING_DEFAULT(ctx.GetLogger(),
                              ctx.GetAlarm(),
                              errorMsg,
                              mKeepingSourceWhenParseFail,
                              pluginType,
                              ctx.GetConfigName(),
                              ctx.GetProjectName(),
                              ctx.GetLogstoreName(),
                              ctx.GetRegion());
    }

    // KeepingSourceWhenParseSucceed
    if (!GetOptionalBoolParam(config, "KeepingSourceWhenParseSucceed", mKeepingSourceWhenParseSucceed, errorMsg)) {
        PARAM_WARNING_DEFAULT(ctx.GetLogger(),
                              ctx.GetAlarm(),
                              errorMsg,
                              mKeepingSourceWhenParseSucceed,
                              pluginType,
                              ctx.GetConfigName(),
                              ctx.GetProjectName(),
                              ctx.GetLogstoreName(),
                              ctx.GetRegion());
    }

    // RenamedSourceKey
    if (!GetOptionalStringParam(config, "RenamedSourceKey", mRenamedSourceKey, errorMsg)) {
        PARAM_WARNING_IGNORE(ctx.GetLogger(),
                             ctx.GetAlarm(),
                             errorMsg,
                             pluginType,
                             ctx.GetConfigName(),
                             ctx.GetProjectName(),
                             ctx.GetLogstoreName(),
                             ctx.GetRegion());
    }
    if (mRenamedSourceKey.empty()) {
        // SourceKey is guranteed to exist in config
        mRenamedSourceKey = config["SourceKey"].asString();
    }

    // CopingRawLog
    if (!GetOptionalBoolParam(config, "CopingRawLog", mCopingRawLog, errorMsg)) {
        PARAM_WARNING_DEFAULT(ctx.GetLogger(),
                              ctx.GetAlarm(),
                              errorMsg,
                              mCopingRawLog,
                              pluginType,
                              ctx.GetConfigName(),
                              ctx.GetProjectName(),
                              ctx.GetLogstoreName(),
                              ctx.GetRegion());
    }

    return true;
}

bool CommonParserOptions::ShouldAddLegacyUnmatchedRawLog(bool parseSuccess) {
    return !parseSuccess && mKeepingSourceWhenParseFail && mCopingRawLog;
}

bool CommonParserOptions::ShouldAddSourceContent(bool parseSuccess) {
    return (((parseSuccess && mKeepingSourceWhenParseSucceed) || (!parseSuccess && mKeepingSourceWhenParseFail)));
}

bool CommonParserOptions::ShouldEraseEvent(bool parseSuccess,
                                           const LogEvent& sourceEvent,
                                           const GroupMetadata& metadata) {
    if (!parseSuccess && !mKeepingSourceWhenParseFail) {
        if (sourceEvent.Empty()) {
            return true;
        }
        size_t size = sourceEvent.Size();
        // "__file_offset__"
        auto offsetKey = metadata.find(EventGroupMetaKey::LOG_FILE_OFFSET_KEY);
        if (size == 1 && (offsetKey != metadata.end() && sourceEvent.cbegin()->first == offsetKey->second)) {
            return true;
        } else if (size == 2 && sourceEvent.HasContent(ProcessorParseContainerLogNative::containerTimeKey)
                   && sourceEvent.HasContent(ProcessorParseContainerLogNative::containerSourceKey)) {
            return true;
        }
    }
    return false;
}

} // namespace logtail
