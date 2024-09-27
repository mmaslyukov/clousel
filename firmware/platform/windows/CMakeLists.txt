
cmake_minimum_required(VERSION 3.10)

project(clousel)
# set(CMAKE_VERBOSE_MAKEFILE ON)
set(CMAKE_BUILD_TYPE Debug)
# set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -g3 -ggdb")

find_package(eclipse-paho-mqtt-c CONFIG REQUIRED)

add_executable(clousel winmain.cpp)

target_include_directories(${PROJECT_NAME} PUBLIC ../../)

target_link_libraries(${PROJECT_NAME} PRIVATE eclipse-paho-mqtt-c::paho-mqtt3a eclipse-paho-mqtt-c::paho-mqtt3c eclipse-paho-mqtt-c::paho-mqtt3as eclipse-paho-mqtt-c::paho-mqtt3cs)
