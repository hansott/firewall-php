#include "Includes.h"

void helper_handle_pre_shell_execution(std::string cmd, EVENT_ID &eventId) {
    eventCache.cmd = cmd;
    eventId = EVENT_PRE_SHELL_EXECUTED;
}

std::string build_command_line_from_array(HashTable*) {
    return "";
}

AIKIDO_HANDLER_FUNCTION(handle_shell_execution) {
    zend_string *cmd = NULL;

    ZEND_PARSE_PARAMETERS_START(0, -1)
    Z_PARAM_OPTIONAL
    Z_PARAM_STR(cmd)
    ZEND_PARSE_PARAMETERS_END();

    if (!cmd) {
        return;
    }

    helper_handle_pre_shell_execution(ZSTR_VAL(cmd), eventId);
}


AIKIDO_HANDLER_FUNCTION(handle_shell_execution_with_array) {
    zval *cmd = nullptr;

    ZEND_PARSE_PARAMETERS_START(0, -1)
    Z_PARAM_OPTIONAL
    Z_PARAM_ZVAL(cmd)
    ZEND_PARSE_PARAMETERS_END();

    if (Z_TYPE_P(cmd) == IS_STRING) {
        zend_string* cmdStr = Z_STR_P(cmd);
        if (!cmdStr) {
            return;
        }
        helper_handle_pre_shell_execution(ZSTR_VAL(cmdStr), eventId);
    }
    else if (Z_TYPE_P(cmd) == IS_ARRAY) {
        HashTable* cmdArr = Z_ARRVAL_P(cmd);
        if (!cmdArr) {
            return;
        }
        helper_handle_pre_shell_execution(build_command_line_from_array(cmdArr), eventId);    
    }
}
