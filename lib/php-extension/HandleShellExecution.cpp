#include "Includes.h"

AIKIDO_HANDLER_FUNCTION(handle_shell_execution) {
    zend_string *cmd = NULL;

    ZEND_PARSE_PARAMETERS_START(0, -1)
    Z_PARAM_OPTIONAL
    Z_PARAM_STR(cmd)
    ZEND_PARSE_PARAMETERS_END();

    if (!cmd) {
        return;
    }

    eventCache.cmd = ZSTR_VAL(cmd);

    eventId = EVENT_PRE_SHELL_EXECUTED;
}