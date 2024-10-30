# get_filename_component(PARENT_DIR ../ ABSOLUTE)
# "wifi/wifi_mode_softap.cpp" "wifi/wifi mode_station.cpp" "wifi/wifi_mode.cpp" "wifi/wifi.cpp"
idf_component_register(SRCS  "espmain.cpp"
                       INCLUDE_DIRS "." "./../../../")
