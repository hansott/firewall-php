#include "HandleUrls.h"
#include "Utils.h"

AIKIDO_HANDLER_FUNCTION(handle_curl_init) {
	zend_string *url = NULL;

	ZEND_PARSE_PARAMETERS_START(0,1)
		Z_PARAM_OPTIONAL
		Z_PARAM_STR(url)
	ZEND_PARSE_PARAMETERS_END();

	// Z_OBJ_P(return_value)
	inputEvent = {
		{ "event", "function_executed" },
		{ "data", {
			{ "function_name", "curl_init" },
			{ "parameters", json::object() }
		} }
	};
	
	if (url) {
		std::string urlString(ZSTR_VAL(url));
		inputEvent["data"]["parameters"]["url"] = urlString;
	}
}

AIKIDO_HANDLER_FUNCTION(handle_curl_setopt) {
	zval *curlHandle = NULL;
	zend_long options = 0;
	zval *zvalue = NULL;

	AIKIDO_LOG_DEBUG("Handling curl setopt!");

	#if PHP_VERSION_ID >= 80000
		ZEND_PARSE_PARAMETERS_START(3, 3)
			Z_PARAM_OBJECT_OF_CLASS(curlHandle)
			Z_PARAM_LONG(options)
			Z_PARAM_ZVAL(zvalue)
		ZEND_PARSE_PARAMETERS_END();
	#else
		ZEND_PARSE_PARAMETERS_START(3, 3)
			Z_PARAM_RESOURCE(curlHandle)
			Z_PARAM_LONG(options)
			Z_PARAM_ZVAL(zvalue)
		ZEND_PARSE_PARAMETERS_END();
	#endif

	AIKIDO_LOG_DEBUG("Handling curl setopt: options = %d, CURLOPT_URL = %d\n", options, CURLOPT_URL);

	if (options == CURLOPT_URL) {
		zend_string *tmp_str;
		zend_string *url = zval_get_tmp_string(zvalue, &tmp_str);

		std::string urlString(ZSTR_VAL(url));
	
		inputEvent = {
			{ "event", "function_executed" },
			{ "data", {
				{ "function_name", "curl_setopt" },
				{ "parameters", {
					{ "url", urlString }
				} }
			} }
		};

		zend_tmp_string_release(tmp_str);
	}
}
