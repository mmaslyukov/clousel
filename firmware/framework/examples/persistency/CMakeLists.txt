
cmake_minimum_required(VERSION 3.10)

project(persistency)

add_executable(persistency persistency.cpp)

target_include_directories(${PROJECT_NAME} PUBLIC ../../../)



