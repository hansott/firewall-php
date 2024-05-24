#include "HandlePDO.h"

ZEND_NAMED_FUNCTION(handle_pdo___construct) {
	AIKIDO_METHOD_HANDLER_START();

	char *data_source;
	size_t data_source_len;
	char *colon;
	char *username=NULL, *password=NULL;
	size_t usernamelen, passwordlen;

	ZEND_PARSE_PARAMETERS_START(1, 4)
		Z_PARAM_STRING(data_source, data_source_len)
		Z_PARAM_OPTIONAL
		Z_PARAM_STRING_OR_NULL(username, usernamelen)
		Z_PARAM_STRING_OR_NULL(password, passwordlen)
		//Z_PARAM_ARRAY_OR_NULL(options)
	ZEND_PARSE_PARAMETERS_END();

	std::string data_source_string(data_source, data_source_len);

	json pdo_construct_event = {
		{ "event", "method_executed" },
		{ "data", {
			{ "class_name", "pdo" },
			{ "method_name", "__construct" },
			{ "parameters", {
				{ "data_source", data_source_string }
			} }
		} }
	};

	if (username) {
		std::string usernameString(username, usernamelen);
		pdo_construct_event["data"]["parameters"]["username"] = usernameString;
	}
	if (password) {
		std::string passwordString(password, passwordlen);
		pdo_construct_event["data"]["parameters"]["password"] = passwordString;
	}

	GoOnEvent(pdo_construct_event);

	AIKIDO_METHOD_HANDLER_END();
}

ZEND_NAMED_FUNCTION(handle_pdo_query) {
	AIKIDO_METHOD_HANDLER_START();

	char *query;
	size_t query_len;

	ZEND_PARSE_PARAMETERS_START(1,-1)
		Z_PARAM_STRING(query, query_len)
	ZEND_PARSE_PARAMETERS_END();

	std::string query_string(query, query_len);

	json pdo_query = {
		{ "event", "method_executed" },
		{ "data", {
			{ "class_name", "pdo" },
			{ "method_name", "query" },
			{ "parameters", {
				{ "query",  query_string }
			} }
		} }
	};

	GoOnEvent(pdo_query);

	AIKIDO_METHOD_HANDLER_END();
}
