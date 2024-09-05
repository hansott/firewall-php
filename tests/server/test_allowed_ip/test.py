import requests
import time
import sys
from testlib import *

'''
1. Sets up the allowed IP address config for route '/test'. Checks that requests are blocked.
2. Changes the config to remote the allowed IP address. Checks that requests are passing.
3. Changes the config again to enable allowed IP address. Checks that requests are blocked.
'''


def run_test(php_port, mock_port):
    response = php_server_get(php_port, "/test")
    assert_response_code_is(response, 403)
    assert_reponse_header_contains(response, "Content-Type", "text")
    assert_reponse_body_contains(response, "Your IP address is not allowed to access this resource! (Your IP: ::1)")

    apply_config(mock_port, "change_config_remove_allowed_ip.json")
        
    response = php_server_get(php_port, "/test")
    assert_response_code_is(response, 200)
    assert_reponse_body_contains(response, "Something")
    
    apply_config(mock_port, "start_config.json")
        
    response = php_server_get(php_port, "/test")
    assert_response_code_is(response, 403)
    assert_reponse_header_contains(response, "Content-Type", "text")
    assert_reponse_body_contains(response, "Your IP address is not allowed to access this resource! (Your IP: ::1)")
    
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
