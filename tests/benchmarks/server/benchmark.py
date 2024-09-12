import requests
import time
import sys
from testlib import *
import json


def run_benchmark():
    for _ in range(1000):
        response = php_server_post("/test", {})
        
    for _ in range(10000):
        response = php_server_post("/test", {}, benchmark=True)
        assert_response_code_is(response, 200)
    
if __name__ == "__main__":
    load_test_args()
    run_benchmark()
    store_benchmark_results()
