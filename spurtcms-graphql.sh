#!/bin/bash

if [ ! -e /etc/systemd/system/spurtcms-graphql.service ]; then
    yes | cp -i spurtcms-graphql.service /etc/systemd/system/
fi

if [ $1 ]; then
    if [ $1 = "start" ]; then
        echo ""
        echo "Application Started Successfully - Visit http://localhost:8081"
        echo "To Stop Application : 'sudo sh ./spurtcms-graphql.sh stop' "
        echo ""
        exec systemctl start spurtcms-graphql.service
    fi

    if [ $1 = "stop" ]; then
        echo ""
        echo "-- Application Stopped Successfully --"
        echo "-- To Start Application Again : 'sudo sh ./spurtcms-graphql.sh start' "
        echo ""
        exec systemctl stop spurtcms-graphql.service
    fi
fi
