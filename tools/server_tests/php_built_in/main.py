import os
import subprocess

def handle_php_built_in(test_name, test_dir, test_lib_dir, loaded_env, server_port, mock_port, valgrind, debug):
    php_server_process_cmd = ['php', '-S', f'localhost:{server_port}', '-t', test_dir]
    if valgrind:
        php_server_process_cmd = ['valgrind', f'--suppressions={test_lib_dir}/valgrind.supp', '--track-origins=yes'] + php_server_process_cmd
        
    return [subprocess.Popen(
        php_server_process_cmd,
        env=loaded_env
    )]
