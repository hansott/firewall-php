#include "Includes.h"
#include "Actions.h"

Action action;

ACTION_STATUS Action::executeThrow(json &event)
{
    int _response_code = event["response_code"].get<int>();
    std::string _message = event["message"].get<std::string>();
    zend_throw_exception(zend_exception_get_default(), _message.c_str(), _response_code);
    return BLOCK;
}

ACTION_STATUS Action::executeExit(json &event)
{
    int _response_code = event["response_code"].get<int>();
    std::string _message = event["message"].get<std::string>();

    // aikido_call_user_function("ob_clean");
    aikido_call_user_function("header_remove");
    aikido_call_user_function_one_param("http_response_code", _response_code);
    aikido_call_user_function_one_param("header", "Content-Type: text/plain");
    aikido_echo(_message);
    exit = true;
    return EXIT;
}

ACTION_STATUS Action::executeStore(json &event)
{
    block = true;
    type = event["type"];
    trigger = event["trigger"];
    ip = event["ip"];
    return CONTINUE;
}

ACTION_STATUS Action::Execute(std::string &event)
{
    if (event.empty())
    {
        return CONTINUE;
    }

    json eventJson = json::parse(event);
    std::string actionType = eventJson["action"];
    
    if (actionType == "throw")
    {
        return executeThrow(eventJson);
    }
    else if (actionType == "exit")
    {
        return executeExit(eventJson);
    }
    else if (actionType == "store")
    {
        return executeStore(eventJson);
    }
    return CONTINUE;
}

void Action::Reset()
{
    exit = false;
    block = false;
    type = "";
    trigger = "";
    ip = "";
}

bool Action::Exit() {
    return exit;
}

bool Action::Block() {
    return block;
}

char* Action::Type() {
    return (char*)type.c_str();
}

char* Action::Trigger() {
    return (char*)trigger.c_str();
}

char* Action::Ip() {
    return (char*)ip.c_str();
}
