#include "HandleUsers.h"
#include "Cache.h"

bool send_user_event(std::string id, std::string username) {
    requestCache.userId = id;
    requestCache.userName = username;

    try
    {
        std::string output;
        GoRequestProcessorOnEvent(EVENT_PRE_USER, output);
        aikido_execute_output(output);
        return true;
    }
    catch (const std::exception &e)
    {
        AIKIDO_LOG_ERROR("Exception encountered in processing user event: %s\n", e.what());
    }
    return false;
}