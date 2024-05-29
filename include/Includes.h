#pragma once

#ifdef HAVE_CONFIG_H
# include "config.h"
#endif

#include "php.h"
#include "zend_exceptions.h"

#include <string>
#include <unordered_map>
#include <curl/curl.h>
#include <functional>

#include "ext/standard/info.h"
#include "libaikido_go.h"
#include "php_aikido.h"

#include "3rdparty/json.hpp"

using namespace std;
using json = nlohmann::json;

/* For compatibility with older PHP versions */
#ifndef ZEND_PARSE_PARAMETERS_NONE
#define ZEND_PARSE_PARAMETERS_NONE() \
	ZEND_PARSE_PARAMETERS_START(0, 0) \
	ZEND_PARSE_PARAMETERS_END()
#endif

#include "GoWrappers.h"
#include "Utils.h"

enum AIKIDO_LOG_LEVEL {
    AIKIDO_LOG_LEVEL_DEBUG,
    AIKIDO_LOG_LEVEL_INFO,
    AIKIDO_LOG_LEVEL_WARN,
    AIKIDO_LOG_LEVEL_ERROR
};

void aikido_log(AIKIDO_LOG_LEVEL level, const char* format, ...);
const char* aikido_log_level_str(AIKIDO_LOG_LEVEL level);

#define AIKIDO_LOG_DEBUG(format, ...)  aikido_log(AIKIDO_LOG_LEVEL_DEBUG, format, ##__VA_ARGS__)
#define AIKIDO_LOG_INFO(format, ...)   aikido_log(AIKIDO_LOG_LEVEL_INFO, format, ##__VA_ARGS__)
#define AIKIDO_LOG_WARN(format, ...)   aikido_log(AIKIDO_LOG_LEVEL_WARN, format, ##__VA_ARGS__)
#define AIKIDO_LOG_ERROR(format, ...)  aikido_log(AIKIDO_LOG_LEVEL_ERROR, format, ##__VA_ARGS__)

#define AIKIDO_HANDLER_FUNCTION(name) void name(INTERNAL_FUNCTION_PARAMETERS, json& inputEvent)

typedef void (*aikido_handler)(INTERNAL_FUNCTION_PARAMETERS, json& inputEvent);

struct PHP_HANDLERS {
	aikido_handler handler;
	zif_handler original_handler;
};

extern unordered_map<std::string, PHP_HANDLERS> HOOKED_FUNCTIONS;

class AIKIDO_METHOD_KEY {
public:
    std::string class_name;
    std::string method_name;

    // Constructor
    AIKIDO_METHOD_KEY(const std::string& class_name, const std::string& method_name) : class_name(class_name), method_name(method_name) {}

    // Equality operator
    bool operator==(const AIKIDO_METHOD_KEY& other) const {
        return class_name == other.class_name && method_name == other.method_name;
    }
};

class AIKIDO_METHOD_KEY_HASH {
public:
    std::size_t operator()(const AIKIDO_METHOD_KEY& k) const {
        return std::hash<std::string>()(k.class_name) ^ (std::hash<std::string>()(k.method_name) << 1);
    }
};

extern unordered_map<AIKIDO_METHOD_KEY, PHP_HANDLERS, AIKIDO_METHOD_KEY_HASH> HOOKED_METHODS;
