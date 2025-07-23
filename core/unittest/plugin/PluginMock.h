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

#include "collection_pipeline/CollectionPipeline.h"
#include "collection_pipeline/plugin/PluginRegistry.h"
#include "collection_pipeline/plugin/creator/StaticFlusherCreator.h"
#include "collection_pipeline/plugin/creator/StaticInputCreator.h"
#include "collection_pipeline/plugin/creator/StaticProcessorCreator.h"
#include "collection_pipeline/plugin/interface/Flusher.h"
#include "collection_pipeline/plugin/interface/HttpFlusher.h"
#include "collection_pipeline/plugin/interface/Input.h"
#include "collection_pipeline/plugin/interface/Processor.h"
#include "collection_pipeline/queue/SLSSenderQueueItem.h"
#include "collection_pipeline/queue/SenderQueueManager.h"
#include "common/StringView.h"
#include "plugin/flusher/sls/FlusherSLS.h"
#include "task_pipeline/Task.h"
#include "task_pipeline/TaskRegistry.h"

namespace logtail {

class ProcessorInnerMock : public Processor {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config) override { return true; }
    void Process(PipelineEventGroup& logGroup) override { ++mCnt; };

    uint32_t mCnt = 0;

protected:
    bool IsSupportedEvent(const PipelineEventPtr& e) const override { return true; };
};

const std::string ProcessorInnerMock::sName = "processor_inner_mock";

class InputMock : public Input {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config, Json::Value& optionalGoPipeline) override {
        if (config.isMember("SupportAck")) {
            mSupportAck = config["SupportAck"].asBool();
        }
        auto processor = PluginRegistry::GetInstance()->CreateProcessor(
            ProcessorInnerMock::sName, mContext->GetPipeline().GenNextPluginMeta(false));
        processor->Init(Json::Value(), *mContext);
        mInnerProcessors.emplace_back(std::move(processor));
        return true;
    }
    bool Start() override { return true; }
    bool Stop(bool isPipelineRemoving) override {
        while (mBlockFlag) {
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
        return true;
    }
    bool SupportAck() const override { return mSupportAck; }

    void Block() { mBlockFlag = true; }
    void Unblock() { mBlockFlag = false; }

    bool mSupportAck = true;

private:
    std::atomic_bool mBlockFlag = false;
};

const std::string InputMock::sName = "input_mock";

class InputSingletonMock1 : public Input {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config, Json::Value& optionalGoPipeline) override {
        if (config.isMember("SupportAck")) {
            mSupportAck = config["SupportAck"].asBool();
        }
        auto processor = PluginRegistry::GetInstance()->CreateProcessor(
            ProcessorInnerMock::sName, mContext->GetPipeline().GenNextPluginMeta(false));
        processor->Init(Json::Value(), *mContext);
        mInnerProcessors.emplace_back(std::move(processor));
        return true;
    }
    bool Start() override { return true; }
    bool Stop(bool isPipelineRemoving) override {
        while (mBlockFlag) {
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
        return true;
    }
    bool SupportAck() const override { return mSupportAck; }

    void Block() { mBlockFlag = true; }
    void Unblock() { mBlockFlag = false; }

    bool mSupportAck = true;

private:
    std::atomic_bool mBlockFlag = false;
};

const std::string InputSingletonMock1::sName = "input_singleton_mock_1";

class InputSingletonMock2 : public Input {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config, Json::Value& optionalGoPipeline) override {
        if (config.isMember("SupportAck")) {
            mSupportAck = config["SupportAck"].asBool();
        }
        auto processor = PluginRegistry::GetInstance()->CreateProcessor(
            ProcessorInnerMock::sName, mContext->GetPipeline().GenNextPluginMeta(false));
        processor->Init(Json::Value(), *mContext);
        mInnerProcessors.emplace_back(std::move(processor));
        return true;
    }
    bool Start() override { return true; }
    bool Stop(bool isPipelineRemoving) override {
        while (mBlockFlag) {
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
        return true;
    }
    bool SupportAck() const override { return mSupportAck; }

    void Block() { mBlockFlag = true; }
    void Unblock() { mBlockFlag = false; }

    bool mSupportAck = true;

private:
    std::atomic_bool mBlockFlag = false;
};

const std::string InputSingletonMock2::sName = "input_singleton_mock_2";

const std::string PROCESSOR_MOCK_LOCAL_CONTENT_KEY = "processor_mock_local_content_key";
const std::string PROCESSOR_MOCK_LOCAL_CONTENT_VALUE = "processor_mock_local_content_value";

