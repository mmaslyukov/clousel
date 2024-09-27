
cmake_minimum_required(VERSION 3.10)

project(broker)

find_package(eclipse-paho-mqtt-c CONFIG REQUIRED)

add_executable(broker broker.cpp)

target_include_directories(${PROJECT_NAME} PUBLIC ../../../)

target_link_libraries(${PROJECT_NAME} PRIVATE eclipse-paho-mqtt-c::paho-mqtt3a eclipse-paho-mqtt-c::paho-mqtt3c eclipse-paho-mqtt-c::paho-mqtt3as eclipse-paho-mqtt-c::paho-mqtt3cs)
