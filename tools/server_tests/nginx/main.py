def handle_nginx_php_fpm(test_dir, test_lib_dir, loaded_env, server_port, mock_port, valgrind, debug):
    nginx_config_dir = "/etc/nginx/sites-available"
    nginx_sites_enabled_dir = "/etc/nginx/sites-enabled"
    php_fpm_pool_dir = "/etc/php/7.4/fpm/pool.d"

    nginx_server_template = """
server {{
    listen 80;
    server_name {server_name};

    root {document_root};
    index index.php;

    location / {{
        try_files $uri $uri/ /index.php?$args;
    }}

    location ~ \.php$ {{
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/run/php/{php_fpm_socket};
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }}

    location ~ /\.ht {{
        deny all;
    }}
}}
"""
