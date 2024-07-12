/* Aikido runtime extension for PHP */
#include "Includes.h"
#include "Utils.h"
#include "Handle.h"

ZEND_DECLARE_MODULE_GLOBALS(aikido)

void* aikido_agent_lib_handle = nullptr;

static PHP_GINIT_FUNCTION(aikido)
{
#if defined(COMPILE_DL_BCMATH) && defined(ZTS)
	ZEND_TSRMLS_CACHE_UPDATE();
#endif

	aikido_log_init();

	aikido_globals->log_level = 3;

	AIKIDO_LOG_DEBUG("GInit started (PID = %d)!\n", getpid());

	aikido_globals->blocking = config_override_with_env_bool(false, "AIKIDO_BLOCKING");

	json initData = {
		{ "token", config_override_with_env("", "AIKIDO_TOKEN") },
		{ "endpoint", config_override_with_env("https://guard.aikido.dev/", "AIKIDO_ENDPOINT") },
		{ "log_level", config_override_with_env("DEBUG", "AIKIDO_LOG_LEVEL") },
		{ "blocking", aikido_globals->blocking }
	};

	aikido_agent_lib_handle = dlopen("/opt/aikido/aikido_agent.so", RTLD_LAZY);
    if (!aikido_agent_lib_handle) {
		AIKIDO_LOG_ERROR("Error loading the Aikido Agent library: %s!\n", dlerror());
        return;
    }

    AgentInitFn agent_init_fn = (AgentInitFn)dlsym(aikido_agent_lib_handle, "AgentInit");
    if (!agent_init_fn) {
		AIKIDO_LOG_ERROR("Error loading symbol 'AgentInit' from the Aikido Agent library: %s!\n", dlerror());
        dlclose(aikido_agent_lib_handle);
        return;
    }

    AIKIDO_LOG_DEBUG("Initializing Aikido Agent...\n");

	std::string initDataString = initData.dump();

	int initOk = agent_init_fn(GoCreateString(initDataString));

	AIKIDO_LOG_DEBUG("Aikido Agent initialized with status: %d!\n", initOk);

	AIKIDO_LOG_DEBUG("GInit finished!\n");
}

static PHP_GSHUTDOWN_FUNCTION(aikido)
{
#if defined(COMPILE_DL_BCMATH) && defined(ZTS)
	ZEND_TSRMLS_CACHE_UPDATE();
#endif
	AIKIDO_LOG_DEBUG("GShutdown started (PID = %d)!\n", getpid());

	if (aikido_agent_lib_handle) {
		AgentUninitFn agent_uninit_fn = (AgentUninitFn)dlsym(aikido_agent_lib_handle, "AgentUninit");
		if (agent_uninit_fn) {
			AIKIDO_LOG_DEBUG("Uninitializing Aikido Agent...\n");
			agent_uninit_fn();
			AIKIDO_LOG_DEBUG("Aikido Agent uninitialized!\n");
		}
		else {
			AIKIDO_LOG_ERROR("Error loading symbol 'AgentUninit' from Aikido Agent library: %s!\n", dlerror());		
		}
		dlclose(aikido_agent_lib_handle);
		aikido_agent_lib_handle = nullptr;
	}
	
	AIKIDO_LOG_INFO("GShutdown finished!\n");
	aikido_log_uninit();
}

