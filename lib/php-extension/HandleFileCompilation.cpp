#include "Includes.h"

zend_op_array* handle_file_compilation(zend_file_handle* file_handle, int type) {
    AIKIDO_LOG_INFO("\"zend_compile_file\" called for \"%s\"!\n", file_handle->filename);

    eventCache.Reset();
    eventCache.functionName = "zend_compile_file";

    EVENT_ID eventId = NO_EVENT_ID;
    helper_handle_pre_file_path_access(file_handle->filename, eventId);

    if (aikido_process_event(eventId) == BLOCK) {
        // exit zend_compile_file handler and do not call the original handler, thus blocking the script file compilation
        return nullptr;
    }

    zend_op_array* op_array = original_file_compilation_handler(file_handle, type);

    helper_handle_post_file_path_access(eventId);
    aikido_process_event(eventId);

    return op_array;
}
