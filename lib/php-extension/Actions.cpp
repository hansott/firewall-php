#include "Includes.h"

ACTION aikido_execute_output_throw(json event) {
    std::string message = event["message"].get<std::string>();
    int code = event["code"].get<int>();
    zend_throw_exception(zend_exception_get_default(), message.c_str(), code);
    return BLOCK;
}

ACTION aikido_execute_output_exit(json event) {
    int _response_code = event["response_code"].get<int>();
    std::string _message = event["message"].get<std::string>();

    //aikido_call_user_function("ob_clean");
    aikido_call_user_function("header_remove");
    aikido_call_user_function_one_param("http_response_code", _response_code);
    aikido_call_user_function_one_param("header", "Content-Type: text/plain");
    aikido_echo(_message);
    aikido_exit();
    return EXIT;
}

ACTION aikido_execute_output(std::string event) {
    if (event.empty()) {
        return CONTINUE;
    }

    json eventJson = json::parse(event);
    if (eventJson["action"] == "throw")
        return aikido_execute_output_throw(eventJson);
    if (eventJson["action"] == "exit")
        return aikido_execute_output_exit(eventJson);
    return CONTINUE;
}