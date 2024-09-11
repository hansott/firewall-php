import requests
import time
import sys
from testlib import *

import requests
import time
import sys
from testlib import *

'''
1. Sets up a simple config with receivedAnyStats = false (so heartbeat will be sent after 1 minute).
2. Sends a get request (on the PHP side a user will be set).
3. Waits for 1 minute and checks if the user is present in the hearbeat request.
'''

def run_test():
    php_server_get("/")
    
    mock_server_wait_for_new_events(70)
    
    events = mock_server_get_events()
    assert_events_length_is(events, 2)
    assert_started_event_is_valid(events[0])
    assert_event_contains_subset_file(events[1], "expect_user.json")
    
if __name__ == "__main__":
    load_test_args()
    run_test()
