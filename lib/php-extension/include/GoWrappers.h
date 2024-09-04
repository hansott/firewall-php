#pragma once

#include "Includes.h"

std::string CppCreateString(GoString g);

GoString GoCreateString(std::string& s);

json GoRequestProcessorOnEvent(json& event);

char* GoContextCallback(int, char*);

bool GoRequestProcessorContextInit();

bool IsBlockingEnabled();