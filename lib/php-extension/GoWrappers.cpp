#include "GoWrappers.h"

std::string CppCreateString(GoString g) {
    std::string s;
    s.assign(g.p, g.n);
    return s;
}

GoString GoCreateString(std::string& s) {
    return GoString { s.c_str(), s.length() };
}

json GoRequestProcessorOnEvent(json& event) {
    std::string eventString = event.dump();
    
    AIKIDO_LOG_DEBUG("Sending event to GO\n");
    
    std::string outputString = CppCreateString(request_processor_on_event_fn(GoCreateString(eventString)));
    
    AIKIDO_LOG_DEBUG("Got event reply: %s\n", outputString.c_str());
    
    json output = json::parse(outputString);
    return output;
}
