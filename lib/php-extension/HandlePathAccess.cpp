#include "HandlePathAccess.h"
#include "Utils.h"
#include "Cache.h"

/* Handles PHP functions that have a file path as first parameter */
AIKIDO_HANDLER_FUNCTION(handle_file_path_access) {
    zend_string *filename = NULL;

    ZEND_PARSE_PARAMETERS_START(0, -1)
        Z_PARAM_OPTIONAL
        Z_PARAM_STR(filename)
    ZEND_PARSE_PARAMETERS_END();

    if (!filename) {
        return;
    }

    // if filename starts with http:// or https://, it's a URL so we treat it as an outgoing request
     if (strncmp(ZSTR_VAL(filename), "http://", 7) == 0 || 
         strncmp(ZSTR_VAL(filename), "https://", 8) == 0) {
        eventId = EVENT_PRE_OUTGOING_REQUEST;
        eventCache.outgoingRequestUrl = ZSTR_VAL(filename);
    }
    else {
        eventId = EVENT_PRE_PATH_ACCESSED;
        eventCache.filename = ZSTR_VAL(filename);
    }
}

/* Handles PHP functions that have a file path as both first and second parameter */
AIKIDO_HANDLER_FUNCTION(handle_file_path_access_2) {
    zend_string *filename = NULL;
    zend_string *filename2 = NULL;

    ZEND_PARSE_PARAMETERS_START(0, -1)
        Z_PARAM_OPTIONAL
        Z_PARAM_STR(filename)
        Z_PARAM_STR(filename2)
    ZEND_PARSE_PARAMETERS_END();

    if (!filename) {
        return;
    }

    eventId = EVENT_PRE_PATH_ACCESSED;
    eventCache.filename = ZSTR_VAL(filename);
    if (filename2)
    {
        eventCache.filename2 = ZSTR_VAL(filename2);
    }
}
