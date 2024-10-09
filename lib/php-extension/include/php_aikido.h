/* Aikido extension for PHP */

#ifndef PHP_AIKIDO_H
# define PHP_AIKIDO_H

extern zend_module_entry aikido_module_entry;
# define phpext_aikido_ptr &aikido_module_entry

# define PHP_AIKIDO_VERSION "1.0.79"

# if defined(ZTS) && defined(COMPILE_DL_AIKIDO)
ZEND_TSRMLS_CACHE_EXTERN()
# endif

ZEND_BEGIN_MODULE_GLOBALS(aikido)
    long log_level;
    bool blocking;
    bool disable;
    bool collect_api_schema;
    bool trust_proxy;
    bool localhost_allowed_by_default;
    std::string socket_path;
    std::string log_level_str;
    std::string sapi_name;
ZEND_END_MODULE_GLOBALS(aikido)

ZEND_EXTERN_MODULE_GLOBALS(aikido)

ZEND_FUNCTION(set_user);

#define AIKIDO_GLOBAL(v) ZEND_MODULE_GLOBALS_ACCESSOR(aikido, v)

#endif	/* PHP_AIKIDO_H */
