import os
import threading
import subprocess
import random
import time
import sys
import json
import argparse
from server_tests.php_built_in.main import handle_php_built_in
from server_tests.apache.main import handle_apache_mod_php
from server_tests.nginx.main import prepare_nginx_php_fpm, pre_nginx_php_fpm, handle_nginx_php_fpm, done_nginx_php_fpm

server_prepare_handlers = {
    "php-built-in": None,
    "apache-mod-php": None,
    "nginx-php-fpm": prepare_nginx_php_fpm
}

server_pre_tests_handlers = {
    "php-built-in": None,
    "apache-mod-php": None,
    "nginx-php-fpm": pre_nginx_php_fpm
}

server_handlers = {
    "php-built-in": handle_php_built_in,
    "apache-mod-php": handle_apache_mod_php,
    "nginx-php-fpm": handle_nginx_php_fpm
}

server_done_handlers = {
    "php-built-in": None,
    "apache-mod-php": None,
    "nginx-php-fpm": done_nginx_php_fpm
}

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


def handle_test_scenario(data, root_tests_dir, test_lib_dir, server, benchmark, valgrind, debug):
    test_name = data["test_name"]
    mock_port = data["mock_port"]
    server_port = data["server_port"]
    try:
        print(f"Running {test_name}...")
        print(f"Starting mock server on port {mock_port} with start_config.json for {test_name}...")
        mock_aikido_core = subprocess.Popen(["python3", "mock_aikido_core.py", str(mock_port), data["config_path"]])
        time.sleep(5)

        print(f"Starting {server} server on port {server_port} for {test_name}...")
        
        server_processes = server_handlers[server](data, test_lib_dir, valgrind)

        time.sleep(5)

        test_script_name = "test.py"
        test_script_cwd = data["test_dir"]
        if benchmark:
            print(f"Running benchmark for {test_name}...")
            test_script_name = "benchmark.py"
            test_script_cwd = root_tests_dir
        else:
            print(f"Running test.py for {test_name}...")
            
        subprocess.run(["python3", test_script_name, str(server_port), str(mock_port), test_name], 
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
        for p in server_processes:
            p.terminate()
            p.wait()
            print(f"PHP server on port {server_port} stopped.")

        if mock_aikido_core:
            mock_aikido_core.terminate()
            mock_aikido_core.wait()
            print(f"Mock server on port {mock_port} stopped.")


def main(root_tests_dir, test_lib_dir, specific_test=None, server="php-built-in", benchmark=False, valgrind=False, debug=False):    
    if specific_test:
        test_dirs = [os.path.join(root_tests_dir, specific_test)]
    else:
        test_dirs = [f.path for f in os.scandir(root_tests_dir) if f.is_dir()]
        
    tests_data = []
    for test_dir in test_dirs:
        mock_port = generate_unique_port()
        test_data = {
            "test_name": os.path.basename(os.path.normpath(test_dir)),
            "test_dir": test_dir,
            "mock_port": mock_port,
            "server_port": generate_unique_port(),
            "config_path": os.path.join(test_dir, "start_config.json"),
            "env_path": os.path.join(test_dir, "env.json")
        }

        env = {
            "AIKIDO_LOG_LEVEL": "DEBUG" if debug else "ERROR",
            "AIKIDO_TOKEN": "AIK_RUNTIME_MOCK",
            "AIKIDO_ENDPOINT": f"http://localhost:{mock_port}/",
            "AIKIDO_REALTIME_ENDPOINT": f"http://localhost:{mock_port}/",
        }
        env.update(load_env_from_json(test_data["env_path"]))
        test_data["env"] = env
        
        if server_prepare_handlers[server] is not None:
            test_data = server_prepare_handlers[server](test_data)
        tests_data.append(test_data)
            
    if server_pre_tests_handlers[server] is not None:
        test_data = server_pre_tests_handlers[server]()

            
    threads = []
    for test_data in tests_data:
        args = (test_data, root_tests_dir, test_lib_dir, server, benchmark, valgrind, debug)
        thread = threading.Thread(target=handle_test_scenario, args=args)
        threads.append(thread)
        thread.start()

    for thread in threads:
        thread.join()
        
    if server_done_handlers[server]:
        server_done_handlers[server]()
            
    print_test_results("Passed tests:", passed_tests)
    print_test_results("Failed tests:", failed_tests)
    assert failed_tests == [], f"Found failed tests: {failed_tests}"
    print("All tests passed!")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Script for running PHP server tests with Aikido Firewall installed.")
    parser.add_argument("root_folder_path", type=str, help="Path to the root folder of the tests to be ran.")
    parser.add_argument("test_lib_dir", type=str, help="Directory for the test libraries.")
    parser.add_argument("--test", type=str, default=None, help="Run a single test from the root folder.")
    parser.add_argument("--benchmark", action="store_true", help="Enable benchmarking.")
    parser.add_argument("--valgrind", action="store_true", help="Enable valgrind.")
    parser.add_argument("--debug", action="store_true", help="Enable debugging logs.")
    parser.add_argument("--server", type=str, choices=["php-built-in", "apache-mod-php", "nginx-php-fpm"], default="php-built-in", help="Enable nginx & php-fpm testing.")

    # Parse arguments
    args = parser.parse_args()

    # Extract values from parsed arguments
    root_folder = os.path.abspath(args.root_folder_path)
    test_lib_dir = os.path.abspath(args.test_lib_dir)
    main(root_folder, test_lib_dir, args.test, args.server, args.benchmark, args.valgrind, args.debug)
