import requests
import time
import sys
from testlib import *
import json

'''
1. Sets up a simple config.
2. Sends one request that will trigger multiple curl reuqests from php.
3. Waits for the heartbeat event and validates it.
'''

def run_test():
    response = php_server_post("/testDetection", json.load(open("test.json", 'r')))
    assert_response_code_is(response, 200)
    
if __name__ == "__main__":
    load_test_args()
    run_test()
