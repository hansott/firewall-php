#include "Includes.h"

zend_op_array* handle_file_compilation(zend_file_handle *file_handle, int type) {
    return original_file_compilation(file_handle, type);   
}
