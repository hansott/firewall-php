#include "GoWrappers.h"

GoString GoCreateString(std::string& s) {
    return GoString { s.c_str(), s.length() };
}

json GoRequestProcessorOnEvent(json& event) {
    std::string eventString = event.dump();

    if (!request_processor_on_event_fn) {
        return json::object();
    }
    
    AIKIDO_LOG_DEBUG("Sending event to GO\n");
    
    char* charPtr = request_processor_on_event_fn(GoCreateString(eventString));
    if (!charPtr) {
        return json::object();
    }
    
    std::string outputString(charPtr);
    free(charPtr);
    
    AIKIDO_LOG_DEBUG("Got event reply: %s\n", outputString.c_str());
    
    json output = json::parse(outputString);
    return output;
}

bool IsBlockingEnabled() {
    if (!request_processor_get_blocking_mode_fn) {
        return false;
    }
    int ret = request_processor_get_blocking_mode_fn();
    if (ret == -1) {
        return AIKIDO_GLOBAL(blocking);
    }
    return ret;
}