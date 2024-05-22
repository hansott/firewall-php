#pragma once

#include "libaikido_go.h"

#include <string>
#include "3rdparty/json.hpp"

using json = nlohmann::json;

std::string CppCreateString(GoString g);

GoString GoCreateString(std::string& s);

json GoOnEvent(json& event);