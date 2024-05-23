#pragma once

#ifdef HAVE_CONFIG_H
# include "config.h"
#endif

#include "php.h"

#include <string>
#include <unordered_map>
#include <curl/curl.h>


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

struct PHP_HANDLERS {
	zif_handler aikido_handler;
	zif_handler original_handler;
};

extern unordered_map<std::string, PHP_HANDLERS> HOOKED_FUNCTIONS;
extern unordered_map<std::string, PHP_HANDLERS> HOOKED_METHODS;
