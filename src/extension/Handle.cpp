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

unordered_map<AIKIDO_METHOD_KEY, PHP_HANDLERS, AIKIDO_METHOD_KEY_HASH> HOOKED_METHODS = {
	AIKIDO_REGISTER_METHOD_HANDLER(pdo, __construct),
	AIKIDO_REGISTER_METHOD_HANDLER(pdo, query)
};

enum ACTION {
	CONTINUE,
	BLOCK
};

ACTION aikido_execute_output(json event) {
	if ( event["action"] == "throw" ) {
		std::string message = event["message"].get<std::string>();
		int code = event["code"].get<int>();
		zend_throw_exception(zend_exception_get_default(), message.c_str(), code);
		return BLOCK;
	}
	return CONTINUE;
}

ZEND_NAMED_FUNCTION(aikido_generic_handler) {
	AIKIDO_LOG_DEBUG("Aikido generic handler started!\n");

	zif_handler original_handler = nullptr;
	
	try {
		zend_execute_data *exec_data = EG(current_execute_data);
		zend_function *func = exec_data->func;
		zend_class_entry* executed_scope = zend_get_executed_scope();
		
		std::string function_name(ZSTR_VAL(func->common.function_name));
		function_name = to_lowercase(function_name);

		aikido_handler handler = nullptr;

		std::string scope_name;

		if (executed_scope) {
			/* A method was executed (executed_scope stores the name of the current class) */

			std::string class_name(ZSTR_VAL(executed_scope->name));
			class_name = to_lowercase(class_name);

			scope_name = class_name + "->" + function_name;

			AIKIDO_METHOD_KEY method_key(class_name, function_name);

			AIKIDO_LOG_DEBUG("Method name: %s\n", scope_name.c_str());

			if (HOOKED_METHODS.find(method_key) == HOOKED_METHODS.end()) {
				return;
			}

			handler = HOOKED_METHODS[method_key].handler;
			original_handler = HOOKED_METHODS[method_key].original_handler;
		}
		else {
			/* A function was executed (executed_scope is null) */
			scope_name = function_name;
			AIKIDO_LOG_DEBUG("Function name: %s\n", scope_name.c_str());
			if (HOOKED_FUNCTIONS.find(function_name) == HOOKED_FUNCTIONS.end()) {
				return;
			}
			handler = HOOKED_FUNCTIONS[function_name].handler;
			original_handler = HOOKED_FUNCTIONS[function_name].original_handler;
		}

		AIKIDO_LOG_DEBUG("Handler called for \"%s\"!\n", scope_name.c_str());

		json inputEvent;
		handler(INTERNAL_FUNCTION_PARAM_PASSTHRU, inputEvent);

		if (!inputEvent.empty()) {
			json outputEvent = GoOnEvent(inputEvent);
			if (AIKIDO_GLOBAL(blocking) == true && aikido_execute_output(outputEvent) == BLOCK) {
				// exit generic handler and do not call the original handler
				// thus blocking the execution 
				return;
			}
		}
	}
	catch (const std::exception& e) {
		AIKIDO_LOG_WARN("Exception encountered in generic handler: %s\n", e.what());
	}
	
	if (original_handler) {
		original_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU);
	}

	AIKIDO_LOG_DEBUG("Aikido generic handler ended!\n");
}