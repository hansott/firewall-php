#include "HandleUrls.h"
#include "HandleShellExecution.h"
#include "HandlePathAccess.h"
#include "HandlePDO.h"
#include "Cache.h"
#include "Utils.h"
#include "Actions.h"

unordered_map<std::string, PHP_HANDLERS> HOOKED_FUNCTIONS = {
	/* Outgoing request */
	AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST(curl_exec),

	/* Shell execution */
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(exec, handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(shell_exec, handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(system, handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(passthru, handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(popen, handle_shell_execution),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(proc_open, handle_shell_execution),

	/* Path access */
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chdir, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chgrp, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chmod, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chown, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(copy, handle_pre_file_path_access_2),
	AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST_EX(file, handle_pre_file_path_access, handle_post_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST_EX(file_get_contents, handle_pre_file_path_access, handle_post_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(file_put_contents, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST_EX(fopen, handle_pre_file_path_access, handle_post_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(lchgrp, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(lchown, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(link, handle_pre_file_path_access_2),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(mkdir, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(move_uploaded_file, handle_pre_file_path_access_2),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(opendir, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(parse_ini_file, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(readfile, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(readlink, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(rename, handle_pre_file_path_access_2),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(rmdir, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(scandir, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(symlink, handle_pre_file_path_access_2),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(touch, handle_pre_file_path_access),
	AIKIDO_REGISTER_FUNCTION_HANDLER_EX(unlink, handle_pre_file_path_access),
};

unordered_map<AIKIDO_METHOD_KEY, PHP_HANDLERS, AIKIDO_METHOD_KEY_HASH> HOOKED_METHODS = {
	AIKIDO_REGISTER_METHOD_HANDLER(pdo, query)
};

ZEND_NAMED_FUNCTION(aikido_generic_handler) {
	AIKIDO_LOG_DEBUG("Aikido generic handler started!\n");

	zif_handler original_handler = nullptr;
	aikido_handler post_handler = nullptr;

	std::string outputEvent;
	bool caughtException = false;

	eventCache.Reset();
	eventCache.functionName = AIKIDO_GET_FUNCTION_NAME();

	if (request_processor_on_event_fn) {
		try {
			zend_execute_data *exec_data = EG(current_execute_data);
			zend_function *func = exec_data->func;
			zend_class_entry* executed_scope = zend_get_executed_scope();
			
			std::string function_name(ZSTR_VAL(func->common.function_name));
			function_name = to_lowercase(function_name);

			aikido_handler handler = nullptr;

			std::string scope_name = function_name;
			AIKIDO_LOG_DEBUG("Function name: %s\n", scope_name.c_str());
			if (HOOKED_FUNCTIONS.find(function_name) != HOOKED_FUNCTIONS.end()) {
				handler = HOOKED_FUNCTIONS[function_name].handler;
				post_handler = HOOKED_FUNCTIONS[function_name].post_handler;
				original_handler = HOOKED_FUNCTIONS[function_name].original_handler;
			}
			else if (executed_scope) {
				/* A method was executed (executed_scope stores the name of the current class) */

				std::string class_name(ZSTR_VAL(executed_scope->name));
				class_name = to_lowercase(class_name);

				scope_name = class_name + "->" + function_name;

				AIKIDO_METHOD_KEY method_key(class_name, function_name);

				AIKIDO_LOG_DEBUG("Method name: %s\n", scope_name.c_str());

				if (HOOKED_METHODS.find(method_key) == HOOKED_METHODS.end()) {
					AIKIDO_LOG_DEBUG("Method not found! Returning!\n");
					return;
				}

				handler = HOOKED_METHODS[method_key].handler;
				post_handler = HOOKED_METHODS[method_key].post_handler;
				original_handler = HOOKED_METHODS[method_key].original_handler;
			}
			else {
				AIKIDO_LOG_DEBUG("Nothing matches the current handler! Returning!\n");
				return;
			}

			AIKIDO_LOG_DEBUG("Calling handler for \"%s\"!\n", scope_name.c_str());

			EVENT_ID eventId = NO_EVENT_ID;
			/*
				The handler for a specific PHP function that we hook can set an event ID
				to be sent to the Go libary (request processor).
				This will notify the Go library that an event has happend in the PHP layer.
				The event ID is initialy empty and it's only sent to Go only if the C++ handler
				for the currently hooked function sets it.
			*/
			handler(INTERNAL_FUNCTION_PARAM_PASSTHRU, eventId);

			if (eventId != NO_EVENT_ID) {
				std::string outputEvent;
				GoRequestProcessorOnEvent(eventId, outputEvent);
				if (IsBlockingEnabled() && action.Execute(outputEvent) == BLOCK) {
					// exit generic handler and do not call the original handler
					// thus blocking the execution 
					return;
				}
			}
		}
		catch (const std::exception& e) {
			caughtException = true;
			AIKIDO_LOG_ERROR("Exception encountered in generic handler: %s\n", e.what());
		}
	}

	if (original_handler) {
		original_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU);

		if (!caughtException && post_handler) {
			EVENT_ID eventId = NO_EVENT_ID;
			/*
				The handler for a specific PHP function that we hook can set an event ID
				to be sent to the Go libary (request processor).
				This will notify the Go library that an event has happend in the PHP layer.
				The event ID is initialy empty and it's only sent to Go only if the C++ handler
				for the currently hooked function sets it.
			*/
			post_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU, eventId);
			if (eventId != NO_EVENT_ID) {
				std::string outputEvent;
				GoRequestProcessorOnEvent(eventId, outputEvent);
				if (IsBlockingEnabled()) {
					action.Execute(outputEvent);
				}
			}
		}
	}

	AIKIDO_LOG_DEBUG("Aikido generic handler ended!\n");
}