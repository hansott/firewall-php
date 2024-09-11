import requests
import time
import sys
from testlib import *

'''
1. Sets up a config with receivedAnyStats = false (so heartbeat will be sent after 1 minute) and a user id to be blocked.
2. Sends a get request (on the PHP side a user will be set). Checks that the user is blocked.
3. Changes to config so the user is no longer blocked.
4. Sends a get request. Checks that the user is not blocked anymore.
5. Repeats steps 1-3.
'''

def run_test():
    response = php_server_get("/test")
    assert_response_code_is(response, 403)
    assert_response_header_contains(response, "Content-Type", "text")
    assert_response_body_contains(response, "You are blocked by Aikido Firewall!")

    apply_config("change_config_remove_blocked_user.json")
        
    response = php_server_get("/test")
    assert_response_code_is(response, 200)
    assert_response_body_contains(response, "User set successfully")
    
    apply_config("start_config.json")
        
    response = php_server_get("/test")
    assert_response_code_is(response, 403)
    assert_response_header_contains(response, "Content-Type", "text")
    assert_response_body_contains(response, "You are blocked by Aikido Firewall!")
    
    
if __name__ == "__main__":
    load_test_args()
    run_test()
