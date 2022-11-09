#!/bin/bash

#keep LAN alive
echo "................................................"
date "+%Y.%m.%d %H:%M:%S"
echo " "
while :
do
    sleep 300s
    ping -c2 google.com


    if [ $? != 0 ]
    then
        echo " "
        echo "No network connection, restarting wlan0"
        sudo reboot

    #end of if statement
    fi
done
echo "................................................"
echo " "