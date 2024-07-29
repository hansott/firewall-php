VERSION=$(grep '# define PHP_AIKIDO_VERSION' include/php_aikido.h | awk -F'"' '{print $2}')

sudo rpm -Uvh --replacepkgs ~/rpmbuild/RPMS/aarch64/aikido-php-firewall-$VERSION-1.aarch64.rpm
