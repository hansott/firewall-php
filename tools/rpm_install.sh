VERSION=$(grep '#define PHP_AIKIDO_VERSION' lib/php-extension/include/php_aikido.h | awk -F'"' '{print $2}')

sudo rpm -Uvh --oldpackage ~/rpmbuild/RPMS/aarch64/aikido-php-firewall-$VERSION-1.aarch64.rpm
