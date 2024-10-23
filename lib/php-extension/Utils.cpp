#include "Includes.h"

std::string ToLowercase(const std::string& str) {
    std::string result = str;
    std::transform(result.begin(), result.end(), result.begin(), [](unsigned char c) { return std::tolower(c); });
    return result;
}

FILE* log_file = nullptr;

void aikido_log_init() {
    std::time_t current_time = std::time(nullptr);
    char time_str[20];
    std::strftime(time_str, sizeof(time_str), "%Y%m%d%H%M%S", std::localtime(&current_time));
    std::string log_file_path = "/var/log/aikido-" + std::string(PHP_AIKIDO_VERSION) + "/aikido-extension-php-" + time_str + "-" + std::to_string(getpid()) + ".log";
    log_file = fopen(log_file_path.c_str(), "w");
}

void aikido_log_uninit() {
    if (log_file) {
        fclose(log_file);
        log_file = nullptr;
    }
}

void aikido_log(AIKIDO_LOG_LEVEL level, const char* format, ...) {
    if (!log_file || level < AIKIDO_GLOBAL(log_level)) {
        return;
    }

    fprintf(log_file, "[AIKIDO][%s] ", aikido_log_level_str(level).c_str());

    va_list args;
    va_start(args, format);
    vfprintf(log_file, format, args);
    va_end(args);

    fflush(log_file);
}

std::string aikido_log_level_str(AIKIDO_LOG_LEVEL level) {
    switch (level) {
        case AIKIDO_LOG_LEVEL_DEBUG:
            return "DEBUG";
        case AIKIDO_LOG_LEVEL_INFO:
            return "INFO";
        case AIKIDO_LOG_LEVEL_WARN:
            return "WARN";
        case AIKIDO_LOG_LEVEL_ERROR:
            return "ERROR";
    }
    return "UNKNOWN";
}

AIKIDO_LOG_LEVEL aikido_log_level_from_str(std::string level) {
    if (level == "ERROR") {
        return AIKIDO_LOG_LEVEL_ERROR;
    }
    if (level == "WARN") {
        return AIKIDO_LOG_LEVEL_WARN;
    }
    if (level == "INFO") {
        return AIKIDO_LOG_LEVEL_INFO;
    }
    if (level == "DEBUG") {
        return AIKIDO_LOG_LEVEL_DEBUG;
    }
    return AIKIDO_LOG_LEVEL_ERROR;
}

bool aikido_echo(std::string message) {
    unsigned int wrote = zend_write(message.c_str(), message.length());  // echo '<message>'
    AIKIDO_LOG_INFO("Called 'echo' -> result %d\n", wrote == message.length());
    return wrote == message.length();
}

bool aikido_call_user_function(std::string function_name, unsigned int params_number, zval* params, zval* return_value, zval* object) {
    zval _function_name;
    zend_string* _function_name_str = zend_string_init(function_name.c_str(), function_name.length(), 0);
    if (!_function_name_str) {
        return false;
    }
    ZVAL_STR(&_function_name, _function_name_str);

    zval* _return_value = return_value;
    zval _temp_return_value;
    if (!return_value) {
        _return_value = &_temp_return_value;
    }

    int _result = call_user_function(EG(function_table), object, &_function_name, _return_value, params_number, params);

    zend_string_release(_function_name_str);

    if (!return_value) {
        zval_ptr_dtor(&_temp_return_value);
    }
    AIKIDO_LOG_INFO("Called user function '%s' -> result %d\n", function_name.c_str(), _result == SUCCESS);
    return _result == SUCCESS;
}

bool aikido_call_user_function_one_param(std::string function_name, long first_param, zval* return_value, zval* object) {
    zval _params[1];
    ZVAL_LONG(&_params[0], first_param);
    return aikido_call_user_function(function_name, 1, _params, return_value, object);
}

bool aikido_call_user_function_one_param(std::string function_name, std::string first_param, zval* return_value, zval* object) {
    zval _params[1];
    zend_string* _first_param = zend_string_init(first_param.c_str(), first_param.length(), 0);
    if (!_first_param) {
        return false;
    }
    ZVAL_STR(&_params[0], _first_param);

    bool ret = aikido_call_user_function(function_name, 1, _params, return_value, object);

    zend_string_release(_first_param);

    return ret;
}

int GetRandomNumber() {
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dis(100000, 999999);
    return dis(gen);
}

std::string GenerateSocketPath() {
    std::time_t current_time = std::time(nullptr);
    char time_str[20];
    std::strftime(time_str, sizeof(time_str), "%Y%m%d%H%M%S", std::localtime(&current_time));
    std::string socket_file_path = "/run/aikido-" + std::string(PHP_AIKIDO_VERSION) + "/aikido-" +
                                   std::string(time_str) + "-" + std::to_string(GetRandomNumber()) + ".sock";
    return socket_file_path;
}

const char* GetEventName(EVENT_ID event) {
    switch (event) {
        case EVENT_PRE_REQUEST:
            return "PreRequest";
        case EVENT_POST_REQUEST:
            return "PostRequest";
        case EVENT_SET_USER:
            return "SetUser";
        case EVENT_GET_BLOCKING_STATUS:
            return "GetBlockingStatus";
        case EVENT_PRE_OUTGOING_REQUEST:
            return "PreOutgoingRequest";
        case EVENT_POST_OUTGOING_REQUEST:
            return "PostOutgoingRequest";
        case EVENT_PRE_SHELL_EXECUTED:
            return "PreShellExecuted";
        case EVENT_PRE_PATH_ACCESSED:
            return "PrePathAccessed";
        case EVENT_PRE_SQL_QUERY_EXECUTED:
            return "PreSqlQueryExecuted";
    }
    return "Unknown";
}

std::string aikido_call_user_function_curl_getinfo(zval* curl_handle, int curl_info_option) {
    zval retval;
    zval params[2];

    ZVAL_COPY(&params[0], curl_handle);
    ZVAL_LONG(&params[1], curl_info_option);

    std::string result = "";
    if (aikido_call_user_function("curl_getinfo", 2, params, &retval)) {
        switch (Z_TYPE(retval)) {
            case IS_LONG:
                result = std::to_string(Z_LVAL(retval));
                break;
            case IS_STRING:
                result = Z_STRVAL(retval);
                break;
        }
    }

    zval_dtor(&params[0]);
    zval_dtor(&params[1]);
    zval_dtor(&retval);

    return result;
}