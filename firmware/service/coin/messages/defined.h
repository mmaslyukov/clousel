#pragma once
#include <infrastructure/config/entry/char_container.h>
#include <infrastructure/json/named.h>

// BMT  - Broker Message Type
// BMTC - Broker Message Type Commnad
// BMTR - Broker Message Type Response
// BMTE - Broker Message Type Event

#define BMTE_HEARTBEAT "EventHeartbeat"
#define BMTR_ACK "EventAck"
#define BMTR_CONFIG "EventConfig"
#define BMTR_SC_COMPLETED "EventCompleted"

#define BMT_COMMAND "MessageCommand"
#define BMTC_PLAY "Play"
#define BMTC_CONFIG_WRITE "ConfigWrite"
#define BMTC_CONFIG_READ "ConfigRead"


using TypeContainer = infra::CharContainer<20>;
using CommandContainer = infra::CharContainer<20>;
using CarouselIdContainer = infra::CharContainer<37>;
using EventIdContainer = infra::CharContainer<37>;
using ErrorContainer = infra::CharContainer<100>;