#pragma once
#include <infrastructure/config/entry/char_container.h>
#include <infrastructure/json/named.h>

#define BMT_HEARTBEAT "EventHeartbeat"
#define BMT_COMMAND "MessageCommand"
#define BMT_ACK "EventAck"

#define BMTC_PLAY "Play"

using TypeContainer = infra::CharContainer<20>;
using CommandContainer = infra::CharContainer<20>;
using CarouselIdContainer = infra::CharContainer<37>;
using EventIdContainer = infra::CharContainer<37>;
using ErrorContainer = infra::CharContainer<100>;