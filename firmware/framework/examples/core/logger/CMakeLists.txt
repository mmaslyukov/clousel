
cmake_minimum_required(VERSION 3.10)

project(logger)

add_executable(logger logger.cpp)

target_include_directories(${PROJECT_NAME} PUBLIC ../../../../)



