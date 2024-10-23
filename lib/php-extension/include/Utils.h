#pragma once

enum AIKIDO_LOG_LEVEL {
    AIKIDO_LOG_LEVEL_DEBUG,
    AIKIDO_LOG_LEVEL_INFO,
    AIKIDO_LOG_LEVEL_WARN,
    AIKIDO_LOG_LEVEL_ERROR
};

void aikido_log_init();

void aikido_log_uninit();

void aikido_log(AIKIDO_LOG_LEVEL level, const char *format, ...);

#if defined(ZEND_DEBUG)
#define AIKIDO_LOG_DEBUG(format, ...) aikido_log(AIKIDO_LOG_LEVEL_DEBUG, format, ##__VA_ARGS__)
#else
/* Disable debugging logs for production builds */
#define AIKIDO_LOG_DEBUG(format, ...)
#endif

#define AIKIDO_LOG_INFO(format, ...) aikido_log(AIKIDO_LOG_LEVEL_INFO, format, ##__VA_ARGS__)
#define AIKIDO_LOG_WARN(format, ...) aikido_log(AIKIDO_LOG_LEVEL_WARN, format, ##__VA_ARGS__)
#define AIKIDO_LOG_ERROR(format, ...) aikido_log(AIKIDO_LOG_LEVEL_ERROR, format, ##__VA_ARGS__)

std::string aikido_log_level_str(AIKIDO_LOG_LEVEL level);

AIKIDO_LOG_LEVEL aikido_log_level_from_str(std::string level);

std::string to_lowercase(const std::string &str);

bool aikido_echo(std::string s);

bool aikido_call_user_function(std::string function_name, unsigned int params_number = 0, zval *params = nullptr, zval *return_value = nullptr, zval *object = nullptr);

bool aikido_call_user_function_one_param(std::string function_name, long first_param, zval *return_value = nullptr, zval *object = nullptr);

bool aikido_call_user_function_one_param(std::string function_name, std::string first_param, zval *return_value = nullptr, zval *object = nullptr);

std::string aikido_call_user_function_curl_getinfo(zval *curl_handle, int curl_info_option);

std::string aikido_generate_socket_path();

const char *GetEventName(EVENT_ID event);
