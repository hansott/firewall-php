#include "HandlePathTraversal.h"
#include "Utils.h"

/* Handles PHP functions that have a file path as first parameter */
AIKIDO_HANDLER_FUNCTION(handle_file_path_access) {
    zend_string *filename = NULL;
    zend_string *filename2 = NULL;

    ZEND_PARSE_PARAMETERS_START(0, -1)
        Z_PARAM_OPTIONAL
        Z_PARAM_STR(filename)
        Z_PARAM_STR(filename2)
    ZEND_PARSE_PARAMETERS_END();

    if (!filename) {
        return;
    }

    std::string filenameString(ZSTR_VAL(filename));

    std::string functionNameString(AIKIDO_GET_FUNCTION_NAME());
    
    inputEvent = {
        { "event", "before_function_executed" },
        { "data", {
            { "function_name", "path_accessed" },
            { "parameters", {
                { "filename", filenameString },
                { "operation", functionNameString }
            } }
        } }
    };

    if (filename2) {
        std::string filenameString2(ZSTR_VAL(filename2));
        inputEvent["data"]["parameters"]["filename2"] = filenameString2;
    }
}

