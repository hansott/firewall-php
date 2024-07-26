#include "HandleUrls.h"
#include "Utils.h"


AIKIDO_HANDLER_FUNCTION(handle_curl_exec) {
    zval *curlHandle = NULL;
    #if PHP_VERSION_ID >= 80000
        ZEND_PARSE_PARAMETERS_START(1, 1)
            Z_PARAM_OBJECT(curlHandle)
            
        ZEND_PARSE_PARAMETERS_END();
    #else
        ZEND_PARSE_PARAMETERS_START(1, 1)
            Z_PARAM_RESOURCE(curlHandle)
        ZEND_PARSE_PARAMETERS_END();
    #endif

	zval retval;
	zval params[2];
	ZVAL_COPY(&params[0], curlHandle);
	ZVAL_LONG(&params[1], CURLINFO_EFFECTIVE_URL);
	zval* fname = NULL;

	fname = (zval*)emalloc(sizeof(zval));

	if (fname == NULL) {
		return;
	}
		
	ZVAL_STRING(fname, "curl_getinfo");

	if (call_user_function(EG(function_table), NULL, fname, &retval, 2, params) == SUCCESS) {
		if (Z_TYPE(retval) == IS_STRING) {
			std::string urlString(Z_STRVAL(retval));
			inputEvent = {
				{ "event", "function_executed" },
				{ "data", {
					{ "function_name", "curl_exec" },
					{ "parameters", {
						{ "url", urlString }
					} }
				} }
			};
		}
	}

	zval_dtor(&retval);
	zval_dtor(&params[0]);
	zval_dtor(&params[1]);
	efree(fname);     
}
