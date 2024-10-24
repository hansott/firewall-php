#include "Includes.h"

RequestProcessor requestProcessor;

std::string RequestProcessor::GetInitData() {
    json initData = {
        {"log_level", AIKIDO_GLOBAL(log_level_str)},
        {"socket_path", AIKIDO_GLOBAL(socket_path)},
        {"trust_proxy", AIKIDO_GLOBAL(trust_proxy)},
        {"localhost_allowed_by_default", AIKIDO_GLOBAL(localhost_allowed_by_default)},
        {"collect_api_schema", AIKIDO_GLOBAL(collect_api_schema)},
        {"sapi", AIKIDO_GLOBAL(sapi_name)}};

    return initData.dump();
}

bool RequestProcessor::ContextInit() {
    return this->requestProcessorContextInitFn(GoContextCallback);
}

bool RequestProcessor::SendEvent(EVENT_ID eventId, std::string& output) {
    if (!this->requestProcessorOnEventFn) {
        return false;
    }

    AIKIDO_LOG_DEBUG("Sending event %s...\n", GetEventName(eventId));

    char* charPtr = requestProcessorOnEventFn(eventId);
    if (!charPtr) {
        AIKIDO_LOG_DEBUG("Got event reply: nullptr\n");
        return true;
    }

    AIKIDO_LOG_DEBUG("Got event reply: %s\n", charPtr);

    output = charPtr;
    free(charPtr);
    return true;
}

void RequestProcessor::SendPreRequestEvent() {
    try {
        std::string outputEvent;
        SendEvent(EVENT_PRE_REQUEST, outputEvent);
        action.Execute(outputEvent);
    } catch (const std::exception& e) {
        AIKIDO_LOG_ERROR("Exception encountered in processing request init metadata: %s\n", e.what());
    }
}

void RequestProcessor::SendPostRequestEvent() {
    try {
        std::string outputEvent;
        SendEvent(EVENT_POST_REQUEST, outputEvent);
        action.Execute(outputEvent);
    } catch (const std::exception& e) {
        AIKIDO_LOG_ERROR("Exception encountered in processing request shutdown metadata: %s\n", e.what());
    }
}

/*
    If the blocking mode is set from agent (different than -1), return that.
        Otherwise, return the env variable AIKIDO_BLOCK.
*/
bool RequestProcessor::IsBlockingEnabled() {
    if (!requestProcessorGetBlockingModeFn) {
        return false;
    }
    int ret = requestProcessorGetBlockingModeFn();
    if (ret == -1) {
        return AIKIDO_GLOBAL(blocking);
    }
    return ret;
}

bool RequestProcessor::Init() {
    if (this->initFailed) {
        return false;
    }

    if (!this->libHandle) {
        std::string requestProcessorLibPath = "/opt/aikido-" + std::string(PHP_AIKIDO_VERSION) + "/aikido-request-processor.so";
        this->libHandle = dlopen(requestProcessorLibPath.c_str(), RTLD_LAZY);
        if (!this->libHandle) {
            AIKIDO_LOG_ERROR("Error loading the Aikido Request Processor library from %s: %s!\n", requestProcessorLibPath.c_str(), dlerror());
            this->initFailed = true;
            return false;
        }

        AIKIDO_LOG_DEBUG("Initializing Aikido Request Processor...\n");

        RequestProcessorInitFn requestProcessorInitFn = (RequestProcessorInitFn)dlsym(libHandle, "RequestProcessorInit");
        this->requestProcessorContextInitFn = (RequestProcessorContextInitFn)dlsym(libHandle, "RequestProcessorContextInit");
        this->requestProcessorOnEventFn = (RequestProcessorOnEventFn)dlsym(libHandle, "RequestProcessorOnEvent");
        this->requestProcessorGetBlockingModeFn = (RequestProcessorGetBlockingModeFn)dlsym(libHandle, "RequestProcessorGetBlockingMode");
        this->requestProcessorUninitFn = (RequestProcessorUninitFn)dlsym(libHandle, "RequestProcessorUninit");
        if (!requestProcessorInitFn ||
            !this->requestProcessorContextInitFn ||
            !this->requestProcessorOnEventFn ||
            !this->requestProcessorGetBlockingModeFn ||
            !this->requestProcessorUninitFn) {
            AIKIDO_LOG_ERROR("Error loading symbols from the Aikido Request Processor library!\n");
            this->initFailed = true;
            return false;
        }

        std::string initDataString = this->GetInitData();
        if (!requestProcessorInitFn(GoCreateString(initDataString))) {
            AIKIDO_LOG_ERROR("Failed to initialize Aikido Request Processor library: %s!\n", dlerror());
            this->initFailed = true;
            return false;
        }

        AIKIDO_LOG_DEBUG("Aikido Request Processor initialized successfully!\n");
    }

    if (!request.Init()) {
        AIKIDO_LOG_WARN("Failed to initialize the current request!\n");
        return false;
    }

    ContextInit();
    SendPreRequestEvent();
    return true;
}

void RequestProcessor::Uninit() {
    if (!request.Init()) {
        AIKIDO_LOG_WARN("Failed to initialize the current request!\n");
        return;
    }
    SendPostRequestEvent();
}

RequestProcessor::~RequestProcessor() {
    if (!this->libHandle) {
        return;
    }
    if (!this->initFailed && this->requestProcessorUninitFn) {
        this->requestProcessorUninitFn();
    }
    dlclose(this->libHandle);
    this->libHandle = nullptr;
}
