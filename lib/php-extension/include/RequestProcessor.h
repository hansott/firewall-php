#pragma once

typedef GoUint8 (*RequestProcessorInitFn)(GoString initJson);
typedef GoUint8 (*RequestProcessorContextInitFn)(ContextCallback);
typedef char* (*RequestProcessorOnEventFn)(GoInt eventId);
typedef int (*RequestProcessorGetBlockingModeFn)();
typedef void (*RequestProcessorUninitFn)();

class RequestProcessor {
   private:
    bool initFailed = false;
    bool requestInitialized = false;
    void* libHandle = nullptr;
    RequestProcessorContextInitFn requestProcessorContextInitFn = nullptr;
    RequestProcessorOnEventFn requestProcessorOnEventFn = nullptr;
    RequestProcessorGetBlockingModeFn requestProcessorGetBlockingModeFn = nullptr;
    RequestProcessorUninitFn requestProcessorUninitFn = nullptr;

   private:
    std::string GetInitData();
    bool ContextInit();
    void SendPreRequestEvent();
    void SendPostRequestEvent();

   public:
    RequestProcessor() = default;

    bool Init();
    bool RequestInit();
    bool SendEvent(EVENT_ID eventId, std::string& output);
    bool IsBlockingEnabled();
    void RequestShutdown();
    void Uninit();

    ~RequestProcessor();
};

extern RequestProcessor requestProcessor;
