import requests
import time
import sys
from testlib import *

'''
1. Sets up the rate limiting config to 30 requests / 10 minutes.
2. Sends 10 requests once every minute, for 3 minutes. Checks that those requests are not blocked.
3. Send another more 10 request. Checks that they all are rate limited.
'''

def run_test(php_port, mock_port):
    for i in range(30):
        response = php_server_get(php_port, "/test")
        assert_response_code_is(response, 200)
        assert_reponse_body_contains(response, "Something")
        
        if i != 0 and i % 10 == 0:
            time.sleep(60)
        
    for _ in range(10):
        response = php_server_get(php_port, "/test")
        assert_response_code_is(response, 429)
        assert_reponse_header_contains(response, "Content-Type", "text")
        assert_reponse_body_contains(response, "This request was rate limited by Aikido Security!")
    
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
