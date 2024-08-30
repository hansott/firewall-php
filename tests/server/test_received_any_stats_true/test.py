import requests
import time
import sys
from testlib import *

'''
1. Sets up receiveAnyStats to true.
2. Checks that the 'started' event is valid.
3. After 1 minute, verifies that the hearbeat events was not yet sent.
'''

def run_test(php_port, mock_port):
    php_server_get(php_port, "/")
    
    events = mock_server_get_events(mock_port)
    assert_events_length_is(events, 1)
    assert_started_event_is_valid(events[0])
    
    mock_server_wait_for_new_events(mock_port, 70)
    
    events = mock_server_get_events(mock_port)
    assert_events_length_is(events, 1)
    
if __name__ == "__main__":
    run_test(int(sys.argv[1]), int(sys.argv[2]))
