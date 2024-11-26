#include "Includes.h"

void helper_handle_pre_shell_execution(std::string cmd, EVENT_ID &eventId) {
    eventCache.cmd = cmd;
    eventId = EVENT_PRE_SHELL_EXECUTED;
}

std::string build_command_line_from_array(HashTable* cmdTokens) {
    std::string cmd;
    zval *cmdTokenVal;
    ZEND_HASH_FOREACH_VAL(cmdTokens, cmdTokenVal) {
        zend_string *cmdToken = zval_get_string(cmdTokenVal);
        if (!cmdToken) return "";
        cmd += ZSTR_VAL(cmdToken);
        cmd += " ";
	} ZEND_HASH_FOREACH_END();
    return cmd;
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
    zval *cmdVal = nullptr;

    ZEND_PARSE_PARAMETERS_START(0, -1)
    Z_PARAM_OPTIONAL
    Z_PARAM_ZVAL(cmdVal)
    ZEND_PARSE_PARAMETERS_END();

    if (Z_TYPE_P(cmdVal) == IS_STRING) {
        zend_string* cmdStr = Z_STR_P(cmdVal);
        if (!cmdStr) {
            return;
        }
        helper_handle_pre_shell_execution(ZSTR_VAL(cmdStr), eventId);
    }
#if PHP_VERSION_ID >= 70400
    else if (Z_TYPE_P(cmdVal) == IS_ARRAY) {
        HashTable* cmdTokens = Z_ARRVAL_P(cmdVal);
        if (!cmdTokens) {
            return;
        }
        std::string cmdString = build_command_line_from_array(cmdTokens);
        if (cmdString.empty()) {
            return;
        }
        helper_handle_pre_shell_execution(cmdString, eventId);    
    }
#endif
}
