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
