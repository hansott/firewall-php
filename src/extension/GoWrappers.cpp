#include "GoWrappers.h"

std::string CppCreateString(GoString g) {
    std::string s;
    s.assign(g.p, g.n);
    return s;
}

GoString GoCreateString(std::string& s) {
    return GoString { s.c_str(), s.length() };
}

json GoOnEvent(json& event) {
    std::string eventString = event.dump();
    
    php_printf("[AIKIDO-C++] Seding event to GO\n");
    
    std::string outputString = CppCreateString(OnEvent(GoCreateString(eventString)));
    
    php_printf("[AIKIDO-C++] Got event reply: %s\n", outputString.c_str());
    
    json output = json::parse(outputString);
    return output;
}