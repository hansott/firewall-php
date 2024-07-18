Name:           aikido-php-firewall
Version:        1.0.1
Release:        1
Summary:        Aikido PHP extension and Agent

License:        GPL
URL:            https://aikido.dev
Source0:        aikido-php-firewall-%{version}.tar.gz

%description
Aikido PHP extension and agent.

%prep
%setup -q

%build
# Build steps if required

%install
mkdir -p %{buildroot}/opt/aikido
cp -rf opt/aikido/* %{buildroot}/opt/aikido

%post
#!/bin/bash

echo "Starting the installation process for Aikido..."

sudo mkdir /var/log/aikido
sudo chmod 777 /var/log/aikido

VERSION="%{version}"
PHP_VERSION=$(php -v | head -n 1 | awk '{print $2}' | cut -d '.' -f1,2)
echo "Found PHP version $PHP_VERSION!"

# Install PHP extension
PHP_EXT_DIR=$(php -i | grep "^extension_dir" | awk '{print $3}')
echo "Installing Aikido extension in $PHP_EXT_DIR..."

if [ -d "$PHP_EXT_DIR" ]; then
    echo "Installing Aikido extension in $PHP_EXT_DIR/aikido.so..."
    ln -sf /opt/aikido/aikido-$VERSION-extension-php-$PHP_VERSION.so $PHP_EXT_DIR/aikido.so
else
    echo "No extension dir. Exiting."
    exit 1
fi

# Installing Aikido mod
PHP_MOD_DIR=$(php -i | grep "Scan this dir for additional .ini files" | awk -F"=> " '{print $2}')
echo "Installing Aikido mod in $PHP_MOD_DIR..."

if [ -d "$PHP_MOD_DIR" ]; then
    echo "Installing Aikido mod in $PHP_MOD_DIR/zz-aikido-firewall.ini..."
    ln -sf /opt/aikido/aikido-dev.ini $PHP_MOD_DIR/zz-aikido-firewall.ini
else
    echo "No mod dir. Exiting."
    exit 1
fi


%preun
#!/bin/bash

echo "Starting the uninstallation process for Aikido..."

VERSION="%{version}"
PHP_VERSION=$(php -v | head -n 1 | awk '{print $2}' | cut -d '.' -f1,2)
PHP_EXT_DIR=$(php -i | grep "^extension_dir" | awk '{print $3}')
PHP_MOD_DIR=$(php -i | grep "Scan this dir for additional .ini files" | awk -F"=> " '{print $2}')

echo "Found PHP version $PHP_VERSION!"

# Function to handle script termination
cleanup() {
    echo "Uninstallation script was terminated unexpectedly."
    exit 1
}

# Trap termination signals
trap cleanup SIGTERM SIGINT

# Uninstall PHP extension
echo "Uinstalling Aikido extension from EXT_DIR $PHP_EXT_DIR..."

if [ -d "$PHP_EXT_DIR" ]; then
    echo "Uinstalling Aikido extension from $EXT_DIR/aikido.so..."
    rm -f $PHP_EXT_DIR/aikido.so
else
    echo "No extension dir. Exiting."
    exit 1
fi

# Uninstalling Aikido mod
echo "Unistalling Aikido mod from $PHP_MOD_DIR..."

if [ -d "$PHP_MOD_DIR" ]; then
    echo "Uninstalling Aikido mod from $PHP_MOD_DIR/zz-aikido-firewall.ini..."
    rm -f $PHP_MOD_DIR/zz-aikido-firewall.ini
else
    echo "No mod dir. Exiting."
    exit 1
fi

# Remove the Aikido socket
SOCKET_PATH="/run/aikido.sock"

if [ -S "$SOCKET_PATH" ]; then
    echo "Removing $SOCKET_PATH ..."
    rm "$SOCKET_PATH"
    if [ $? -eq 0 ]; then
        echo "Socket removed successfully."
    else
        echo "Failed to remove the socket."
    fi
else
    echo "Socket $SOCKET_PATH does not exist."
fi

sudo rm -rf /var/log/aikido

echo "Uninstallation process for Aikido completed."

%files
/opt/aikido

%changelog
* Fri Jun 21 2024 Aikido <hello@aikido.dev> - %{version}-1
- New package version
