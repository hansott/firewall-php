#include <stdlib.h>

typedef char* (*ContextCallback)(int);

static char* call(ContextCallback callback, int context_id) { return callback(context_id); }

enum
{
    CONTEXT_REMOTE_ADDRESS,
    CONTEXT_HTTPS,
    CONTEXT_METHOD,
    CONTEXT_ROUTE,
    CONTEXT_URL,
    CONTEXT_QUERY,
    CONTEXT_STATUS_CODE,
    CONTEXT_HEADERS,
    CONTEXT_HEADER_X_FORWARDED_FOR,
    CONTEXT_HEADER_USER_AGENT,
    CONTEXT_COOKIES,
    CONTEXT_BODY
};