Name:           aikido
Version:        1.0.0
Release:        1%{?dist}
Summary:        Aikido PHP extension and agent

License:        GPL
URL:            https://aikido.dev
Source0:        aikido-1.0.0.tar.gz

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

VERSION="%{version}"

declare -A php_api_versions

php_api_versions=(
    ["8.1"]="20210902"
)

# Find all installed PHP versions
PHP_VERSIONS=$(ls /etc/php)
for PHP_VERSION in $PHP_VERSIONS; do
    PHP_API_VERSION=${php_api_versions[$PHP_VERSION]}
    EXT_DIR="/usr/lib/php/$PHP_API_VERSION"
    if [ -d "$EXT_DIR" ]; then
        # Install PHP extension
        echo "Installing Aikido extension for PHP $PHP_VERSION (API VERSION $PHP_API_VERSION)..."
        if [ -f "/opt/aikido/aikido-$VERSION-extension-php-$PHP_VERSION.so" ]; then
            echo "Installing Aikido extension in $EXT_DIR/aikido.so"
            ln -sf /opt/aikido/aikido-$VERSION-extension-php-$PHP_VERSION.so $EXT_DIR/aikido.so
        else
            echo "Extension file /opt/aikido/aikido-$VERSION-extension-php-$PHP_VERSION.so does not exist. Skipping."
        fi

        # Installing Aikido mod
        MOD_PATHS=(
            "/etc/php/$PHP_VERSION/mods-available"
            "/etc/php/$PHP_VERSION/cli/conf.d"
            "/etc/php/$PHP_VERSION/apache2/conf.d"
        )

        for MOD_PATH in "${MOD_PATHS[@]}"; do
            if [ -d "$MOD_PATH" ]; then
                if [[ "$MOD_PATH" == *"conf.d"* ]]; then
                    echo "Installing Aikido mod in $MOD_PATH/20-aikido.ini..."
                    ln -sf /opt/aikido/aikido-dev.ini $MOD_PATH/20-aikido.ini
                else
                    echo "Installing Aikido mod in $MOD_PATH/aikido.ini..."
                    ln -sf /opt/aikido/aikido-dev.ini $MOD_PATH/aikido.ini
                fi
            else
                echo "Mod path $MOD_PATH does not exist for PHP $PHP_VERSION. Skipping."
            fi
        done
    else
        echo "PHP extension directory $EXT_DIR does not exist for PHP $PHP_VERSION. Skipping."
    fi
done

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


%preun
#!/bin/bash

echo "Starting the uninstallation process for Aikido..."

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

declare -A php_api_versions

php_api_versions=(
    ["8.1"]="20210902"
)

# Remove the PHP extension
if [ -d /etc/php ]; then
    echo "Removing the Aikido PHP extension..."
    PHP_VERSIONS=$(ls /etc/php)
    for PHP_VERSION in $PHP_VERSIONS; do
        PHP_API_VERSION=${php_api_versions[$PHP_VERSION]}
        EXT_DIR="/usr/lib/php/$PHP_API_VERSION"
        if [ -d "$EXT_DIR" ]; then
            echo "Uninstalling Aikido extension for PHP $PHP_VERSION (API VERSION $PHP_API_VERSION)..."
            if [ -L "$EXT_DIR/aikido.so" ]; then
                echo "Uninstalling Aikido extension from $EXT_DIR/aikido.so..."
                unlink $EXT_DIR/aikido.so
            else
                echo "Aikido extension not found for PHP $PHP_VERSION."
            fi
        else
            echo "PHP extension directory $EXT_DIR does not exist for PHP $PHP_VERSION. Skipping."
            continue
        fi

        # Uninstalling Aikido mod
        MOD_PATHS=(
            "/etc/php/$PHP_VERSION/mods-available"
            "/etc/php/$PHP_VERSION/cli/conf.d"
            "/etc/php/$PHP_VERSION/apache2/conf.d"
        )

        for MOD_PATH in "${MOD_PATHS[@]}"; do
            if [ -d "$MOD_PATH" ]; then
                if [[ "$MOD_PATH" == *"conf.d"* ]]; then
                    echo "Uninstalling Aikido mod from $MOD_PATH/20-aikido.ini..."
                    rm -f $MOD_PATH/20-aikido.ini
                else
                    echo "Uninstalling Aikido mod from $MOD_PATH/aikido.ini..."
                    rm -f $MOD_PATH/aikido.ini
                fi
            else
                echo "Mod path $MOD_PATH does not exist for PHP $VERSION. Skipping."
            fi
        done
    done
else
    echo "/etc/php directory not found. Skipping PHP extension removal."
fi

# Remove the Aikido log file
if [ -f "/var/log/aikido.log" ]; then
    echo "Removing /var/log/aikido.log..."
    rm -f /var/log/aikido.log
else
    echo "/var/log/aikido.log does not exist. Skipping."
fi

echo "Uninstallation process for Aikido completed."

%files
/etc/systemd/system/aikido.service
/opt/aikido

%changelog
* Fri Jun 21 2024 Aikido <hello@aikido.dev> - 1.0.0-1
- Initial package
