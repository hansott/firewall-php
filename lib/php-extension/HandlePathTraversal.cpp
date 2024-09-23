#include "HandlePathTraversal.h"
#include "Utils.h"
#include "Cache.h"

/* Handles PHP functions that have a file path as first parameter */
AIKIDO_HANDLER_FUNCTION(handle_file_path_access) {
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

    eventCache.filename = ZSTR_VAL(filename);
    if (filename2) {
        eventCache.filename2 = ZSTR_VAL(filename2);
    }
    eventId = EVENT_PRE_PATH_ACCESSED;
}
