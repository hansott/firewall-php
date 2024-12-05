#include "Includes.h"

bool SendUserEvent(std::string id, std::string username) {
    requestCache.userId = id;
    requestCache.userName = username;

    try {
        std::string output;
        requestProcessor.SendEvent(EVENT_SET_USER, output);
        action.Execute(output);
        return true;
    } catch (const std::exception &e) {
        AIKIDO_LOG_ERROR("Exception encountered in processing user event: %s\n", e.what());
    }
    return false;
}

// Exports the "\aikido\set_user" function, to be called from PHP user code.
// Receives two parameters: id and name (both strings).
// Returns true if the setting of the user succeeded, false otherwise.
ZEND_FUNCTION(set_user) {
    if (AIKIDO_GLOBAL(disable) == true) {
        RETURN_BOOL(false);
    }

    requestProcessor.ReloadConfig();

    char *id;
    size_t id_len;
    char *name;
    size_t name_len;

    ZEND_PARSE_PARAMETERS_START(2, 2)
    Z_PARAM_STRING(id, id_len)
    Z_PARAM_STRING(name, name_len)
    ZEND_PARSE_PARAMETERS_END();

    RETURN_BOOL(SendUserEvent(std::string(id, id_len), std::string(name, name_len)));
}
