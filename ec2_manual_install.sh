rpm -e aikido-php-firewall
systemctl restart php-fpm.service
wget https://aikido-firewall.s3.eu-west-1.amazonaws.com/v$1/linux_x86_64/aikido-php-firewall-$1-1.x86_64.rpm
rpm -ivh aikido-php-firewall-$1-1.x86_64.rpm

ls -l /opt/aikido
ls -l /etc/php.d/ | grep aikido
ls -l /lib64/php8.2/modules/ | grep aikido
ls -l /run/aikido.sock

systemctl restart php-fpm.service

tail -f /var/log/aikido/aikido_go_*