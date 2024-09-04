import requests
import time
import sys
from testlib import *

'''
1. Sets up a simple config.
2. Sends multiple requests to different routes.
3. Waits for the heartbeat event and validates the reporting.
'''

def run_test(php_port, mock_port):
    response = php_server_post(php_port, "/testDetection", {"folder": "../../../.."})
    assert_response_code_is(response, 500)
    
    mock_server_wait_for_new_events(mock_port, 5)
    
    events = mock_server_get_events(mock_port)
    assert_events_length_is(events, 2)
    assert_started_event_is_valid(events[0])
    assert_event_contains_subset_file(events[1], "expect_detection.json")
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
