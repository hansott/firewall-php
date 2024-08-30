import requests
import time
import sys
from testlib import *

def run_test(php_port, mock_port):
    for _ in range(5):
        response = php_server_get(php_port, "/")
        assert_response_code_is(response, 200)
        
    time.sleep(0.5)
        
    for _ in range(5):
        response = php_server_get(php_port, "/")
        assert_response_code_is(response, 429)
        assert_reponse_header_contains(response, "Content-Type", "text")
        assert_reponse_body_contains(response, "This request was rate limited by Aikido Security!")
    
    for _ in range(100):
        response = php_server_get(php_port, "/test")
        assert_response_code_is(response, 200)
        
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
