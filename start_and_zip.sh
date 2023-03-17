#!/usr/bin/env bash


COLOR_SUFFIX="\033[0m"
RED_PREFIX="\033[31m"
GREEN_PREFIX="\033[32m"
YELLOW_PREFIX="\033[33m"

bin_dir="./bin"
#Automatically created when there is no bin, logs folder
if [ ! -d $bin_dir ]; then
  mkdir -p $bin_dir
fi
./build_all_service.sh

  if [ $? -ne 0 ]; then
        exit -1
        else
    cd   bin_dir
   ./runapi_generator
   cd ..
   tar -czvf OpenIM-Server.zip OpenIM服务器API/
  fi