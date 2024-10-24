#pragma once

enum AIKIDO_LOG_LEVEL {
    AIKIDO_LOG_LEVEL_DEBUG,
    AIKIDO_LOG_LEVEL_INFO,
    AIKIDO_LOG_LEVEL_WARN,
    AIKIDO_LOG_LEVEL_ERROR
};

#if defined(ZEND_DEBUG)
#define AIKIDO_LOG_DEBUG(format, ...) logger.Write(AIKIDO_LOG_LEVEL_DEBUG, format, ##__VA_ARGS__)
#else
/* Disable debugging logs for production builds */
#define AIKIDO_LOG_DEBUG(format, ...)
#endif

#define AIKIDO_LOG_INFO(format, ...) logger.Write(AIKIDO_LOG_LEVEL_INFO, format, ##__VA_ARGS__)
#define AIKIDO_LOG_WARN(format, ...) logger.Write(AIKIDO_LOG_LEVEL_WARN, format, ##__VA_ARGS__)
#define AIKIDO_LOG_ERROR(format, ...) logger.Write(AIKIDO_LOG_LEVEL_ERROR, format, ##__VA_ARGS__)

class Log {
   private:
    FILE* logFile = nullptr;

   public:
    Log();
    ~Log();

    void Write(AIKIDO_LOG_LEVEL level, const char* format, ...);

    static std::string ToString(AIKIDO_LOG_LEVEL level);

    static AIKIDO_LOG_LEVEL ToLevel(std::string level);
};

extern Log logger;