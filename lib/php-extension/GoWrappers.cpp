#include "Includes.h"

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

/*
    Callback wrapper called by the RequestProcessor (GO) whenever it needs data from PHP (C++ extension).
*/
char* GoContextCallback(int context_id) {
    if (!server) {
        AIKIDO_LOG_WARN("_SERVER variables is not initialized!\n");
        return nullptr;
    }
    std::string ctx;
    std::string ret;
    switch (context_id) {
        case CONTEXT_REMOTE_ADDRESS:
            ctx = "REMOTE_ADDRESS";
            ret = extract_server_var("REMOTE_ADDR");
            break;
        case CONTEXT_METHOD:
            ctx = "METHOD";
            ret = extract_server_var("REQUEST_METHOD");
            break;
        case CONTEXT_ROUTE:
            ctx = "ROUTE";
            ret = extract_route();
            break;
        case CONTEXT_STATUS_CODE:
            ctx = "STATUS_CODE";
            ret = extract_status_code();
            break;
        case CONTEXT_BODY:
            ctx = "BODY";
            ret = extract_body();
            break;
        case CONTEXT_HEADER_X_FORWARDED_FOR:
            ctx = "HEADER_X_FORWARDED_FOR";
            ret = extract_server_var("HTTP_X_FORWARDED_FOR");
            break;
        case CONTEXT_COOKIES:
            ctx = "COOKIES";
            ret = extract_server_var("HTTP_COOKIE");
            break;
        case CONTEXT_QUERY:
            ctx = "QUERY";
            ret = extract_server_var("QUERY_STRING");
            break;
        case CONTEXT_HTTPS:
            ctx = "HTTPS";
            ret = extract_server_var("HTTPS");
            break;
        case CONTEXT_URL:
            ctx = "URL";
            ret = extract_url();
            break;
        case CONTEXT_HEADERS:
            ctx = "HEADERS";
            ret = extract_headers();
            break;
        case CONTEXT_HEADER_USER_AGENT:
            ctx = "USER_AGENT";
            ret = extract_server_var("HTTP_USER_AGENT");
            break;
    }

    if (!ret.length()) {
        AIKIDO_LOG_DEBUG("Context callback %s -> NULL\n", ctx.c_str());
        return nullptr;
    }

    AIKIDO_LOG_DEBUG("Context callback %s -> %s\n", ctx.c_str(), ret.c_str());
    return strdup(ret.c_str());
}

bool GoRequestProcessorContextInit() {
    if (!request_processor_context_init_fn) {
        return false;
    }

    return request_processor_context_init_fn(GoContextCallback);
}

/*
    If the blocking mode is set from agent (different than -1), return that.
	Otherwise, return the env variable AIKIDO_BLOCKING.
*/
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