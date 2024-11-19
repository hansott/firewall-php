#pragma once

class ScopedTimer {
   private:
    std::string key;
    std::chrono::high_resolution_clock::time_point start;
    
   public:
    ScopedTimer();
    ScopedTimer(std::string key);
    void SetSink(std::string key);
    ~ScopedTimer();
};

class SinkStats {
    public:
     uint64_t attacksDetected = 0;
     uint64_t attacksBlocked = 0;
     uint64_t interceptorThrewError = 0;
     uint64_t withoutContext = 0;
     std::vector<int64_t> timings;

    void IncrementAttacksDetected();
    void IncrementAttacksBlocked();
    void IncrementInterceptorThrewError();
    void IncrementWithoutContext();
};

extern std::unordered_map<std::string, SinkStats> stats;
