rm -rf ~/rpmbuild
rpmdev-setuptree

VERSION=$(grep '# define PHP_AIKIDO_VERSION' lib/php-extension/include/php_aikido.h | awk -F'"' '{print $2}')

mkdir -p ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION

cp -rf package/rpm/opt ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/

cp -f package/rpm/aikido.spec ~/rpmbuild/SPECS/

cp build/aikido-agent.so ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido/aikido-agent.so
cp build/aikido-request-processor.so ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido/aikido-request-processor.so
cp build/modules/aikido-extension-php-8.0.so ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido/aikido-extension-php-8.0.so

mv ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido-$VERSION

sed -i "s/aikido.so/aikido-$VERSION.so/" ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido-$VERSION/aikido.ini

cd ~/rpmbuild/SOURCES
tar czvf ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION.tar.gz *
rm -rf ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION

rpmbuild --define "debug_package %{nil}" -ba ~/rpmbuild/SPECS/aikido.spec
