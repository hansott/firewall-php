#include <stdlib.h>

typedef char* (*ContextCallback)(int);

static char* call(ContextCallback callback, int context_id) { return callback(context_id); }

enum {
    CONTEXT_REMOTE_ADDRESS,
    CONTEXT_ROUTE,
    CONTEXT_METHOD,
    CONTEXT_STATUS_CODE,
    CONTEXT_BODY,
    CONTEXT_HEADERS,
    CONTEXT_HEADER_X_FORWARDED_FOR,
    CONTEXT_COOKIES,
    CONTEXT_QUERY,
    CONTEXT_HTTPS,
    CONTEXT_URL
};