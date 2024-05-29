/* Aikido runtime extension for PHP */
#include "Includes.h"
#include "Utils.h"
#include "Handle.h"

ZEND_DECLARE_MODULE_GLOBALS(aikido)

PHP_INI_BEGIN()
	STD_PHP_INI_ENTRY("aikido.log_level", "-1", PHP_INI_ALL, OnUpdateLong, log_level, zend_aikido_globals, aikido_globals)
	STD_PHP_INI_ENTRY("aikido.token", "AIK_RUNTIME_UNSET", PHP_INI_ALL, OnUpdateString, token, zend_aikido_globals, aikido_globals)
	STD_PHP_INI_ENTRY("aikido.blocking", "0", PHP_INI_ALL, OnUpdateBool, blocking, zend_aikido_globals, aikido_globals)
PHP_INI_END()

bool aikido_global_init() {
	if (AIKIDO_GLOBAL(log_level) < AIKIDO_LOG_LEVEL_DEBUG ||
		AIKIDO_GLOBAL(log_level) > AIKIDO_LOG_LEVEL_ERROR) {
		AIKIDO_GLOBAL(log_level) = AIKIDO_LOG_LEVEL_ERROR;
	}

	const char* log_level_str = aikido_log_level_str((AIKIDO_LOG_LEVEL)AIKIDO_GLOBAL(log_level));
	
	AIKIDO_LOG_DEBUG("Config:\n");
	AIKIDO_LOG_DEBUG("Log level: %s\n", log_level_str);
	AIKIDO_LOG_DEBUG("Blocking: %d\n", AIKIDO_GLOBAL(blocking));

	json initData = {
		{ "log_level", log_level_str },
		{ "token", AIKIDO_GLOBAL(token) },
		{ "blocking", AIKIDO_GLOBAL(blocking) }
	};

	return GoInit(initData);
}

static PHP_GINIT_FUNCTION(aikido)
{
#if defined(COMPILE_DL_BCMATH) && defined(ZTS)
	ZEND_TSRMLS_CACHE_UPDATE();
#endif
}

PHP_MINIT_FUNCTION(aikido)
{
	REGISTER_INI_ENTRIES();

	if (!aikido_global_init()) {
		return FAILURE;
	}

	for ( auto& it : HOOKED_FUNCTIONS ) {
		zend_function* function_data = (zend_function*)zend_hash_str_find_ptr(CG(function_table), it.first.c_str(), it.first.length());
		if (function_data != NULL) {
			it.second.original_handler = function_data->internal_function.handler;
			function_data->internal_function.handler = aikido_generic_handler;
			AIKIDO_LOG_DEBUG("Hooked function \"%s\" (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
		}
	}

	for ( auto& it : HOOKED_METHODS ) {
		zend_class_entry *class_entry = (zend_class_entry *)zend_hash_str_find_ptr(CG(class_table), it.first.class_name.c_str(), it.first.class_name.length());
		if (class_entry != NULL) {
			zend_function *method = (zend_function*)zend_hash_str_find_ptr(&class_entry->function_table, it.first.method_name.c_str(), it.first.method_name.length());
			if (method != NULL) {
				it.second.original_handler = method->internal_function.handler;
				method->internal_function.handler = aikido_generic_handler;
				AIKIDO_LOG_DEBUG("Hooked method \"%s->%s\" (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
			}
   		}
	}

	return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(aikido)
{
	UNREGISTER_INI_ENTRIES();
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
	NULL,						/* PHP_RINIT - Request initialization */
	NULL,						/* PHP_RSHUTDOWN - Request shutdown */
	PHP_MINFO(aikido),			/* PHP_MINFO - Module info */
	PHP_AIKIDO_VERSION,			/* Version */
	PHP_MODULE_GLOBALS(aikido),	/* Module globals */
	PHP_GINIT(aikido),			/* PHP_GINIT – Globals initialization */
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
