#include "HandleShouldBlockRequest.h"
#include "Actions.h"

zend_class_entry *aikido_request_block_ce = nullptr;

ZEND_FUNCTION(should_block_request)
{
    if (AIKIDO_GLOBAL(disable) == true) {
        return;
    }

    if (!aikido_request_block_ce) {
        return;
    }

    object_init_ex(return_value, aikido_request_block_ce);

    zend_object *obj = Z_OBJ_P(return_value);
    if (!obj) {
        return;
    }    
    zend_update_property_bool(aikido_request_block_ce, obj, "block", sizeof("block") - 1, action.Block());
    zend_update_property_string(aikido_request_block_ce, obj, "type", sizeof("type") - 1, action.Type());
    zend_update_property_string(aikido_request_block_ce, obj, "trigger", sizeof("trigger") - 1, action.Trigger());
    zend_update_property_string(aikido_request_block_ce, obj, "ip", sizeof("ip") - 1, action.Ip());
}

void RegisterRequestBlockObject() {
    zend_class_entry ce;
    INIT_CLASS_ENTRY(ce, "AikidoRequestBlock", NULL); // Register class without methods
    aikido_request_block_ce = zend_register_internal_class(&ce);

    zend_declare_property_bool(aikido_request_block_ce, "block", sizeof("block") - 1, 0, ZEND_ACC_PUBLIC);
    zend_declare_property_string(aikido_request_block_ce, "type", sizeof("type") - 1, "", ZEND_ACC_PUBLIC);
    zend_declare_property_string(aikido_request_block_ce, "trigger", sizeof("trigger") - 1, "", ZEND_ACC_PUBLIC);
    zend_declare_property_string(aikido_request_block_ce, "ip", sizeof("ip") - 1, "", ZEND_ACC_PUBLIC);
}

