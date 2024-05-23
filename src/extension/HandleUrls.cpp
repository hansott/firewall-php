#include "HandleUrls.h"
#include "Utils.h"

ZEND_NAMED_FUNCTION(handle_curl_init) {
	AIKIDO_HANDLER_START();

	zend_string *url = NULL;

	ZEND_PARSE_PARAMETERS_START(0,1)
		Z_PARAM_OPTIONAL
		Z_PARAM_STR_OR_NULL(url)
	ZEND_PARSE_PARAMETERS_END();

	AIKIDO_HANDLER_END();
	
	if (Z_TYPE_P(return_value) != IS_FALSE) {
		// Z_OBJ_P(return_value)
		json curl_init_event = {
			{ "event", "function_executed" },
			{ "data", {
				{ "function_name", "curl_init" },
				{ "parameters", json::object() }
			} }
		};
		if (url) {
			std::string urlString(ZSTR_VAL(url));
			curl_init_event["data"]["parameters"]["url"] = urlString;
		}
		GoOnEvent(curl_init_event);
	}
}

ZEND_NAMED_FUNCTION(handle_curl_setopt) {
	AIKIDO_HANDLER_START();

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
	
		json curl_setopt_event = {
			{ "event", "function_executed" },
			{ "data", {
				{ "function_name", "curl_setopt" },
				{ "parameters", {
					{ "url", urlString }
				} }
			} }
		};

		GoOnEvent(curl_setopt_event);

		zend_tmp_string_release(tmp_str);
	}

	AIKIDO_HANDLER_END();
}
