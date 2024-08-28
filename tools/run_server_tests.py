import os
import threading
import subprocess
import random
import time
import sys

used_ports = set()

def generate_unique_port():
    while True:
        port = random.randint(1024, 65535)
        if port not in used_ports:
            used_ports.add(port)
            return port


def handle_test_scenario(test_scenario, test_lib_dir):
    try:
        # Generate unique ports for mock server and PHP server.
        mock_port = generate_unique_port()
        php_port = generate_unique_port()

        config_path = os.path.join(test_scenario, 'config.json')

        print(f"Starting mock server on port {mock_port} with config {config_path} for {test_scenario}...")
        mock_aikido_core = subprocess.Popen(['python', 'mock_aikido_core.py', str(mock_port), config_path])
        time.sleep(2)

        print(f"Starting PHP server on port {php_port} for {test_scenario}...")
        env = os.environ.copy()
        env.update({
            'AIKIDO_TOKEN': 'AIK_RUNTIME_MOCK',
            'AIKIDO_ENDPOINT': f'http://localhost:{mock_port}/',
            'AIKIDO_CONFIG_ENDPOINT': f'http://localhost:{mock_port}/'
        })
        php_server_process = subprocess.Popen(
            ['php', '-S', f'localhost:{php_port}', '-t', test_scenario],
            env=env
        )
        time.sleep(2)

        print(f"Running test.py in {test_scenario}...")
        subprocess.run(['python', 'test.py', str(php_port), str(mock_port)], 
                       env=dict(os.environ, PYTHONPATH=f"{test_lib_dir}:$PYTHONPATH"),
                       cwd=test_scenario,
                       check=True, timeout=180)

    except subprocess.CalledProcessError as e:
        print(f"Error in testing scenario {test_scenario}:")
        print(e.output)

        sys.exit(-1)
        
    except subprocess.TimeoutExpired:
        print(f"Error in testing scenario {test_scenario}:")
        print(f"Execution timed out.")
        sys.exit(-1)
        
    finally:
        if php_server_process:
            php_server_process.terminate()
            php_server_process.wait()
            print(f"PHP server on port {php_port} stopped.")

        if mock_aikido_core:
            mock_aikido_core.terminate()
            mock_aikido_core.wait()
            print(f"Mock server on port {mock_port} stopped.")


def main(root_tests_dir, test_lib_dir):
    test_dirs = [f.path for f in os.scandir(root_tests_dir) if f.is_dir()]
    threads = []

    for test_dir in test_dirs:
        thread = threading.Thread(target=handle_test_scenario, args=(test_dir,test_lib_dir))
        threads.append(thread)
        thread.start()

    # Wait for all threads to complete.
    for thread in threads:
        thread.join()

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python script.py <root_folder_path> <test_lib_dir>")
        exit(1)

    root_folder = sys.argv[1]
    test_lib_dir = sys.argv[2]
    main(root_folder, os.path.abspath(test_lib_dir))
