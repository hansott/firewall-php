/* Aikido runtime extension for PHP */
#include "Includes.h"
#include "Utils.h"
#include "Handle.h"

PHP_MINIT_FUNCTION(aikido)
{
#if defined(COMPILE_DL_MY_EXTENSION) && defined(ZTS)
	ZEND_TSRMLS_CACHE_UPDATE();
#endif

	for ( auto& it : HOOKED_FUNCTIONS ) {
		zend_function* function_data = (zend_function*)zend_hash_str_find_ptr(CG(function_table), it.first.c_str(), it.first.length());
		if (function_data != NULL) {
			it.second.original_handler = function_data->internal_function.handler;
			function_data->internal_function.handler = aikido_generic_handler;
			php_printf("[AIKIDO-C++] Hooked function \"%s\" (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
		}
	}

	for ( auto& it : HOOKED_METHODS ) {
		zend_class_entry *class_entry = (zend_class_entry *)zend_hash_str_find_ptr(CG(class_table), it.first.class_name.c_str(), it.first.class_name.length());
		if (class_entry != NULL) {
			zend_function *method = (zend_function*)zend_hash_str_find_ptr(&class_entry->function_table, it.first.method_name.c_str(), it.first.method_name.length());
			if (method != NULL) {
				it.second.original_handler = method->internal_function.handler;
				method->internal_function.handler = aikido_generic_handler;
				php_printf("[AIKIDO-C++] Hooked method \"%s->%s\" (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
			}
   		}
	}

	return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(aikido)
{
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
	STANDARD_MODULE_PROPERTIES
};

#ifdef COMPILE_DL_AIKIDO
# ifdef ZTS
ZEND_TSRMLS_CACHE_DEFINE()
# endif
ZEND_GET_MODULE(aikido)
#endif
