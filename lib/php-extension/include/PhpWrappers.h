#pragma once

bool CallPhpEcho(std::string s);

bool CallPhpFunction(std::string function_name, unsigned int params_number = 0, zval *params = nullptr, zval *return_value = nullptr, zval *object = nullptr);

bool CallPhpFunctionWithOneParam(std::string function_name, long first_param, zval *return_value = nullptr, zval *object = nullptr);

bool CallPhpFunctionWithOneParam(std::string function_name, std::string first_param, zval *return_value = nullptr, zval *object = nullptr);

std::string CallPhpFunctionCurlGetInfo(zval *curl_handle, int curl_info_option);
