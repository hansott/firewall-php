#include "HandleShouldBlockRequest.h"
#include "Actions.h"

zend_class_entry *aikido_block_request_status_class = nullptr;

ZEND_FUNCTION(should_block_request)
{
    if (AIKIDO_GLOBAL(disable) == true) {
        return;
    }

    if (!aikido_block_request_status_class) {
        return;
    }

    object_init_ex(return_value, aikido_block_request_status_class);
    #if PHP_VERSION_ID >= 80000
        zend_object *obj = Z_OBJ_P(return_value);
        if (!obj) {
            return;
        }
    #else
        zval *obj = return_value;
    #endif
    zend_update_property_bool(aikido_block_request_status_class, obj, "block", sizeof("block") - 1, action.Block());
    zend_update_property_string(aikido_block_request_status_class, obj, "type", sizeof("type") - 1, action.Type());
    zend_update_property_string(aikido_block_request_status_class, obj, "trigger", sizeof("trigger") - 1, action.Trigger());
    zend_update_property_string(aikido_block_request_status_class, obj, "ip", sizeof("ip") - 1, action.Ip());
}

void RegisterAikidoBlockRequestStatusClass() {
    zend_class_entry ce;
    INIT_CLASS_ENTRY(ce, "AikidoBlockRequestStatus", NULL); // Register class without methods
    aikido_block_request_status_class = zend_register_internal_class(&ce);

    zend_declare_property_bool(aikido_block_request_status_class, "block", sizeof("block") - 1, 0, ZEND_ACC_PUBLIC);
    zend_declare_property_string(aikido_block_request_status_class, "type", sizeof("type") - 1, "", ZEND_ACC_PUBLIC);
    zend_declare_property_string(aikido_block_request_status_class, "trigger", sizeof("trigger") - 1, "", ZEND_ACC_PUBLIC);
    zend_declare_property_string(aikido_block_request_status_class, "ip", sizeof("ip") - 1, "", ZEND_ACC_PUBLIC);
}

