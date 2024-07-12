#pragma once

#include "Includes.h"

std::string CppCreateString(GoString g);

GoString GoCreateString(std::string& s);

json GoRequestProcessorOnEvent(json& event);

bool GoRequestProcessorInit(json& initData);

void GoRequestProcessorUninit();

bool GoAgentInit(json& initData);

void GoAgentUninit();
