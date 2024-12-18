#pragma once

ZEND_FUNCTION(set_user);

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_aikido_set_user, 0, 1, _IS_BOOL, 0)
    ZEND_ARG_TYPE_INFO(0, id, IS_STRING, 0)
    ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
ZEND_END_ARG_INFO()

bool send_user_event(std::string id, std::string username);
