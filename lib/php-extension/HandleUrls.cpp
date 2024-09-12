#include "HandleUrls.h"
#include "Utils.h"
#include "Cache.h"

AIKIDO_HANDLER_FUNCTION(handle_pre_curl_exec) {
    zval *curlHandle = NULL;

	// Curl handles changed between PHP 7 & PHP 8 - so we need different extraction
    #if PHP_VERSION_ID >= 80000
        ZEND_PARSE_PARAMETERS_START(1, 1)
            Z_PARAM_OBJECT(curlHandle) 
        ZEND_PARSE_PARAMETERS_END();
    #else
        ZEND_PARSE_PARAMETERS_START(1, 1)
            Z_PARAM_RESOURCE(curlHandle)
        ZEND_PARSE_PARAMETERS_END();
    #endif

	// Prepare parameters for curl_getinfo php call
	zval retval;
	zval params[2];
	ZVAL_COPY(&params[0], curlHandle);
	ZVAL_LONG(&params[1], CURLINFO_EFFECTIVE_URL);

	// Call curl_getinfo to extract the URL
	if (!aikido_call_user_function("curl_getinfo", 2, params, &retval)) return;

	eventId = EVENT_PRE_OUTGOING_REQUEST;
	eventCache.outgoingRequestUrl = Z_STRVAL(retval);

	zval_dtor(&retval);
	zval_dtor(&params[0]);
	zval_dtor(&params[1]);
}


AIKIDO_HANDLER_FUNCTION(handle_post_curl_exec) {
	zval *curlHandle = NULL;

	// Curl handles changed between PHP 7 & PHP 8 - so we need different extraction
	#if PHP_VERSION_ID >= 80000
		ZEND_PARSE_PARAMETERS_START(1, 1)
			Z_PARAM_OBJECT(curlHandle) 
		ZEND_PARSE_PARAMETERS_END();
	#else
		ZEND_PARSE_PARAMETERS_START(1, 1)
			Z_PARAM_RESOURCE(curlHandle)
		ZEND_PARSE_PARAMETERS_END();
	#endif

	// Prepare parameters for curl_getinfo php call
	zval retval;
	zval params[2];
	ZVAL_COPY(&params[0], curlHandle);
	ZVAL_LONG(&params[1], CURLINFO_PRIMARY_PORT);

	// Call curl_getinfo to extract the PORT
	if (!aikido_call_user_function("curl_getinfo", 2, params, &retval)) return;

	eventId = EVENT_POST_OUTGOING_REQUEST;
	eventCache.outgoingRequestPort = Z_LVAL(retval);
	
	zval_dtor(&retval);
	zval_dtor(&params[0]);
	zval_dtor(&params[1]);
}
