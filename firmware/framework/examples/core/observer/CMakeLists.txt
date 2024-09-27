
cmake_minimum_required(VERSION 3.10)

project(observer)

add_executable(observer observer.cpp)

target_include_directories(${PROJECT_NAME} PUBLIC ../../../../)



