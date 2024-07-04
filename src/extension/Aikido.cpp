/* Aikido runtime extension for PHP */
#include "Includes.h"
#include "Utils.h"
#include "Handle.h"

ZEND_DECLARE_MODULE_GLOBALS(aikido)

static PHP_GINIT_FUNCTION(aikido)
{
#if defined(COMPILE_DL_BCMATH) && defined(ZTS)
	ZEND_TSRMLS_CACHE_UPDATE();
#endif

	aikido_log_init();

	aikido_globals->log_level = 3;

	AIKIDO_LOG_INFO("GInit started!\n");

	aikido_globals->blocking = false;

	aikido_globals->blocking = config_override_with_env_bool(aikido_globals->blocking, "AIKIDO_BLOCKING");

	json initData = {
		{ "version", PHP_AIKIDO_VERSION },
		{ "log_level", "DEBUG" },
	};

	GoInit(initData);

	AIKIDO_LOG_INFO("GInit finished!\n");
}

static PHP_GSHUTDOWN_FUNCTION(aikido)
{
#if defined(COMPILE_DL_BCMATH) && defined(ZTS)
	ZEND_TSRMLS_CACHE_UPDATE();
#endif
	AIKIDO_LOG_INFO("GShutdown started!\n");
	
	GoUninit();

	AIKIDO_LOG_INFO("GShutdown finished!\n");
	aikido_log_uninit();
}

PHP_MINIT_FUNCTION(aikido)
{
	AIKIDO_LOG_INFO("MInit started!\n");
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
	AIKIDO_LOG_INFO("MShutdown started!\n");
	AIKIDO_LOG_INFO("MShutdown finished!\n");
	return SUCCESS;
}

PHP_RINIT_FUNCTION(aikido) {
	AIKIDO_LOG_INFO("RInit started!\n");
	AIKIDO_LOG_INFO("RInit finished!\n");

	return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(aikido) {
	AIKIDO_LOG_INFO("RShutdown started!\n");
	AIKIDO_LOG_INFO("RShutdown finished!\n");
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
