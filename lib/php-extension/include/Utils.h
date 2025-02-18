#pragma once

#include "Includes.h"

std::string ToLowercase(const std::string& str);

std::string GetRandomNumber();

std::string GetTime();

std::string GetDateTime();

std::string GenerateSocketPath();

const char* GetEventName(EVENT_ID event);

std::string NormalizeJson(const std::string& jsonStr);
