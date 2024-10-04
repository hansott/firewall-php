import requests
import time
import sys
from testlib import *


def run_test():        
    response = php_server_post("/", {})
    assert_response_code_is(response, 200)
    assert_response_body_contains(response, "test_instance_id")
    
    events = mock_server_get_events()
    assert_events_length_is(events, 1)
    assert_started_event_is_valid(events[0])
    
    response = php_server_post("/", {"url": "169.254.169.254"})
    assert_response_code_is(response, 500)
    
    mock_server_wait_for_new_events(5)
    
    events = mock_server_get_events()
    assert_events_length_is(events, 2)
    assert_started_event_is_valid(events[0])
    assert_event_contains_subset_file(events[1], "expect_detection_blocked.json")
        
if __name__ == "__main__":
    load_test_args()
    run_test()
