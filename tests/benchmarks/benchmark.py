import requests
import time
import sys
from testlib import *
import json


def run_benchmark():
    for _ in range(1):
        response = php_server_post("/" + generate_random_string(20), {})
        assert_response_code_is(response, 200)
    
if __name__ == "__main__":
    load_test_args()
    run_benchmark()
    store_benchmark_results()
