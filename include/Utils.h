#pragma once

#include "Includes.h"

/* 
	Macro for registering an Aikido handler in the HOOKED_FUNCTIONS map.
	It takes as parameters the PHP function name to be hooked and C++ function 
		that should be called when that PHP function is executed.
	The nullptr part is a placeholder where the original function handler from
		the Zend framework will be stored at initialization when we run the hooking.
*/
#define AIKIDO_REGISTER_FUNCTION_HANDLER_EX(function_name, function_pointer) { std::string(#function_name), { function_pointer, nullptr } }

/*
	Shorthand version of AIKIDO_REGISTER_FUNCTION_HANDLER_EX that constructs automatically the C++ function to be called.
	For example, if function name is curl_init this macro will store { "curl_init", { handle_curl_init, nullptr } }.
*/
#define AIKIDO_REGISTER_FUNCTION_HANDLER(function_name) { std::string(#function_name), { handle_##function_name, nullptr } }

/*
	Similar to AIKIDO_REGISTER_FUNCTION_HANDLER, but for methods.
*/
#define AIKIDO_REGISTER_METHOD_HANDLER(class_name, method_name) { AIKIDO_METHOD_KEY(std::string(#class_name), std::string(#method_name)), { handle_ ## class_name ## _ ## method_name, nullptr } }


#define AIKIDO_GET_FUNCTION_NAME() (ZSTR_VAL(execute_data->func->common.function_name))

#define AIKIDO_FUNCTION_HANDLER_START() php_printf("[AIKIDO-C++] Handler called for \"%s\"!\n", AIKIDO_GET_FUNCTION_NAME());

#define AIKIDO_FUNCTION_HANDLER_END() HOOKED_FUNCTIONS[AIKIDO_GET_FUNCTION_NAME()].original_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU);

#define AIKIDO_METHOD_HANDLER_START() \
    zend_class_entry *executed_ce = zend_get_executed_scope(); \
    const char *class_name = executed_ce ? ZSTR_VAL(executed_ce->name) : "None"; \
    zend_function *executed_method = EG(current_execute_data)->func; \
    const char *method_name = executed_method->common.function_name ? ZSTR_VAL(executed_method->common.function_name) : "None"; \
	php_printf("[AIKIDO-C++] Handler called for \"%s->%s\"!\n", class_name, method_name);

inline std::string to_lowercase(const std::string& str) {
    std::string result = str;
    std::transform(result.begin(), result.end(), result.begin(), [](unsigned char c){ return std::tolower(c); });
    return result;
}

#define AIKIDO_METHOD_HANDLER_END() HOOKED_METHODS[AIKIDO_METHOD_KEY(to_lowercase(std::string(class_name)), to_lowercase(std::string(method_name)))].original_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU);

