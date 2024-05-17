/* Aikido runtime extension for PHP */

#ifdef HAVE_CONFIG_H
# include "config.h"
#endif

#include "php.h"
#include "ext/standard/info.h"
#include "php_aikido.h"
#include <unordered_map>

ZEND_NAMED_FUNCTION(handle_file_get_contents);

struct FUNCTION_HANDLERS {
	zif_handler aikido_handler;
	zif_handler original_handler;
};

#define AIKIDO_REGISTER_HANDLER(function_name) { #function_name, { handle_##function_name, nullptr } }

std::unordered_map<const char*, FUNCTION_HANDLERS> HOOKED_FUNCTIONS = {
	AIKIDO_REGISTER_HANDLER(file_get_contents)
};

#define AIKIDO_HANDLER_START(function_name) php_printf("[AIKIDO] Handler called for \"" #function_name "\"!\n");
#define AIKIDO_HANDLER_END(function_name) HOOKED_FUNCTIONS[#function_name].original_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU);

ZEND_NAMED_FUNCTION(handle_file_get_contents) {
	AIKIDO_HANDLER_START(file_get_contents);
	AIKIDO_HANDLER_END(file_get_contents);
}

/* For compatibility with older PHP versions */
#ifndef ZEND_PARSE_PARAMETERS_NONE
#define ZEND_PARSE_PARAMETERS_NONE() \
	ZEND_PARSE_PARAMETERS_START(0, 0) \
	ZEND_PARSE_PARAMETERS_END()
#endif

PHP_MINIT_FUNCTION(aikido)
{
#if defined(COMPILE_DL_MY_EXTENSION) && defined(ZTS)
    ZEND_TSRMLS_CACHE_UPDATE();
#endif

	for ( auto& it : HOOKED_FUNCTIONS ) {
		zend_function* function_data = (zend_function*)zend_hash_str_find_ptr(CG(function_table), it.first, strlen(it.first));
		if (function_data != NULL) {
        	it.second.original_handler = function_data->internal_function.handler;
        	function_data->internal_function.handler = it.second.aikido_handler;
			php_printf("[AIKIDO] Hooked function \"%s\" using aikido handler %p (original handler %p)!\n", it.first, it.second.aikido_handler, it.second.original_handler);
    	}
	}

	return SUCCESS;
}

PHP_MSHUTDOWN_FUNCTION(my_extension)
{
    return SUCCESS;
}

PHP_RINIT_FUNCTION(aikido)
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
	ext_functions,			    /* zend_function_entry */
	PHP_MINIT(aikido),			/* PHP_MINIT - Module initialization */
	NULL,						/* PHP_MSHUTDOWN - Module shutdown */
	PHP_RINIT(aikido),			/* PHP_RINIT - Request initialization */
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
