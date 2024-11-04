#include "Includes.h"

void PhpLifecycle::ModuleInit() {
    this->mainPID = getpid();
    AIKIDO_LOG_INFO("Main PID is: %u\n", this->mainPID);
    if (!AIKIDO_GLOBAL(agent).Init()) {
        AIKIDO_LOG_INFO("Aikido Agent initialization failed!\n");
    } else {
        AIKIDO_LOG_INFO("Aikido Agent initialization succeeded!\n");
    }
}

void PhpLifecycle::RequestInit() {
    action.Reset();
    requestCache.Reset();
    requestProcessor.RequestInit();
}

void PhpLifecycle::RequestShutdown() {
    requestProcessor.RequestShutdown();
}

void PhpLifecycle::ModuleShutdown() {
    if (this->mainPID == getpid()) {
        AIKIDO_GLOBAL(agent).Uninit();
    } else {
        requestProcessor.Uninit();
    }
}

PhpLifecycle phpLifecycle;