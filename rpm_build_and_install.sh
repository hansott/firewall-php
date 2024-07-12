rm -rf ~/rpmbuild
rpmdev-setuptree

VERSION="2.0.1"

mkdir -p ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION
cp -rf package/rpm/opt ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/

cp -f package/rpm/aikido.spec ~/rpmbuild/SPECS/

cp build/aikido_agent.so ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido/aikido_agent_$VERSION.so
cp build/aikido_request_processor.so ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido/aikido_agent_$VERSION.so
cp build/modules/aikido.so ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION/opt/aikido/aikido-$VERSION-extension-php-8.1.so

cd ~/rpmbuild/SOURCES
tar czvf ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION.tar.gz *
rm -rf ~/rpmbuild/SOURCES/aikido-php-firewall-$VERSION

rpmbuild --define "debug_package %{nil}" -ba ~/rpmbuild/SPECS/aikido.spec

rpm -ivh ~/rpmbuild/RPMS/x86_64/aikido-php-firewall-$VERSION-1.x86_64.rpm
