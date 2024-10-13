import os
import subprocess
import re

nginx_config_dir = "/etc/nginx/conf.d"
php_fpm_pool_dir = "/etc/php-fpm.d"
socket_folder = "/run/php-fpm"

nginx_conf_template = """
server {{
    listen {port};
    server_name {name};

    root {test_dir};
    index index.php;

    location / {{
        try_files $uri $uri/ /index.php?$args;
    }}

    location ~ \.php$ {{
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        fastcgi_pass unix:/run/php-fpm/php-fpm-{name}.sock;
        fastcgi_index index.php;
        include fastcgi.conf;
    }}
}}
"""

php_fpm_conf_template = """[{name}]
user = {user}
group = {user}
listen = /run/php-fpm/php-fpm-{name}.sock
listen.owner = nginx
listen.group = nginx
pm = dynamic
pm.max_children = 5
pm.start_servers = 2
pm.min_spare_servers = 1
pm.max_spare_servers = 3
clear_env = no

"""

def enable_config_line(file_path, line_to_check, comment_ch):
    # Read the nginx.conf file
    with open(file_path, 'r') as file:
        lines = file.readlines()

    # Prepare a regex pattern to match the commented line
    commented_line_pattern = r"\s*" + comment_ch + r"\s*" + re.escape(line_to_check.strip()) + r"\s*"

    # Initialize a flag to track changes
    changes_made = False

    # Iterate through the lines and check for the commented line
    with open(file_path, 'w') as file:
        for line in lines:
            if re.match(commented_line_pattern, line):
                # If the line is commented, uncomment it
                file.write(line.replace(comment_ch, "", 1).lstrip())
                changes_made = True
            else:
                # Otherwise, write the line as-is
                file.write(line)

    if changes_made:
        print(f"The line '{line_to_check}' was uncommented.")
    else:
        print(f"The line '{line_to_check}' was already uncommented or not found.")


def nginx_create_conf_file(test_name, test_dir, server_port):
    nginx_config = nginx_conf_template.format(
        name = test_name,
        port = server_port,
        test_dir = test_dir
    )

    nginx_config_file = os.path.join(nginx_config_dir, f"{test_name}.conf")
    with open(nginx_config_file, "w") as fpm_file:
        fpm_file.write(nginx_config)

    print(f"Configured nginx config for {test_name}")


def php_fpm_create_conf_file(test_name, user, env):
    php_fpm_config = php_fpm_conf_template.format(
        name = test_name,
        user = "ttimcu",
    )
    
    for e in env:
        php_fpm_config += f"env[{e}] = {env[e]}\n"

    php_fpm_config_file = os.path.join(php_fpm_pool_dir, f"{test_name}.conf")
    with open(php_fpm_config_file, "w") as fpm_file:
        fpm_file.write(php_fpm_config)

    print(f"Configured PHP-FPM config for {test_name}")

def handle_nginx_php_fpm(test_name, test_dir, test_lib_dir, env, server_port, mock_port, valgrind, debug):
    enable_config_line("/etc/nginx/nginx.conf", "include /etc/nginx/conf.d/*.conf;", '#')
    nginx_create_conf_file(test_name, test_dir, server_port)

    enable_config_line("/etc/php-fpm.conf", "include=/etc/php-fpm.d/*.conf", ';')
    php_fpm_create_conf_file(test_name, "ttimcu", env)
    
    subprocess.run(['systemctl', 'restart', 'nginx.service'], check=True)
    subprocess.run(['systemctl', 'restart', 'php-fpm.service'], check=True)
    return []
    