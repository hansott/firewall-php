Name:           aikido-php-firewall
Version:        1.20.0
Release:        1
Summary:        Aikido PHP extension and agent

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
mkdir -p %{buildroot}/etc/systemd/system/
cp -f etc/systemd/system/aikido.service %{buildroot}/etc/systemd/system/aikido.service

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
echo "Installing Aikido extension in EXT_DIR $PHP_EXT_DIR..."

if [ -d "$PHP_EXT_DIR" ]; then
    echo "Installing Aikido extension in $EXT_DIR/aikido.so..."
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

echo "Installing Aikido agent $VERSION..."
if [ -f "/opt/aikido/aikido-$VERSION" ]; then
    ln -sf /opt/aikido/aikido-$VERSION /opt/aikido/aikido
else
    echo "Aikido agent directory /opt/aikido/aikido-$VERSION does not exist. Exiting."
    exit 1
fi

echo "Registering Aikido agent to run as service..."
sudo systemctl daemon-reload
sudo systemctl enable aikido.service
sudo systemctl start aikido.service

sleep 10

echo "Installing SE Linux module for allowing access to /run/aikido.sock..."
sudo semodule -i /opt/aikido/aikido.pp
sudo chcon -t var_run_t /run/aikido.sock

%preun
#!/bin/bash

echo "Starting the uninstallation process for Aikido..."

VERSION="%{version}"
PHP_VERSION=$(php -v | head -n 1 | awk '{print $2}' | cut -d '.' -f1,2)

echo "Found PHP version $PHP_VERSION!"

# Function to handle script termination
cleanup() {
    echo "Uninstallation script was terminated unexpectedly."
    exit 1
}

# Trap termination signals
trap cleanup SIGTERM SIGINT

# Stop the service if it is running
echo "Stopping the Aikido service..."
sudo systemctl stop aikido.service

sleep 10

echo "Disabling the Aikido service..."
sudo systemctl disable aikido.service
sudo systemctl daemon-reload

# Unlink the main executable if it exists
if [ -L /opt/aikido/aikido ]; then
    echo "Unlinking the Aikido main executable..."
    unlink /opt/aikido/aikido
else
    echo "Aikido main executable not found. Skipping unlink step."
fi


# Uninstall PHP extension
PHP_EXT_DIR=$(php -i | grep "^extension_dir" | awk '{print $3}')
echo "Uinstalling Aikido extension from EXT_DIR $PHP_EXT_DIR..."

if [ -d "$PHP_EXT_DIR" ]; then
    echo "Uinstalling Aikido extension from $EXT_DIR/aikido.so..."
    rm -f $PHP_EXT_DIR/aikido.so
else
    echo "No extension dir. Exiting."
    exit 1
fi

# Uninstalling Aikido mod
PHP_MOD_DIR=$(php -i | grep "Scan this dir for additional .ini files" | awk -F"=> " '{print $2}')
echo "Unistalling Aikido mod from $PHP_MOD_DIR..."

if [ -d "$PHP_MOD_DIR" ]; then
    echo "Uninstalling Aikido mod from $PHP_MOD_DIR/zz-aikido-firewall.ini..."
    rm -f $PHP_MOD_DIR/zz-aikido-firewall.ini
else
    echo "No mod dir. Exiting."
    exit 1
fi

# Remove the Aikido log file
if [ -f "/var/log/aikido.log" ]; then
    echo "Removing /var/log/aikido.log..."
    rm -f /var/log/aikido.log
else
    echo "/var/log/aikido.log does not exist. Skipping."
fi

# Remove semodule
sudo semodule -r aikido

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
/etc/systemd/system/aikido.service
/opt/aikido

%changelog
* Fri Jun 21 2024 Aikido <hello@aikido.dev> - %{version}-1
- New package version
