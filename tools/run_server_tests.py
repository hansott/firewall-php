import os
import threading
import subprocess
import random
import time
import sys
import json

used_ports = set()

def generate_unique_port():
    while True:
        port = random.randint(1024, 65535)
        if port not in used_ports:
            used_ports.add(port)
            return port

def load_env_from_json(file_path):
    if not os.path.exists(file_path):
        return {}

    with open(file_path) as f:
        env_vars = json.load(f)
        return env_vars

def handle_test_scenario(test_dir, test_lib_dir):
    try:
        # Generate unique ports for mock server and PHP server.
        mock_port = generate_unique_port()
        php_port = generate_unique_port()

        test_name = os.path.basename(os.path.normpath(test_dir))

        config_path = os.path.join(test_dir, 'start_config.json')
        env_file_path = os.path.join(test_dir, 'env.json')

        print(f"Running {test_name}...")
        print(f"Starting mock server on port {mock_port} with start_config.json for {test_name}...")
        mock_aikido_core = subprocess.Popen(['python', 'mock_aikido_core.py', str(mock_port), config_path])
        time.sleep(2)

        print(f"Starting PHP server on port {php_port} for {test_name}...")
        env = os.environ.copy()
        env.update(load_env_from_json(env_file_path))
        env.update({
            'AIKIDO_LOG_LEVEL': 'DEBUG',
            'AIKIDO_TOKEN': 'AIK_RUNTIME_MOCK',
            'AIKIDO_ENDPOINT': f'http://localhost:{mock_port}/',
            'AIKIDO_CONFIG_ENDPOINT': f'http://localhost:{mock_port}/',
        })
        php_server_process = subprocess.Popen(
            ['php', '-S', f'localhost:{php_port}', '-t', test_dir],
            env=env
        )
        time.sleep(2)

        print(f"Running test.py for {test_name}...")
        subprocess.run(['python', 'test.py', str(php_port), str(mock_port)], 
                       env=dict(os.environ, PYTHONPATH=f"{test_lib_dir}:$PYTHONPATH"),
                       cwd=test_dir,
                       check=True, timeout=600)

    except subprocess.CalledProcessError as e:
        print(f"Error in testing scenario {test_name}:")
        print(f"Test output: {e.output}")

        sys.exit(1)
        
    except subprocess.TimeoutExpired:
        print(f"Error in testing scenario {test_name}:")
        print(f"Execution timed out.")
        sys.exit(1)
        
    finally:
        if php_server_process:
            php_server_process.terminate()
            php_server_process.wait()
            print(f"PHP server on port {php_port} stopped.")

        if mock_aikido_core:
            mock_aikido_core.terminate()
            mock_aikido_core.wait()
            print(f"Mock server on port {mock_port} stopped.")


def main(root_tests_dir, test_lib_dir, specific_test=None):
    if specific_test:
        specific_test = os.path.join(root_tests_dir, specific_test)
        handle_test_scenario(specific_test, test_lib_dir)
    else:
        test_dirs = [f.path for f in os.scandir(root_tests_dir) if f.is_dir()]
        for test_dir in test_dirs:
            handle_test_scenario(test_dir, test_lib_dir)

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python script.py <root_folder_path> <test_lib_dir> [specific_test]")
        exit(1)

    root_folder = os.path.abspath(sys.argv[1])
    test_lib_dir = os.path.abspath(sys.argv[2])
    specific_test = sys.argv[3] if len(sys.argv) > 3 else None
    main(root_folder, test_lib_dir, specific_test)
