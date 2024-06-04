#include "Utils.h"

std::string to_lowercase(const std::string& str) {
    std::string result = str;
    std::transform(result.begin(), result.end(), result.begin(), [](unsigned char c){ return std::tolower(c); });
    return result;
}

void aikido_log(AIKIDO_LOG_LEVEL level, const char* format, ...) {
    if (level > AIKIDO_GLOBAL(log_level)) {
        return;
    }

    printf("[AIKIDO][%s][C++] ", aikido_log_level_str(level));

    va_list args;
    va_start(args, format);
    vprintf(format, args);
    va_end(args);
}

const char* aikido_log_level_str(AIKIDO_LOG_LEVEL level) {
    switch (level) {
        case AIKIDO_LOG_LEVEL_DEBUG:
            return "DEBUG";
        case AIKIDO_LOG_LEVEL_INFO:
            return "INFO";
        case AIKIDO_LOG_LEVEL_WARN:
            return "WARN";
        case AIKIDO_LOG_LEVEL_ERROR:
            return "ERROR";
    }
    return "UNKNOWN";
}

std::string get_environment_variable(const std::string& env_key) {
    const char* env_value = getenv(env_key.c_str());
    if (!env_value) return "";
    return env_value;
}

std::string config_override_with_env(const std::string previous_value, const std::string& env_key) {
	std::string env_value = get_environment_variable(env_key.c_str());
	if (!env_value.empty()) {
        return env_value;
	}
    return previous_value;
}

std::string get_hostname() {
    char hostname[HOST_NAME_MAX + 1];

    if (gethostname(hostname, sizeof(hostname)) == 0) {
        return hostname;
    }

    return "";
}

utsname get_os_info() {
    utsname buffer = {};
    uname(&buffer);
    return buffer;
}

std::string get_ip_address() {
    std::string ip_address;

    ifaddrs* ifAddrStruct = nullptr;
    ifaddrs* ifa = nullptr;
    void* tmpAddrPtr = nullptr;

    getifaddrs(&ifAddrStruct);

    for (ifa = ifAddrStruct; ifa != NULL; ifa = ifa->ifa_next) {
        if (!ifa->ifa_addr || (ifa->ifa_flags & IFF_LOOPBACK)) {
            continue;
        }
        if (ifa->ifa_addr->sa_family == AF_INET) {
            // IPv4 Address
            tmpAddrPtr=&((sockaddr_in*)ifa->ifa_addr)->sin_addr;
            char addressBuffer[INET_ADDRSTRLEN];
            inet_ntop(AF_INET, tmpAddrPtr, addressBuffer, INET_ADDRSTRLEN);
            ip_address = addressBuffer;
            break;
        }
        else if (ifa->ifa_addr->sa_family == AF_INET6) {
            // IPv6 Address
            tmpAddrPtr=&((sockaddr_in6*)ifa->ifa_addr)->sin6_addr;
            char addressBuffer[INET6_ADDRSTRLEN];
            inet_ntop(AF_INET6, tmpAddrPtr, addressBuffer, INET6_ADDRSTRLEN);
            ip_address = addressBuffer;
            break;
        } 
    }
    if (ifAddrStruct != nullptr) {
        freeifaddrs(ifAddrStruct);
    }
    return ip_address;
}
