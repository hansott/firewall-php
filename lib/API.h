#include <stdlib.h>

enum EVENT_ID {
    NO_EVENT_ID,

    EVENT_PRE_REQUEST,
    EVENT_POST_REQUEST,
    EVENT_PRE_USER,
    EVENT_PRE_OUTGOING_REQUEST,
    EVENT_POST_OUTGOING_REQUEST,
    EVENT_PRE_SHELL_EXECUTED,
    EVENT_PRE_PATH_ACCESSED,
    EVENT_PRE_SQL_QUERY_EXECUTED,

    MAX_EVENT_ID
};

enum CALLBACK_ID
{
    NO_CALLBACK_ID,

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
    CONTEXT_BODY,

    CONTEXT_USER_ID,
    CONTEXT_USER_NAME,

    FUNCTION_NAME,

    OUTGOING_REQUEST_URL,
    OUTGOING_REQUEST_EFFECTIVE_URL,
    OUTGOING_REQUEST_PORT,
    OUTGOING_REQUEST_RESOLVED_IP,

    CMD,

    FILENAME,
    FILENAME2,

    SQL_QUERY,
    SQL_DIALECT,

    MODULE,

    MAX_CALLBACK_ID
};

typedef char *(*ContextCallback)(int);

static char *call(ContextCallback callback, int callback_id) { return callback(callback_id); }
