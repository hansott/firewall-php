#include "Includes.h"

void inject_ip_blocking_check(zend_op_array* op_array) {
    zend_op* opline = zend_arena_alloc(&CG(arena), sizeof(zend_op));
    memset(opline, 0, sizeof(zend_op));
    opline->opcode = ZEND_DO_FCALL;
    opline->op1_type = IS_CONST;
    opline->op1.constant = zend_add_literal(op_array, zend_string_init("aikido_ip_blocking_check", sizeof("aikido_ip_blocking_check") - 1, 0));
    opline->result_type = IS_UNUSED;
    opline->result.var = 0;
    zend_llist_add_element(&op_array->opcodes, opline);
}

zend_op_array* handle_file_compilation(zend_file_handle* file_handle, int type) {
    eventCache.Reset();
    switch (type) {
        case ZEND_INCLUDE:
            eventCache.functionName = "include";
            break;
        case ZEND_INCLUDE_ONCE:
            eventCache.functionName = "include_once";
            break;
        case ZEND_REQUIRE:
            eventCache.functionName = "require";
            break;
        case ZEND_REQUIRE_ONCE:
            eventCache.functionName = "require_once";
            break;
        default:
            return original_file_compilation_handler(file_handle, type);
    }

    ScopedTimer scopedTimer(eventCache.functionName);

    char* filename = PHP_GET_CHAR_PTR(file_handle->filename);
    
    AIKIDO_LOG_DEBUG("\"%s\" called for \"%s\"!\n", eventCache.functionName.c_str(), filename);

    EVENT_ID eventId = NO_EVENT_ID;
    helper_handle_pre_file_path_access(filename, eventId);

    if (aikido_process_event(eventId, eventCache.functionName) == BLOCK) {
        // exit zend_compile_file handler and do not call the original handler, thus blocking the script file compilation
        return nullptr;
    }

    scopedTimer.Stop();
    zend_op_array* op_array = original_file_compilation_handler(file_handle, type);
    scopedTimer.Start();

    inject_ip_blocking_check(op_array);

    eventId = NO_EVENT_ID;
    helper_handle_post_file_path_access(eventId);
    aikido_process_event(eventId, eventCache.functionName);

    return op_array;
}
