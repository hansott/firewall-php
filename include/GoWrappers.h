#pragma once

#include "Includes.h"

std::string CppCreateString(GoString g);

GoString GoCreateString(std::string& s);

json GoOnEvent(json& event);

bool GoInit(json& initData);
