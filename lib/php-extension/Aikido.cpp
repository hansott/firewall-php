/* Aikido runtime extension for PHP */
#include "Includes.h"
#include "Utils.h"
#include "Handle.h"

ZEND_DECLARE_MODULE_GLOBALS(aikido)

void* aikido_agent_lib_handle = nullptr;
zval* server = nullptr;

#if PHP_VERSION_ID < 80000
	bool exit_current_request = false;

	static void (*original_zend_execute_ex)(zend_execute_data *execute_data) = NULL;

	void aikido_zend_execute_ex(zend_execute_data *execute_data) {
		if (exit_current_request) {
			AIKIDO_LOG_INFO("Current request is marked for exit. Bailing out...\n");
			zend_bailout();
		}
		original_zend_execute_ex(execute_data);
	}
#endif

PHP_MINIT_FUNCTION(aikido)
{
	aikido_log_init();

	std::string log_level = get_env_string("AIKIDO_LOG_LEVEL", "DEBUG");
	std::string token = get_env_string("AIKIDO_TOKEN", "");
	std::string endpoint = get_env_string("AIKIDO_ENDPOINT", "https://guard.aikido.dev/");
	std::string config_endpoint = get_env_string("AIKIDO_CONFIG_ENDPOINT", "https://runtime.aikido.dev/");
	bool blocking = get_env_bool("AIKIDO_BLOCKING", false);
	bool localhost_allowed_by_default = get_env_bool("AIKIDO_LOCALHOST_ALLOWED_BY_DEFAULT", true);

	AIKIDO_GLOBAL(log_level) = aikido_log_level_from_str(log_level);
	AIKIDO_GLOBAL(blocking) = blocking;
	AIKIDO_GLOBAL(socket_path) = aikido_generate_socket_path();

	AIKIDO_LOG_INFO("MINIT started!\n");

	for ( auto& it : HOOKED_FUNCTIONS ) {
		zend_function* function_data = (zend_function*)zend_hash_str_find_ptr(CG(function_table), it.first.c_str(), it.first.length());
		if (function_data == NULL) {
			AIKIDO_LOG_WARN("Function \"%s\" does not exist!\n", it.first.c_str());
			continue;
		}
		if (it.second.original_handler) {
			AIKIDO_LOG_WARN("Function \"%s\" already hooked (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
			continue;
		}

		it.second.original_handler = function_data->internal_function.handler;
		function_data->internal_function.handler = aikido_generic_handler;
		AIKIDO_LOG_INFO("Hooked function \"%s\" (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
	}

	for ( auto& it : HOOKED_METHODS ) {
		zend_class_entry *class_entry = (zend_class_entry *)zend_hash_str_find_ptr(CG(class_table), it.first.class_name.c_str(), it.first.class_name.length());
		if (class_entry == NULL) {
			AIKIDO_LOG_WARN("Class \"%s\" does not exist!\n", it.first.class_name.c_str());
			continue;
		}

		zend_function *method = (zend_function*)zend_hash_str_find_ptr(&class_entry->function_table, it.first.method_name.c_str(), it.first.method_name.length());
		if (method == NULL) {
			AIKIDO_LOG_WARN("Method \"%s->%s\" does not exist!\n", it.first.class_name.c_str(), it.first.method_name.c_str());
			continue;
		}

		if (it.second.original_handler) {
			AIKIDO_LOG_WARN("Method \"%s->%s\" already hooked (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
			continue;
		}

		it.second.original_handler = method->internal_function.handler;
		method->internal_function.handler = aikido_generic_handler;
		AIKIDO_LOG_INFO("Hooked method \"%s->%s\" (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
	}

	#if PHP_VERSION_ID < 80000
		original_zend_execute_ex = zend_execute_ex;
		zend_execute_ex = aikido_zend_execute_ex;
	#endif

	std::string sapi_name(sapi_module.name);
	AIKIDO_LOG_INFO("SAPI: %s\n", sapi_name.c_str());

	/* If SAPI name is "cli" run in "simple" mode */ 
	if (sapi_name == "cli") {
		AIKIDO_LOG_INFO("MINIT finished earlier because we run in CLI mode!\n");
		return SUCCESS;
	}

	json initData = {
		{"token", token},
		{"socket_path", AIKIDO_GLOBAL(socket_path)},
		{"platform_name", sapi_name},
		{"platform_version", PHP_VERSION},
		{"endpoint", endpoint},
		{"config_endpoint", config_endpoint},
		{"log_level", log_level},
		{"blocking", blocking},
		{"localhost_allowed_by_default", localhost_allowed_by_default}};

	std::string aikido_agent_lib_handle_path = "/opt/aikido-" + std::string(PHP_AIKIDO_VERSION) + "/aikido-agent.so";
	aikido_agent_lib_handle = dlopen(aikido_agent_lib_handle_path.c_str(), RTLD_LAZY);
    if (!aikido_agent_lib_handle) {
		AIKIDO_LOG_ERROR("Error loading the Aikido Agent library from %s: %s!\n", aikido_agent_lib_handle_path.c_str(), dlerror());
        return SUCCESS;
    }

    AgentInitFn agent_init_fn = (AgentInitFn)dlsym(aikido_agent_lib_handle, "AgentInit");
    if (!agent_init_fn) {
		AIKIDO_LOG_ERROR("Error loading symbol 'AgentInit' from the Aikido Agent library: %s!\n", dlerror());
        dlclose(aikido_agent_lib_handle);
        return SUCCESS;
    }

    AIKIDO_LOG_INFO("Initializing Aikido Agent...\n");

	std::string initDataString = initData.dump();

	int initOk = agent_init_fn(GoCreateString(initDataString));
	if (!initOk) {
		AIKIDO_LOG_INFO("Aikido Agent initialization failed!\n");
	}
	else {
		AIKIDO_LOG_INFO("Aikido Agent initialization succeded!\n");
	}

	AIKIDO_LOG_INFO("MINIT finished!\n");
	return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(aikido)
{
	AIKIDO_LOG_DEBUG("MSHUTDOWN started!\n");

	std::string sapi_name(sapi_module.name);
	AIKIDO_LOG_DEBUG("SAPI: %s\n", sapi_name.c_str());

	/* If SAPI name is "cli" run in "simple" mode */
	if (sapi_name == "cli") {
		AIKIDO_LOG_INFO("MSHUTDOWN finished earlier because we run in CLI mode!\n");
		aikido_log_uninit();
		return SUCCESS;
	}

	if (aikido_agent_lib_handle) {
		AgentUninitFn agent_uninit_fn = (AgentUninitFn)dlsym(aikido_agent_lib_handle, "AgentUninit");
		if (agent_uninit_fn) {
			AIKIDO_LOG_INFO("Uninitializing Aikido Agent library...\n");
			agent_uninit_fn();
			AIKIDO_LOG_INFO("Aikido Agent library uninitialized!\n");
		}
		else {
			AIKIDO_LOG_ERROR("Error loading symbol 'AgentUninit' from Aikido Agent library: %s!\n", dlerror());		
		}
		dlclose(aikido_agent_lib_handle);
		aikido_agent_lib_handle = nullptr;
	}

	AIKIDO_LOG_DEBUG("MSHUTDOWN finished!\n");

	aikido_log_uninit();

	return SUCCESS;
}

void* aikido_request_processor_lib_handle = nullptr;
bool request_processor_loading_failed = false;
RequestProcessorContextInitFn request_processor_context_init_fn = nullptr;
RequestProcessorOnEventFn request_processor_on_event_fn = nullptr;
RequestProcessorGetBlockingModeFn request_processor_get_blocking_mode_fn = nullptr;

PHP_RINIT_FUNCTION(aikido) {
	AIKIDO_LOG_DEBUG("RINIT started!\n");

	if (!aikido_request_processor_lib_handle && !request_processor_loading_failed) {
		std::string aikido_request_processor_lib_path = "/opt/aikido-" + std::string(PHP_AIKIDO_VERSION) + "/aikido-request-processor.so";
		aikido_request_processor_lib_handle = dlopen(aikido_request_processor_lib_path.c_str(), RTLD_LAZY);
		if (!aikido_request_processor_lib_handle) {
			AIKIDO_LOG_ERROR("Error loading the Aikido Request Processor library from %s: %s!\n", aikido_request_processor_lib_path.c_str(), dlerror());
			request_processor_loading_failed = true;
			return SUCCESS;
		}

		json initData = {
			{"log_level", aikido_log_level_str((AIKIDO_LOG_LEVEL)AIKIDO_GLOBAL(log_level))},
			{"socket_path", AIKIDO_GLOBAL(socket_path)},
			{"trust_proxy", get_env_bool("AIKIDO_TRUST_PROXY", true)},
			{"localhost_allowed_by_default", get_env_bool("AIKIDO_LOCALHOST_ALLOWED_BY_DEFAULT", true)},
			{"sapi", sapi_module.name}};

		std::string initDataString = initData.dump();

		AIKIDO_LOG_DEBUG("Initializing Aikido Request Processor...\n");

		RequestProcessorInitFn request_processor_init_fn = (RequestProcessorInitFn)dlsym(aikido_request_processor_lib_handle, "RequestProcessorInit");
		request_processor_context_init_fn = (RequestProcessorContextInitFn)dlsym(aikido_request_processor_lib_handle, "RequestProcessorContextInit");
		request_processor_on_event_fn = (RequestProcessorOnEventFn)dlsym(aikido_request_processor_lib_handle, "RequestProcessorOnEvent");
		request_processor_get_blocking_mode_fn = (RequestProcessorGetBlockingModeFn)dlsym(aikido_request_processor_lib_handle, "RequestProcessorGetBlockingMode");
		if (!request_processor_init_fn ||
			!request_processor_context_init_fn ||
			!request_processor_on_event_fn ||
			!request_processor_get_blocking_mode_fn ||
			!request_processor_init_fn(GoCreateString(initDataString)))
		{
			AIKIDO_LOG_ERROR("Failed to initialize Aikido Request Processor library: %s!\n", dlerror());
			dlclose(aikido_request_processor_lib_handle);
			aikido_request_processor_lib_handle = nullptr;
			request_processor_loading_failed = true;
			return SUCCESS;
		}

		AIKIDO_LOG_DEBUG("Aikido Request Processor initialized successfully!\n");
	}

	if (!request_processor_loading_failed) {
		zend_string *server_str = zend_string_init("_SERVER", sizeof("_SERVER") - 1, 0);
		if (server_str) {
			/* Guarantee that "_SERVER" global variable is initialized for the current request */
			zend_is_auto_global(server_str);
			zend_string_release(server_str);

			server = zend_hash_str_find(&EG(symbol_table), "_SERVER", sizeof("_SERVER") - 1);

			GoRequestProcessorContextInit();

			if (send_request_init_metadata_event() == EXIT) {
				#if PHP_VERSION_ID < 80000
					AIKIDO_LOG_INFO("Marking current request for exit!\n");
					exit_current_request = true;
				#endif
			}
		}
	}
	
	AIKIDO_LOG_DEBUG("RINIT finished!\n");
	return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(aikido) {
	AIKIDO_LOG_DEBUG("RSHUTDOWN started!\n");

	#if PHP_VERSION_ID < 80000
		exit_current_request = false;
	#endif

	/*
	if (aikido_request_processor_lib_handle)) {
		RequestProcessorUninitFn request_processor_uninit_fn = (RequestProcessorUninitFn)dlsym(aikido_globals->aikido_request_processor_lib_handle, "AgentUninit");
		if (!request_processor_init_fn || !request_processor_init_fn(GoCreateString(initDataString))) {
			AIKIDO_LOG_INFO("Failed to initialized request processor lib!\n");
			dlclose(aikido_request_processor_lib_handle);
			aikido_request_processor_lib_handle = nullptr;
			request_processor_failed = true;
			return SUCCESS;
		}
	}
	*/

	send_request_shutdown_metadata_event();

	AIKIDO_LOG_DEBUG("RSHUTDOWN finished!\n");
	return SUCCESS;
}


PHP_MINFO_FUNCTION(aikido)
{
	php_info_print_table_start();
	php_info_print_table_row(2, "aikido support", "enabled");
	php_info_print_table_end();
}

// Exports the "\aikido\set_user" function, to be called from PHP user code.
// Receives two parameters: id and name (both strings).
// Returns true if the setting of the user succeeded, false otherwise.
ZEND_FUNCTION(set_user) {
	char *id;
	size_t id_len;
	char *name;
	size_t name_len;

	// parse parameters
	ZEND_PARSE_PARAMETERS_START(2, 2)
		Z_PARAM_STRING(id, id_len)
		Z_PARAM_STRING(name, name_len)
	ZEND_PARSE_PARAMETERS_END();

	RETURN_BOOL(send_user_event(std::string(id, id_len), std::string(name, name_len)));	
}

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_aikido_set_user, 0, 2, _IS_BOOL, 0)
	ZEND_ARG_TYPE_INFO(0, id, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
ZEND_END_ARG_INFO()

static const zend_function_entry ext_functions[] = {
	ZEND_NS_FE("aikido", set_user, arginfo_aikido_set_user)
	ZEND_FE_END
};

zend_module_entry aikido_module_entry = {
	STANDARD_MODULE_HEADER,
	"aikido",					/* Extension name */
	ext_functions,				/* zend_function_entry */
	PHP_MINIT(aikido),			/* PHP_MINIT - Module initialization */
	PHP_MSHUTDOWN(aikido),		/* PHP_MSHUTDOWN - Module shutdown */
	PHP_RINIT(aikido),			/* PHP_RINIT - Request initialization */
	PHP_RSHUTDOWN(aikido),		/* PHP_RSHUTDOWN - Request shutdown */
	PHP_MINFO(aikido),			/* PHP_MINFO - Module info */
	PHP_AIKIDO_VERSION,			/* Version */
	PHP_MODULE_GLOBALS(aikido),	/* Module globals */
	NULL,						/* PHP_GINIT – Globals initialization */
	NULL,						/* PHP_GSHUTDOWN – Globals shutdown */
	NULL,
	STANDARD_MODULE_PROPERTIES_EX
};

#ifdef COMPILE_DL_AIKIDO
# ifdef ZTS
ZEND_TSRMLS_CACHE_DEFINE()
# endif
ZEND_GET_MODULE(aikido)
#endif