PHP_MINIT_FUNCTION(aikido)
{
	AIKIDO_LOG_INFO("MInit started (PID = %d)!\n", getpid());
	for ( auto& it : HOOKED_FUNCTIONS ) {
		zend_function* function_data = (zend_function*)zend_hash_str_find_ptr(CG(function_table), it.first.c_str(), it.first.length());
		if (function_data == NULL) {
			AIKIDO_LOG_DEBUG("Function \"%s\" does not exist!\n", it.first.c_str());
			continue;
		}
		if (it.second.original_handler) {
			AIKIDO_LOG_DEBUG("Function \"%s\" already hooked (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
			continue;
		}

		it.second.original_handler = function_data->internal_function.handler;
		function_data->internal_function.handler = aikido_generic_handler;
		AIKIDO_LOG_DEBUG("Hooked function \"%s\" (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
	}

	for ( auto& it : HOOKED_METHODS ) {
		zend_class_entry *class_entry = (zend_class_entry *)zend_hash_str_find_ptr(CG(class_table), it.first.class_name.c_str(), it.first.class_name.length());
		if (class_entry == NULL) {
			AIKIDO_LOG_DEBUG("Class \"%s\" does not exist!\n", it.first.class_name.c_str());
			continue;
		}

		zend_function *method = (zend_function*)zend_hash_str_find_ptr(&class_entry->function_table, it.first.method_name.c_str(), it.first.method_name.length());
		if (method == NULL) {
			AIKIDO_LOG_DEBUG("Method \"%s->%s\" does not exist!\n", it.first.class_name.c_str(), it.first.method_name.c_str());
			continue;
		}

		if (it.second.original_handler) {
			AIKIDO_LOG_DEBUG("Method \"%s->%s\" already hooked (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
			continue;
		}

		it.second.original_handler = method->internal_function.handler;
		method->internal_function.handler = aikido_generic_handler;
		AIKIDO_LOG_DEBUG("Hooked method \"%s->%s\" (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
	}

	AIKIDO_LOG_INFO("MInit finished!\n");
	return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(aikido)
{
	/* Unregister Aikido-specific (log level, blocking, token, ...) entries in php.ini */
	AIKIDO_LOG_INFO("MShutdown started (PID = %d)!\n", getpid());
	AIKIDO_LOG_INFO("MShutdown finished!\n");
	return SUCCESS;
}

void* aikido_request_processor_lib_handle = nullptr;
bool request_processor_loading_failed = false;
RequestProcessorOnEventFn request_processor_on_event_fn = nullptr;

PHP_RINIT_FUNCTION(aikido) {
	AIKIDO_LOG_DEBUG("RInit started (PID = %d)!\n", getpid());

	if (!aikido_request_processor_lib_handle && !request_processor_loading_failed) {
		aikido_request_processor_lib_handle = dlopen("/opt/aikido/aikido_request_processor.so", RTLD_LAZY);
		if (!aikido_request_processor_lib_handle) {
			AIKIDO_LOG_ERROR("Error loading the Aikido Request Processor library: %s!\n", dlerror());
			request_processor_loading_failed = true;
			return SUCCESS;
		}
		
		json initData = {
			{ "log_level", "DEBUG" }
		};

		std::string initDataString = initData.dump();

		AIKIDO_LOG_DEBUG("Initializing Aikido Request Processor for PID %d...\n", getpid());

		RequestProcessorInitFn request_processor_init_fn = (RequestProcessorInitFn)dlsym(aikido_request_processor_lib_handle, "RequestProcessorInit");
		request_processor_on_event_fn = (RequestProcessorOnEventFn)dlsym(aikido_request_processor_lib_handle, "RequestProcessorOnEvent");
		if (!request_processor_init_fn || 
			!request_processor_on_event_fn ||
			!request_processor_init_fn(GoCreateString(initDataString))) {
			AIKIDO_LOG_ERROR("Failed to initialize Aikido Request Processor library: %s!\n", dlerror());
			dlclose(aikido_request_processor_lib_handle);
			aikido_request_processor_lib_handle = nullptr;
			request_processor_loading_failed = true;
			return SUCCESS;
		}
	}

	AIKIDO_LOG_DEBUG("RInit finished!\n");
	return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(aikido) {
	AIKIDO_LOG_DEBUG("RShutdown started (PID = %d)!\n", getpid());

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

	AIKIDO_LOG_DEBUG("RShutdown finished!\n");
	return SUCCESS;
}


PHP_MINFO_FUNCTION(aikido)
{
	php_info_print_table_start();
	php_info_print_table_row(2, "aikido support", "enabled");
	php_info_print_table_end();
}

static const zend_function_entry ext_functions[] = {
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
	PHP_GINIT(aikido),			/* PHP_GINIT – Globals initialization */
	PHP_GSHUTDOWN(aikido),		/* PHP_GSHUTDOWN – Globals shutdown */
	NULL,
	STANDARD_MODULE_PROPERTIES_EX
};

#ifdef COMPILE_DL_AIKIDO
# ifdef ZTS
ZEND_TSRMLS_CACHE_DEFINE()
# endif
ZEND_GET_MODULE(aikido)
#endif
