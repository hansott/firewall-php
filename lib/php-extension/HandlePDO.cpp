#include "HandlePDO.h"
#include "Utils.h"
#include "Cache.h"

AIKIDO_HANDLER_FUNCTION(handle_pre_pdo_query) {
	zend_string *query = NULL;

	ZEND_PARSE_PARAMETERS_START(0,-1)
		Z_PARAM_OPTIONAL
		Z_PARAM_STR(query)
	ZEND_PARSE_PARAMETERS_END();

	if (!query) {
		return;
	}

	zval *pdo_object = getThis();
	if (!pdo_object) {
		return;
	}

	eventId = EVENT_PRE_EXECUTED_PDO_QUERY;
	eventCache.sql_query = ZSTR_VAL(query);
	eventCache.sql_dialect = "unknown";
	
	zval retval;
	if (aikido_call_user_function_one_param("getAttribute", PDO_ATTR_DRIVER_NAME, &retval, pdo_object)) {
		if (Z_TYPE(retval) == IS_STRING)
		{
			eventCache.sql_dialect = Z_STRVAL_P(&retval);
		}
    }

    zval_ptr_dtor(&retval);	
}
