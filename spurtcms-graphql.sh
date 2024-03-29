#!/bin/bash

if [ ! -e /etc/systemd/system/spurtcms-graphql.service ]; then
    yes | cp -i spurtcms-graphql.service /etc/systemd/system/
fi

if [ $1 ]; then
    if [ $1 == "start" ]; then
        echo "Application Started Successfully - Visit http://localhost:8081"
        echo "To Stop Application  - "sudo ./spurtcms-graphql.sh stop" "
        exec systemctl start spurtcms-graphql.service
    fi

    if [ $1 == "stop" ]; then
        echo "-- Application Stopped Successfully --"
        echo "-- To Start Application Again - "sudo ./spurtcms-graphql.sh start" "
        exec systemctl stop spurtcms-graphql.service
    fi
fi
