import requests
import time
import sys
from testlib import *

'''
1. Sets up a simple config and env AIKIDO_BLOCKING=1.
2. Sends an attack request to a route, that will cause sending a detection event.
3. Checks that the detection event was submitted and is valid.
'''



def check_path_traversal(php_port, mock_port, response_code, response_body=""):
    response = php_server_post(php_port, "/testDetection", {"folder": "../../../.."})
    assert_response_code_is(response, response_code)
    assert_response_body_contains(response, response_body)
    
    mock_server_wait_for_new_events(mock_port, 5)
    
    events = mock_server_get_events(mock_port)
    assert_events_length_is(events, 2)
    assert_started_event_is_valid(events[0])
    assert_event_contains_subset_file(events[1], "expect_detection.json")

def run_test(php_port, mock_port):
    check_path_traversal(php_port, mock_port, 500)
    
    apply_config(mock_port, "change_config_disable_blocking.json")
    check_path_traversal(php_port, mock_port, 200, "File opened!")
    
    apply_config(mock_port, "start_config.json")
    check_path_traversal(php_port, mock_port, 500)
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
