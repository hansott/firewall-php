#include "HandleShellExecution.h"
#include "Utils.h"

ZEND_NAMED_FUNCTION(handle_shell_execution) {
	AIKIDO_FUNCTION_HANDLER_START();

	zend_string *cmd = NULL;

	ZEND_PARSE_PARAMETERS_START(1,-1)
		Z_PARAM_OPTIONAL
		Z_PARAM_STR(cmd)
	ZEND_PARSE_PARAMETERS_END();

	std::string cmdString(ZSTR_VAL(cmd));

	std::string functionNameString(AIKIDO_GET_FUNCTION_NAME());
	
	json shell_execution_event = {
		{ "event", "function_executed" },
		{ "data", {
			{ "function_name", functionNameString },
			{ "parameters", {
				{ "cmd", cmdString }
			} }
		} }
	};

	GoOnEvent(shell_execution_event);

	AIKIDO_FUNCTION_HANDLER_END();
}