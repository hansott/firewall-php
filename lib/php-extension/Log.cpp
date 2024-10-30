#include "Includes.h"

Log::Log() {
    std::time_t currentTime = std::time(nullptr);
    char timeStr[20];
    std::strftime(timeStr, sizeof(timeStr), "%Y%m%d%H%M%S", std::localtime(&currentTime));
    std::string logFilePath = "/var/log/aikido-" + std::string(PHP_AIKIDO_VERSION) + "/aikido-extension-php-" + timeStr + "-" + std::to_string(getpid()) + ".log";
    this->logFile = fopen(logFilePath.c_str(), "w");
}

Log::~Log() {
    if (!this->logFile) {
        return;
    }

    fclose(this->logFile);
    this->logFile = nullptr;
}

void Log::Write(AIKIDO_LOG_LEVEL level, const char* format, ...) {
    if (!logFile || level < AIKIDO_GLOBAL(log_level)) {
        return;
    }

    fprintf(logFile, "[AIKIDO][%s] ", ToString(level).c_str());

    va_list args;
    va_start(args, format);
    vfprintf(logFile, format, args);
    va_end(args);

    fflush(logFile);
}

std::string Log::ToString(AIKIDO_LOG_LEVEL level) {
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

AIKIDO_LOG_LEVEL Log::ToLevel(std::string level) {
    if (level == "ERROR") {
        return AIKIDO_LOG_LEVEL_ERROR;
    }
    if (level == "WARN") {
        return AIKIDO_LOG_LEVEL_WARN;
    }
    if (level == "INFO") {
        return AIKIDO_LOG_LEVEL_INFO;
    }
    if (level == "DEBUG") {
        return AIKIDO_LOG_LEVEL_DEBUG;
    }
    return AIKIDO_LOG_LEVEL_ERROR;
}
