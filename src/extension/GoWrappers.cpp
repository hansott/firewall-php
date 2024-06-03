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
    
    AIKIDO_LOG_DEBUG("Sending event to GO\n");
    
    std::string outputString = CppCreateString(OnEvent(GoCreateString(eventString)));
    
    AIKIDO_LOG_DEBUG("Got event reply: %s\n", outputString.c_str());
    
    json output = json::parse(outputString);
    return output;
}

bool GoInit(json& initData) {
    std::string initDataString = initData.dump();
    
    AIKIDO_LOG_DEBUG("Sending init data to GO\n");
    
    bool initOk = Init(GoCreateString(initDataString));
    
    AIKIDO_LOG_DEBUG("Got init status: %d\n", initOk);
    
    return initOk;
}

void GoUninit() {
    Uninit();
}
