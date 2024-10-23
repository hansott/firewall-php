import requests
import time
import sys
from testlib import *

def run_test():
    for userId, userName in [("123", "Tudor"), ("0", "Test"), ("50", "Willem")]:
        for _ in range(25):
            response = php_server_post("/", { "userId": userId, "userName": userName })
            assert_response_code_is(response, 200)
            assert_response_header_contains(response, "Content-Type", "text")
            assert_response_body_contains(response, "Request successful")
                
    
if __name__ == "__main__":
    load_test_args()
    run_test()
