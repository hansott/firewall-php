import requests
import time
import sys
from testlib import *

def run_test(php_port, mock_port):
    php_server_get(php_port, "/")
    
    time.sleep(60 + 10)
    
    events = mock_server_get_events(mock_port)
    print(events)
    assert_events_length_is(events, 2)
    assert_started_event_is_valid(events[0])
    assert_event_contains_subset_file(events[1], "expect_user.json")
    
    print("All assertions passed successfully.")

if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
