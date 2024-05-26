#include "HandlePDO.h"

ZEND_NAMED_FUNCTION(handle_pdo___construct) {
	AIKIDO_METHOD_HANDLER_START();

	char *data_source;
	size_t data_source_len;
	char *colon;

	ZEND_PARSE_PARAMETERS_START(1, -1)
		Z_PARAM_STRING(data_source, data_source_len)
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
