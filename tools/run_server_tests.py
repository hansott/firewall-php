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


def handle_test_scenario(test_scenario):
    try:
        # Generate unique ports for mock server and PHP server.
        mock_port = generate_unique_port()
        php_port = generate_unique_port()

        print(f"Starting mock server on port {mock_port} for {test_scenario}...")

        mock_server_process = subprocess.Popen(['python', 'mock_server.py', str(mock_port)])
        time.sleep(2)

        print(f"Running mock_setup.py with port {mock_port} from {test_scenario}...")
        mock_setup_process = subprocess.run(['python', os.path.join(test_scenario, 'mock_setup.py'), str(mock_port)],
                                            check=True)

        print(f"Starting PHP server on port {php_port} for {test_scenario}...")
        env = os.environ.copy()
        env.update({
            'AIKIDO_TOKEN': 'AIK_RUNTIME_MOCK',
            'AIKIDO_ENDPOINT': f'http://localhost:{mock_port}/',
            'AIKIDO_CONFIG_ENDPOINT': f'http://localhost:{mock_port}/'
        })

        php_server_process = subprocess.Popen(
            ['php', '-S', f'localhost:{php_port}', '-t', subfolder],
            env=env
        )
        time.sleep(2)

        print(f"Running test.py in {test_scenario}...")
        subprocess.run(['python', os.path.join(test_scenario, 'test.py')], check=True, timeout=180)

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

        if mock_server_process:
            mock_server_process.terminate()
            mock_server_process.wait()
            print(f"Mock server on port {mock_port} stopped.")


def main(root_tests_dir):
    test_dirs = [f.path for f in os.scandir(root_tests_dir) if f.is_dir()]
    threads = []

    for test_dir in test_dirs:
        thread = threading.Thread(target=handle_test_scenario, args=(test_dir,))
        threads.append(thread)
        thread.start()

    # Wait for all threads to complete.
    for thread in threads:
        thread.join()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python script.py <root_folder_path>")
        exit(1)

    root_folder = sys.argv[1]
    main(root_folder)
