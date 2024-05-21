/* Aikido runtime extension for PHP */

#ifdef HAVE_CONFIG_H
# include "config.h"
#endif

#include "php.h"
#include "ext/standard/info.h"
#include "php_aikido.h"
#include <unordered_map>
#include <set>
#include <string>
#include <curl/curl.h>
#include "libaikido_go.h"

#include "GoWrappers.h"

using namespace std;

ZEND_NAMED_FUNCTION(handle_file_get_contents);
ZEND_NAMED_FUNCTION(handle_curl_init);
ZEND_NAMED_FUNCTION(handle_curl_setopt);

struct FUNCTION_HANDLERS {
	zif_handler aikido_handler;
	zif_handler original_handler;
};

#define AIKIDO_REGISTER_HANDLER(function_name) { #function_name, { handle_##function_name, nullptr } }

unordered_map<const char*, FUNCTION_HANDLERS> HOOKED_FUNCTIONS = {
	AIKIDO_REGISTER_HANDLER(file_get_contents),
	AIKIDO_REGISTER_HANDLER(curl_init),
	AIKIDO_REGISTER_HANDLER(curl_setopt)
};

unordered_map<void*, string> curlHandlers;
set<string> outgoingHostnames;

#define AIKIDO_HANDLER_START(function_name) php_printf("[AIKIDO] Handler called for \"" #function_name "\"!\n");
#define AIKIDO_HANDLER_END(function_name) HOOKED_FUNCTIONS[#function_name].original_handler(INTERNAL_FUNCTION_PARAM_PASSTHRU);

ZEND_NAMED_FUNCTION(handle_file_get_contents) {
	AIKIDO_HANDLER_START(file_get_contents);
	AIKIDO_HANDLER_END(file_get_contents);
}

ZEND_NAMED_FUNCTION(handle_curl_init) {
	AIKIDO_HANDLER_START(curl_init);

	zend_string *url = NULL;

	ZEND_PARSE_PARAMETERS_START(0,1)
		Z_PARAM_OPTIONAL
		Z_PARAM_STR_OR_NULL(url)
	ZEND_PARSE_PARAMETERS_END();

	AIKIDO_HANDLER_END(curl_init);
	
	if (Z_TYPE_P(return_value) != IS_FALSE) {
		// Z_OBJ_P(return_value)
		if (url) {
			std::string urlString(ZSTR_VAL(url));
			outgoingHostnames.insert(GetHostname(urlString));
		}
	}
}

ZEND_NAMED_FUNCTION(handle_curl_setopt) {
	AIKIDO_HANDLER_START(curl_setopt);

	zval *curlHandle = NULL;
	zend_long options = 0;
	zval *zvalue = NULL;

	ZEND_PARSE_PARAMETERS_START(3, 3)
		Z_PARAM_OBJECT(curlHandle)
		Z_PARAM_LONG(options)
		Z_PARAM_ZVAL(zvalue)
	ZEND_PARSE_PARAMETERS_END();

	if (options == CURLOPT_URL) {
		zend_string *tmp_str;
		zend_string *url = zval_get_tmp_string(zvalue, &tmp_str);

		std::string urlString(ZSTR_VAL(url));
	
		outgoingHostnames.insert(GetHostname(urlString));
	
		zend_tmp_string_release(tmp_str);
	}

	AIKIDO_HANDLER_END(curl_setopt);
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

PHP_MSHUTDOWN_FUNCTION(aikido)
{
	php_printf("[AIKIDO] List of outgoing hostnames:\n");
	for (auto hostname: outgoingHostnames) {
		php_printf("[AIKIDO] - %s\n", hostname.c_str());
	}
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
