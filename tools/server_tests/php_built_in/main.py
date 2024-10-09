import os
import subprocess

def handle_php_built_in(test_dir, test_lib_dir, loaded_env, server_port, mock_port, valgrind, debug):
    env = os.environ.copy()
    env.update({
        'AIKIDO_LOG_LEVEL': 'DEBUG' if debug else 'ERROR',
        'AIKIDO_TOKEN': 'AIK_RUNTIME_MOCK',
        'AIKIDO_ENDPOINT': f'http://localhost:{mock_port}/',
        'AIKIDO_REALTIME_ENDPOINT': f'http://localhost:{mock_port}/',
    })
    env.update(load_env_from_json(env_file_path))

    php_server_process_cmd = ['php', '-S', f'localhost:{server_port}', '-t', test_dir]
    if valgrind:
        php_server_process_cmd = ['valgrind', f'--suppressions={test_lib_dir}/valgrind.supp', '--track-origins=yes'] + php_server_process_cmd
        
    return [subprocess.Popen(
        php_server_process_cmd,
        env=env
    )]
