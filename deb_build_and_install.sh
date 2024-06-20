sudo apt purge aikido-firewall-php
dpkg --build package/deb aikido-firewall-php.deb
sudo dpkg -i aikido-firewall-php.deb
