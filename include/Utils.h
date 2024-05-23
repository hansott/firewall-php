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
#define AIKIDO_REGISTER_METHOD_HANDLER(class_name, method_name) { std::string(#class_name "_" #method_name), { handle_ ## class_name ## _ ## method_name, nullptr } }


#define AIKIDO_GET_FUNCTION_NAME() (ZSTR_VAL(execute_data->func->common.function_name))

#define AIKIDO_HANDLER_START() php_printf("[AIKIDO-C++] Handler called for \"%s\"!\n", AIKIDO_GET_FUNCTION_NAME());

#define AIKIDO_HANDLER_END() HOOKED_FUNCTIONS[AIKIDO_GET_FUNCTION_NAME()].original_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU);
