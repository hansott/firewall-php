#include "GoWrappers.h"

std::string CreateCString(GoString g) {
    std::string s;
    s.assign(g.p, g.n);
    return s;
}

GoString CreateGoString(std::string& s) {
    return GoString { s.c_str(), s.length() };
}

std::string GetHostname(std::string& url) {
    return CreateCString(AikidoLibNormalizeDomain(CreateGoString(url)));
}