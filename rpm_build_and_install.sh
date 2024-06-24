rm -rf ~/rpmbuild
rpmdev-setuptree

mkdir -p ~/rpmbuild/SOURCES/aikido-1.0.0
cp -rf package/rpm/etc ~/rpmbuild/SOURCES/aikido-1.0.0/
cp -rf package/rpm/opt ~/rpmbuild/SOURCES/aikido-1.0.0/

cp -f package/rpm/aikido.spec ~/rpmbuild/SPECS/

cp build/aikido ~/rpmbuild/SOURCES/aikido-1.0.0/opt/aikido/aikido-1.0.0
cp build/modules/aikido.so ~/rpmbuild/SOURCES/aikido-1.0.0/opt/aikido/aikido-1.0.0-extension-php-8.1.so

cd ~/rpmbuild/SOURCES
tar czvf ~/rpmbuild/SOURCES/aikido-1.0.0.tar.gz *
rm -rf ~/rpmbuild/SOURCES/aikido-1.0.0
cd ~/Work/firewall-php

rpmbuild --define "debug_package %{nil}" -ba ~/rpmbuild/SPECS/aikido.spec