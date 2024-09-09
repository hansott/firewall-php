import requests
import time
import sys
from testlib import *

'''
1. Sets up the bypassed IP address config for route '/test'. Rate limiting is set to 10 req / min. Checks that requests are not rate limited blocked.
2. Changes the config to remove the bypassed IP address. Checks that requests are rate limiting.
3. Changes the config again to enable the bypassed IP address. Checks that requests are not rate limited blocked.
'''


def run_test(php_port, mock_port):
    for _ in range(100):
        response = php_server_get(php_port, "/test")
        assert_response_code_is(response, 200)
        
    apply_config(mock_port, "change_config_remove_bypassed_ip.json")

    for i in range(100):
        response = php_server_get(php_port, "/test")
        if i < 10:
            assert_response_code_is(response, 200)
        else:
            assert_response_code_is(response, 429)
            assert_response_header_contains(response, "Content-Type", "text")
            assert_response_body_contains(response, "This request was rate limited by Aikido Security!")

    apply_config(mock_port, "start_config.json")
    
    for _ in range(100):
        response = php_server_get(php_port, "/test")
        assert_response_code_is(response, 200)
    
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
