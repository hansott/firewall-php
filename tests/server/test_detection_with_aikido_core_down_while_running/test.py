import requests
import time
import sys
from testlib import *

'''
1. Start with normal config
2. Do a path traversal attack and check if event was submitted.
3. Get mock server down (to simulate that aikido core went down while running).
4. Do another path traversal. Check that is blocked but not reported.
5. Get mock server up (to simulate aikido core coming back up).
6. Do another path traversal. Check that it is blocked and reported.
'''

def run_test():
    assert_response_code_is(php_server_post("/testDetection", {"folder": "../../../.."}), 500)
    
    assert_events_length_is(mock_server_get_events(), 2)
    
    mock_server_down()
    
    time.sleep(70)
    
    assert_events_length_is(mock_server_get_events(), 2)
    assert_response_code_is(php_server_post("/testDetection", {"folder": "../../../.."}), 500)
    assert_events_length_is(mock_server_get_events(), 2)
    
    mock_server_up()
    
    assert_response_code_is(php_server_post("/testDetection", {"folder": "../../../.."}), 500)
    assert_events_length_is(mock_server_get_events(), 3)
        
if __name__ == "__main__":
    load_test_args()
    run_test()
