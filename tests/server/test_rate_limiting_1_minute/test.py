import requests
import time
import sys
from testlib import *

'''
1. Sets up the rate limiting config to 5 requests / minute for route '/'.
2. Sends 5 requests to '/'. Checks that those requests are not blocked.
3. Send another more 5 request to '/'. Checks that they all are rate limited.
4. Sends 100 requests to another route '/tests'. Checks that those requests are not blocked.
'''

def run_test(php_port, mock_port):
    for _ in range(5):
        response = php_server_get(php_port, "/")
        assert_response_code_is(response, 200)
        
    time.sleep(10)
        
    for _ in range(5):
        response = php_server_get(php_port, "/")
        assert_response_code_is(response, 429)
        assert_response_header_contains(response, "Content-Type", "text")
        assert_response_body_contains(response, "This request was rate limited by Aikido Security!")
    
    for _ in range(100):
        response = php_server_get(php_port, "/test")
        assert_response_code_is(response, 200)
        
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
