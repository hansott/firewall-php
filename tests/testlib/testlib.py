import requests
import time
import sys
import json
 
php_port = 0
mock_port = 0

def load_ports_from_args():
    global php_port, mock_port
    php_port = int(sys.argv[1])
    mock_port = int(sys.argv[2])
    print(f"Loaded ports: php_port={php_port}, mock_port={mock_port}")

def localhost_get_request(port, route=""):
    r = requests.get(f"http://localhost:{port}{route}")
    time.sleep(0.01)
    return r

def localhost_post_request(port, route, data):
    r = requests.post(f"http://localhost:{port}{route}", json=data)
    time.sleep(0.01)
    return r

def php_server_get(route=""):
    return localhost_get_request(php_port, route)

def php_server_post(route, data):
    return localhost_post_request(php_port, route, data)

def mock_server_get(route=""):
    return localhost_get_request(mock_port, route)

def mock_server_post(route, data):
    return localhost_post_request(mock_port, route, data)

def mock_server_get_events():
    return mock_server_get("/mock/events").json()

def mock_server_set_config(config):
    return mock_server_post("/mock/config", config)

def mock_server_set_config_file(config_file):
    config = None
    with open(config_file, 'r') as f:
        config = json.load(f)
    return mock_server_post("/mock/config", config)

def apply_config(config_file):
    mock_server_set_config_file(config_file)
    time.sleep(120)

def assert_events_length_is(events, length):
    assert isinstance(events, list), "Error: Events is not a list."
    assert len(events) == length, f"Error: Events list contains {len(events)} elements and not {length} elements."

def assert_event_contains_subset(event, event_subset, dry_mode=False):
    """
    Recursively checks that all keys and values in the subset JSON exist in the event JSON
    and have the same values. If a key in the subset is a list, all its elements must exist in the
    corresponding list in the event.

    :param event: The event JSON dictionary
    :param subset: The subset JSON dictionary
    :raises AssertionError: If the subset is not fully contained within the event
    """
    def result(assertion_error):
        if dry_mode:
            return False
        raise assertion_error
    
    #print(f"Search {event_subset} in {event} (dry_mode = {dry_mode})")
    
    if isinstance(event_subset, dict):
        found_all_keys = True
        for key, value in event_subset.items():
            if key not in event:
                return result(AssertionError(f"Key '{key}' not found in '{event}'."))
            if not assert_event_contains_subset(event[key], value, dry_mode):
                found_all_keys = False
        return found_all_keys
    elif isinstance(event_subset, list):
        if not isinstance(event, list):
            return result(AssertionError(f"Expected a list in event but found '{event}'."))
        for event_subset_item in event_subset:
            found_item = False
            for event_item in event:
                if assert_event_contains_subset(event_item, event_subset_item, dry_mode=True):
                    found_item = True
                    break
            if not found_item:
                return result(AssertionError(f"Item '{event_subset_item}' not found in {event}."))
    else:
        if event_subset != event:
            return result(AssertionError(f"Value mismatch: {event_subset} != {event}"))
        
    return True
        
def assert_event_contains_subset_file(event, event_subset_file):
    event_subset = None
    with open(event_subset_file, 'r') as file:
        event_subset = json.load(file)
    assert event_subset
    print(event_subset, event)
    assert_event_contains_subset(event, event_subset)

def assert_started_event_is_valid(event):
    assert_event_contains_subset(event, {"type": "started", "agent": { "library": "firewall-php" } })
    
def assert_response_code_is(response, status_code):
    assert response.status_code == status_code, f"Status codes are not the same: {response.status_code} vs {status_code}"
    
def assert_response_header_contains(response, header, value):
    assert header in response.headers, f"Header '{header}' is not part of response headers: {response.headers}"
    assert value in response.headers[header], f"Header '{header}' does not contain '{value}' but '{response.headers[header]}'"

def assert_response_body_contains(response, text):
    assert text in response.text, f"Test '{text}' is not part of response body: {response.text}"
    
def mock_server_wait_for_new_events(max_wait_time):
    initial_number_of_events = len(mock_server_get_events())
    while max_wait_time > 0:
        if len(mock_server_get_events()) > initial_number_of_events:
            return True
        time.sleep(5)
        max_wait_time -= 5
        
    return False