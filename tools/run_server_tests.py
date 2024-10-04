import os
import threading
import subprocess
import random
import time
import sys
import json
import argparse

used_ports = set()
passed_tests = []
failed_tests = []

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
    
def print_test_results(s, tests):
    if not len(tests):
        return
    
    print(s)
    for t in tests:
        print(f"\t- {t}")

def handle_test_scenario(root_tests_dir, test_dir, test_lib_dir, benchmark, valgrind, debug):
    try:
        # Generate unique ports for mock server and PHP server.
        mock_port = generate_unique_port()
        php_port = generate_unique_port()

        test_name = os.path.basename(os.path.normpath(test_dir))

        config_path = os.path.join(test_dir, 'start_config.json')
        env_file_path = os.path.join(test_dir, 'env.json')

        print(f"Running {test_name}...")
        print(f"Starting mock server on port {mock_port} with start_config.json for {test_name}...")
        mock_aikido_core = subprocess.Popen(['python3', 'mock_aikido_core.py', str(mock_port), config_path])
        time.sleep(5)

        print(f"Starting PHP server on port {php_port} for {test_name}...")
        env = os.environ.copy()
        env.update({
            'AIKIDO_LOG_LEVEL': 'DEBUG' if debug else 'ERROR',
            'AIKIDO_TOKEN': 'AIK_RUNTIME_MOCK',
            'AIKIDO_ENDPOINT': f'http://localhost:{mock_port}/',
            'AIKIDO_REALTIME_ENDPOINT': f'http://localhost:{mock_port}/',
        })
        env.update(load_env_from_json(env_file_path))

        php_server_process_cmd = ['php', '-S', f'localhost:{php_port}', '-t', test_dir]
        if valgrind:
            php_server_process_cmd = ['valgrind', f'--supressions={test_lib_dir}/valgrind.supp'] + php_server_process_cmd
            
        php_server_process = subprocess.Popen(
            php_server_process_cmd,
            env=env
        )
        time.sleep(5)

        test_script_name = "test.py"
        test_script_cwd = test_dir
        if benchmark:
            print(f"Running benchmark for {test_name}...")
            test_script_name = "benchmark.py"
            test_script_cwd = root_tests_dir
        else:
            print(f"Running test.py for {test_name}...")
            
        subprocess.run(['python3', test_script_name, str(php_port), str(mock_port), test_name], 
                       env=dict(os.environ, PYTHONPATH=f"{test_lib_dir}:$PYTHONPATH"),
                       cwd=test_script_cwd,
                       check=True, timeout=600, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        
        passed_tests.append(test_name)

    except subprocess.CalledProcessError as e:
        print(f"Error in testing scenario {test_name}:")
        print(f"Exception output: {e.output}")
        print(f"Test exit code: {e.returncode}")
        print(f"Test stdout: {e.stdout.decode()}")
        print(f"Test stderr: {e.stderr.decode()}")
        failed_tests.append(test_name)
        
    except subprocess.TimeoutExpired:
        print(f"Error in testing scenario {test_name}:")
        print(f"Execution timed out.")
        failed_tests.append(test_name)
        
    finally:
        if php_server_process:
            php_server_process.terminate()
            php_server_process.wait()
            print(f"PHP server on port {php_port} stopped.")

        if mock_aikido_core:
            mock_aikido_core.terminate()
            mock_aikido_core.wait()
            print(f"Mock server on port {mock_port} stopped.")


def main(root_tests_dir, test_lib_dir, specific_test=None, benchmark=False, valgrind=False, debug=False):
    if specific_test:
        specific_test = os.path.join(root_tests_dir, specific_test)
        handle_test_scenario(root_tests_dir, specific_test, test_lib_dir, benchmark, valgrind, debug)
    else:
        run_parallel = True
        if benchmark or valgrind:
            run_parallel = False
            
        test_dirs = [f.path for f in os.scandir(root_tests_dir) if f.is_dir()]
        threads = []
        
        for test_dir in test_dirs:
            args = (root_tests_dir, test_dir, test_lib_dir, benchmark, valgrind, debug)
            if run_parallel:
                thread = threading.Thread(target=handle_test_scenario, args=args)
                threads.append(thread)
                thread.start()
            else:
                handle_test_scenario(*args)
        
        if run_parallel:
            for thread in threads:
                thread.join()
            
    print_test_results("Passed tests:", passed_tests)
    print_test_results("Failed tests:", failed_tests)
    assert failed_tests == [], f"Found failed tests: {failed_tests}"
    print("All tests passed!")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Script for running PHP server tests with Aikido Firewall installed.")
    parser.add_argument('root_folder_path', type=str, help='Path to the root folder of the tests to be ran.')
    parser.add_argument('test_lib_dir', type=str, help='Directory for the test libraries.')
    parser.add_argument('--test', type=str, default=None, help='Run a single test from the root folder.')
    parser.add_argument('--benchmark', action='store_true', help='Enable benchmarking.')
    parser.add_argument('--valgrind', action='store_true', help='Enable valgrind.')
    parser.add_argument('--debug', action='store_true', help='Enable debugging logs.')

    # Parse arguments
    args = parser.parse_args()

    # Extract values from parsed arguments
    root_folder = os.path.abspath(args.root_folder_path)
    test_lib_dir = os.path.abspath(args.test_lib_dir)
    main(root_folder, test_lib_dir, args.test, args.benchmark, args.valgrind, args.debug)
