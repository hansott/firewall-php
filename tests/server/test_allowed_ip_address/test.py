import requests
import time
import sys
from testlib import *

def run_test(php_port, mock_port):
    response = php_server_get(php_port, "/")
    assert_response_code_is(response, 403)
    assert_reponse_header_contains(response, "Content-Type", "text")
    assert_reponse_body_contains(response, "Your IP address is not allowed to access this resource!")    
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
