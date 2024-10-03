import requests
import time
import sys
from testlib import *
import json


def run_benchmark():
    for _ in range(10000):
        response = php_server_post("/test?userId=12345&status=pending&orderId=80", 
                                   { "orderPlaced": True },
                                   { "Content-Type": "application/json", "X-API-Key": "abcdef12345" },
                                   benchmark=True)
        assert_response_code_is(response, 200)
    
if __name__ == "__main__":
    load_test_args()
    benchmark_warmup()
    run_benchmark()
    benchmark_store_results()
