VERSION=$(grep '# define PHP_AIKIDO_VERSION' include/php_aikido.h | awk -F'"' '{print $2}')

sudo rpm -Uvh https://aikido-firewall.s3.eu-west-1.amazonaws.com/v$VERSION/linux_x86_64/aikido-php-firewall-$VERSION-1.x86_64.rpm