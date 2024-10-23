#pragma once

std::string CppCreateString(GoString);

GoString GoCreateString(std::string&);

char* GoContextCallback(int callbackId);
