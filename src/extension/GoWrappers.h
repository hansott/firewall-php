#pragma once

#include "libaikido_go.h"

#include <string>

std::string CreateCString(GoString g);

GoString CreateGoString(std::string& s);

std::string GetHostname(std::string& url);