#include "Includes.h"

ZEND_BEGIN_ARG_INFO(arginfo_aikido_should_block_request, 0)
    // No arguments
ZEND_END_ARG_INFO()

ZEND_FUNCTION(should_block_request);

void RegisterAikidoBlockRequestStatusClass();
