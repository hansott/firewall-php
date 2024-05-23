#include "HandleUrls.h"
#include "HandleShellExecution.h"
#include "HandlePDO.h"

#include "Utils.h"

unordered_map<std::string, PHP_HANDLERS> HOOKED_FUNCTIONS = {
	AIKIDO_REGISTER_FUNCTION_HANDLER(curl_init),
	AIKIDO_REGISTER_FUNCTION_HANDLER(curl_setopt),

	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(exec,       handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(shell_exec, handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(system,     handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(passthru,   handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(popen,      handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(proc_open,  handle_shell_execution)
};

unordered_map<std::string, PHP_HANDLERS> HOOKED_METHODS = {
	AIKIDO_REGISTER_METHOD_HANDLER(PDO, __construct),
	AIKIDO_REGISTER_METHOD_HANDLER(PDO, query)
};