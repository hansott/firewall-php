/* aikido extension for PHP */

#ifndef PHP_AIKIDO_H
# define PHP_AIKIDO_H

extern zend_module_entry aikido_module_entry;
# define phpext_aikido_ptr &aikido_module_entry

# define PHP_AIKIDO_VERSION "0.1.0"

# if defined(ZTS) && defined(COMPILE_DL_AIKIDO)
ZEND_TSRMLS_CACHE_EXTERN()
# endif

#endif	/* PHP_AIKIDO_H */
