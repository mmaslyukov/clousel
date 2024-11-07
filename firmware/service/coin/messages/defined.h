#pragma once
#include <infrastructure/config/entry/char_container.h>
#include <infrastructure/json/named.h>

// BMF - Broker Message Field Commnad(Request)
// BMTC - Broker Message Type Commnad(Request)
// BMTR - Broker Message Type Response
// BMTE - Broker Message Type Event

#define BMTE_HEARTBEAT "Evt.Heartbeat"
#define BMTE_COMPLETED "Evt.Completed"

#define BMTC_PLAY "Req.Play"
#define BMTC_CONFIG_WRITE "Req.Config.Write"
#define BMTC_CONFIG_READ "Req.Config.Read"

#define BMTR_ACK "Res.Ack"
#define BMTR_CONFIG "Res.Config"

#define BMF_TYPE "Type"
#define BMF_EVENT_ID "EvtId"
#define BMF_CAROUSEL_ID "CarId"
#define BMF_CORRELATION_ID "CorId"
#define BMF_SEQUENCE_NUM "SeqNum"
#define BMF_ERROR "Error"
#define BMF_CONFIG "Config"

using TypeContainer = infra::CharContainer<20>;
using CarouselIdContainer = infra::CharContainer<37>;
using EventIdContainer = infra::CharContainer<37>;
using ErrorContainer = infra::CharContainer<100>;