class ProcessorMock : public Processor {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config) override {
        mLocalContentKey = PROCESSOR_MOCK_LOCAL_CONTENT_KEY;
        mLocalContentValue = PROCESSOR_MOCK_LOCAL_CONTENT_VALUE;
        return true;
    }
    void Process(PipelineEventGroup& logGroup) override {
        for (auto& e : logGroup.MutableEvents()) {
            if (e.Is<LogEvent>()) {
                auto& logEvent = e.Cast<LogEvent>();
                logEvent.SetContentNoCopy(StringView(mLocalContentKey.data(), mLocalContentKey.size()),
                                          StringView(mLocalContentValue.data(), mLocalContentValue.size()));
            }
        }
        while (mBlockFlag) {
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
        ++mCnt;
    };

    void Block() { mBlockFlag = true; }
    void Unblock() { mBlockFlag = false; }

    uint32_t mCnt = 0;

protected:
    bool IsSupportedEvent(const PipelineEventPtr& e) const override { return true; };

    std::atomic_bool mBlockFlag = false;
    std::string mLocalContentKey;
    std::string mLocalContentValue;
};

const std::string ProcessorMock::sName = "processor_mock";

class FlusherMock : public Flusher {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config, Json::Value& optionalGoPipeline) override {
        GenerateQueueKey("mock");
        SenderQueueManager::GetInstance()->CreateQueue(mQueueKey, mPluginID, *mContext);
        return true;
    }
    bool Send(PipelineEventGroup&& g) override { return mIsValid; }
    bool Flush(size_t key) override {
        mFlushedQueues.push_back(key);
        return true;
    }
    bool FlushAll() override { return mIsValid; }

    bool mIsValid = true;
    std::vector<size_t> mFlushedQueues;
};

const std::string FlusherMock::sName = "flusher_mock";

class FlusherHttpMock : public HttpFlusher {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config, Json::Value& optionalGoPipeline) override {
        GenerateQueueKey("mock");
        SenderQueueManager::GetInstance()->CreateQueue(mQueueKey, mPluginID, *mContext);
        return true;
    }
    bool Send(PipelineEventGroup&& g) override { return mIsValid; }
    bool Flush(size_t key) override {
        mFlushedQueues.push_back(key);
        return true;
    }
    bool FlushAll() override { return mIsValid; }
    bool BuildRequest(SenderQueueItem* item,
                      std::unique_ptr<HttpSinkRequest>& req,
                      bool* keepItem,
                      std::string* errMsg) override {
        if (item->mData == "invalid_keep") {
            *keepItem = true;
            return false;
        }
        if (item->mData == "invalid_discard") {
            *keepItem = false;
            return false;
        }
        req = std::make_unique<HttpSinkRequest>(
            "", false, "", 80, "", "", std::map<std::string, std::string>(), "", nullptr);
        return true;
    }
    void OnSendDone(const HttpResponse& response, SenderQueueItem* item) override {}

    bool mIsValid = true;
    std::vector<size_t> mFlushedQueues;
};

const std::string FlusherHttpMock::sName = "flusher_http_mock";

class TaskMock : public Task {
public:
    static const std::string sName;

    const std::string& Name() const override { return sName; }
    bool Init(const Json::Value& config) override {
        if (config.isMember("Valid")) {
            return config["Valid"].asBool();
        }
        return true;
    }
    void Start() override { mIsRunning = true; }
    void Stop(bool isRemoving) { mIsRunning = false; }

    bool mIsRunning = false;
};

const std::string TaskMock::sName = "task_mock";

void LoadPluginMock() {
    PluginRegistry::GetInstance()->RegisterInputCreator(new StaticInputCreator<InputMock>());
    PluginRegistry::GetInstance()->RegisterInputCreator(new StaticInputCreator<InputSingletonMock1>(), true);
    PluginRegistry::GetInstance()->RegisterInputCreator(new StaticInputCreator<InputSingletonMock2>(), true);
    PluginRegistry::GetInstance()->RegisterProcessorCreator(new StaticProcessorCreator<ProcessorInnerMock>());
    PluginRegistry::GetInstance()->RegisterProcessorCreator(new StaticProcessorCreator<ProcessorMock>());
    PluginRegistry::GetInstance()->RegisterFlusherCreator(new StaticFlusherCreator<FlusherMock>());
    PluginRegistry::GetInstance()->RegisterFlusherCreator(new StaticFlusherCreator<FlusherHttpMock>());
}

void LoadTaskMock() {
    TaskRegistry::GetInstance()->RegisterCreator(TaskMock::sName, []() { return std::make_unique<TaskMock>(); });
}

} // namespace logtail
