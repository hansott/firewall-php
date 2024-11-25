#include "Includes.h"

void helper_handle_pre_shell_execution(char *filename, EVENT_ID &eventId) {
    eventCache.cmd = ZSTR_VAL(cmd);
    eventId = EVENT_PRE_SHELL_EXECUTED;
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
    zend_string *cmd = nullptr;
	HashTable *cmdHashTable = nullptr;

    ZEND_PARSE_PARAMETERS_START(0, -1)
    Z_PARAM_OPTIONAL
    Z_PARAM_ARRAY_HT_OR_STR(command_ht, command_str)
    ZEND_PARSE_PARAMETERS_END();

    if (cmd) {
        return;
    }

    helper_handle_pre_shell_execution(ZSTR_VAL(cmd), eventId);
}
