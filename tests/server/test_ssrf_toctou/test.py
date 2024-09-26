import requests
import time
import sys
from testlib import *

'''
1. Sets up a simple config and env AIKIDO_BLOCKING=1.
2. Sends an attack request to a route, that will cause sending a detection event.
3. Checks that the detection event was submitted and is valid.
'''

def check_ssrf(response_code, response_body, event_id, expected_json):
    response = php_server_post("/testDetection", {"url": f"http://app.example.local:{get_mock_port()}/mock/events"})
    assert_response_code_is(response, response_code)
    assert_response_body_contains(response, response_body)
    
    mock_server_wait_for_new_events(5)
    
    events = mock_server_get_events()
    assert_events_length_is(events, event_id + 1)
    assert_started_event_is_valid(events[0])
    assert_event_contains_subset_file(events[event_id], expected_json)

def run_test():
    add_to_hosts_file("app.example.local", "127.0.0.1")
    
    check_ssrf(500, "", 1, "expect_detection_blocked.json")
    
    apply_config("change_config_disable_blocking.json")
    check_ssrf(200, "Got URL content!", 2, "expect_detection_not_blocked.json")
    
    apply_config("start_config.json")
    check_ssrf(500, "", 3, "expect_detection_blocked.json")
    
if __name__ == "__main__":
    load_test_args()
    run_test()
