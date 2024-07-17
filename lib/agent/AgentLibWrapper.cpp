#include "libaikido_agent.h"

extern "C" {

GoUint8 WrapAgentInit(GoString initJson) {
    return AgentInit(initJson);
}

// Function to uninitialize the agent
void WrapUninitAgent() {
    AgentUninit();
}

} // extern "C"