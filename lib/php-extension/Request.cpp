#include "Includes.h"

Request request;

bool Request::Init() {
    zend_string* serverString = zend_string_init("_SERVER", sizeof("_SERVER") - 1, 0);
    if (!serverString) {
        return false;
    }

    /* Guarantee that "_SERVER" PHP global variable is initialized for the current request */
    if (!zend_is_auto_global(serverString)) {
        zend_string_release(serverString);
        return false;
    }

    zend_string_release(serverString);

    /* Search for the "_SERVER" PHP global variable and store it */
    this->server = zend_hash_str_find(&EG(symbol_table), "_SERVER", sizeof("_SERVER") - 1);
    return this->server != NULL;
}

bool Request::Ok() {
    return this->server != NULL;
}

std::string Request::GetVar(const char* var) {
    if (!this->server) {
        return "";
    }
    zval* data = zend_hash_str_find(Z_ARRVAL_P(this->server), var, strlen(var));
    if (!data) {
        return "";
    }
    return Z_STRVAL_P(data);
}

std::string Request::GetRoute() {
    if (!this->server) {
        return "";
    }
    std::string route = GetVar("REQUEST_URI");
    size_t pos = route.find("?");
    if (pos != std::string::npos) {
        route = route.substr(0, pos);
    }
    return route;
}

std::string Request::GetStatusCode() {
    return std::to_string(SG(sapi_headers).http_response_code);
}

std::string Request::GetUrl() {
    return (IsHttps() ? "https://" : "http://") + GetVar("HTTP_HOST") + GetVar("REQUEST_URI");
}

std::string Request::GetBody() {
    long maxlen = PHP_STREAM_COPY_ALL;
    zend_string* contents;
    php_stream* stream;

    stream = php_stream_open_wrapper("php://input", "rb", 0 | REPORT_ERRORS, NULL);
    if ((contents = php_stream_copy_to_mem(stream, maxlen, 0)) != NULL) {
        php_stream_close(stream);
        return std::string(ZSTR_VAL(contents));
    }
    php_stream_close(stream);
    return "";
}

std::string Request::GetHeaders() {
    if (!this->server) {
        return "";
    }
    std::map<std::string, std::string> headers;
    zend_string* key;
    zval* val;
    ZEND_HASH_FOREACH_STR_KEY_VAL(Z_ARRVAL_P(this->server), key, val) {
        if (key && Z_TYPE_P(val) == IS_STRING) {
            std::string header_name(ZSTR_VAL(key));
            std::string http_header_key;
            std::string http_header_value(Z_STRVAL_P(val));

            if (header_name.find("HTTP_") == 0) {
                http_header_key = header_name.substr(5);
            } else if (header_name == "CONTENT_TYPE" || header_name == "CONTENT_LENGTH" || header_name == "AUTHORIZATION") {
                http_header_key = header_name;
            }

            if (!http_header_key.empty()) {
                std::transform(http_header_key.begin(), http_header_key.end(), http_header_key.begin(), ::tolower);
                headers[http_header_key] = http_header_value;
            }
        }
    }
    ZEND_HASH_FOREACH_END();

    json headers_json;
    for (auto const& [key, val] : headers) {
        headers_json[key] = val;
    }
    return headers_json.dump();
}

bool Request::IsHttps() {
    return GetVar("HTTPS") != "" ? true : false;
}
