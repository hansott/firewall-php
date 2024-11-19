#include "Includes.h"

std::unordered_map<std::string, SinkStats> stats;

ScopedTimer::ScopedTimer() : start(std::chrono::high_resolution_clock::now()) {}

ScopedTimer::ScopedTimer(std::string key) : key(key), start(std::chrono::high_resolution_clock::now()) {}

void ScopedTimer::SetSink(std::string key) {
    this->key = key;
}

ScopedTimer::~ScopedTimer() {
    if (key.empty()) {
        return;
    }
    int64_t duration = std::chrono::duration_cast<std::chrono::nanoseconds>(std::chrono::high_resolution_clock::now() - start).count();
    SinkStats& sinkStats = stats[key];
    sinkStats.timings.push_back(duration);
}

void SinkStats::IncrementAttacksDetected() {
    attacksDetected++;
}

void SinkStats::IncrementAttacksBlocked() {
    attacksBlocked++;
}

void SinkStats::IncrementInterceptorThrewError() {
    interceptorThrewError += 1;
}

void SinkStats::IncrementWithoutContext() {
    withoutContext += 1;
}