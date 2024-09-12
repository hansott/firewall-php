#pragma once

#include "Includes.h"

std::string CppCreateString(GoString);

GoString GoCreateString(std::string&);

bool GoRequestProcessorOnEvent(EVENT_ID, std::string&);

char *GoContextCallback(int, char *);

bool GoRequestProcessorContextInit();

bool IsBlockingEnabled();