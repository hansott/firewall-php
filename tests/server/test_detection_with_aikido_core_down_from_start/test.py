import requests
import time
import sys
from testlib import *

'''
1. Simulate aikido-core is down by setting AIKIDO_ENDPOINT & AIKIDO_REALTIME_ENDPOINT to some invalid values.
2. Check that path traversal detection is still emitted and blocked.
3. Check no events are submitted.
'''

def run_test():
    assert_response_code_is(php_server_post("/testDetection", {"folder": "../../../.."}), 500)
    
    events = mock_server_get_events()
    assert_events_length_is(events, 0)
    
    time.sleep(65)
    
    assert_response_code_is(php_server_post("/testDetection", {"folder": "../../../.."}), 500)
    assert_events_length_is(events, 0)
        
if __name__ == "__main__":
    load_test_args()
    run_test()
