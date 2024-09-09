import requests
import time
import sys
from testlib import *

'''
1. Sets up a simple config and env AIKIDO_BLOCKING=1.
2. Sends 200 attack requests to a route, that will cause sending a detection event.
3. Checks that there are no more than 100 detection events submited.
'''

def run_test(php_port, mock_port):
    for _ in range(200):
        response = php_server_post(php_port, "/testDetection", {"folder": "../../../.."})
        assert_response_code_is(response, 500)   
        time.sleep(0.01)
        
    time.sleep(5)
        
    events = mock_server_get_events(mock_port)
    assert_events_length_is(events, 101)
    assert_started_event_is_valid(events[0])
    for e in events[1:101]:
        assert_event_contains_subset(e, {"type": "detected_attack"})

    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
