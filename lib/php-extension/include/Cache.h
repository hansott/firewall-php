#pragma once

#include "Includes.h"

class RequestCache {
public:
    std::string userId;
    std::string userName;

    RequestCache() = default;
    void Reset();
};

class EventCache {
public:
    std::string functionName;

    std::string filename;
    std::string filename2;

    std::string cmd;

    std::string outgoingRequestUrl;
    std::string outgoingRequestPort;
    
    std::string sql_query;
    std::string sql_dialect;

    EventCache() = default;
    void Reset();
};

extern RequestCache requestCache;
extern EventCache eventCache;
