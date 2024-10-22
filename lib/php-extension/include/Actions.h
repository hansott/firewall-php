#pragma once

#include "Includes.h"

class Action {
    private:
        bool exit = false;
        bool block = false;
        std::string type;
        std::string trigger;
        std::string ip;

    private:
        ACTION_STATUS executeThrow(json &event);

        ACTION_STATUS executeExit(json &event);

        ACTION_STATUS executeStore(json &event);

    public:
        Action() = default;
        ~Action() = default;

        ACTION_STATUS Execute(std::string &event);
        void Reset();

        bool Exit();
        bool Block();
        char* Type();
        char* Trigger();
        char* Ip();
};

extern Action action;
