#!/bin/sh
mysql --password=950703 -e 'use market'
status="$?"
if [ "$status"="0" ]
then echo "market created already"
else mysql --password=950703 -e 'create database market'
fi
if [ "$status"="0" ]
then echo "database market created"
else echo "failed"
fi
go run /code/test.go host=localhost
