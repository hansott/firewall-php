import requests
import time
import sys
import json

def localhost_get_request(port, route=""):
    return requests.get(f"http://localhost:{port}{route}")

def php_server_get(port, route=""):
    return localhost_get_request(port, route)

def mock_server_get(port, route=""):
    return localhost_get_request(port, route)

def mock_server_get_events(port):
    return mock_server_get(port, "/mock/events").json()

def assert_events_length_is(events, length):
    assert isinstance(events, list), "Error: Events is not a list."
    assert len(events) == length, f"Error: The events list does not contain exactly {length} elements."

def assert_event_contains_subset(event, event_subset):
    """
    Recursively checks if event_subset is a subset of event.
    
    :param event: The superset JSON.
    :param event_subset: The subset JSON.
    :return: True if event_subset is a subset of event, otherwise False.
    """
    if isinstance(event_subset, dict) and isinstance(event, dict):
        # Check if all keys and their corresponding values in event_subset are in event
        return all(key in event and assert_event_contains_subset(event_subset[key], event[key]) for key in event_subset)

    elif isinstance(event_subset, list) and isinstance(event, list):
        # Check if all elements of event_subset are in event
        return all(any(assert_event_contains_subset(item, subitem) for subitem in event) for item in event_subset)

    else:
        print(f"Checking {event_subset} == {event}...")
        # Base case for non-iterable types
        assert event_subset == event
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