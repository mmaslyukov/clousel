
cmake_minimum_required(VERSION 3.10)

project(strategy)

add_executable(strategy strategy.cpp)

target_include_directories(${PROJECT_NAME} PUBLIC ../../../)



