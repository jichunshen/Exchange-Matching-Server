#!/bin/bash
int=10000
num=$1
up=$[10000+$num]
size=$[93+39*$num]
amount=100
echo "$size" > single.xml
echo -e "<?xml version="1.0" encoding="UTF-8"?>" >> single.xml
echo -e "<create>" >> single.xml
while(($int<=$up))
do
    echo -e "  <account id=\"$int\" balance=\"1000\"/>" >> single.xml
    let "int++"
done
sleep 2
echo -e "</create>" >> single.xml
#nc vcm-3855.vm.duke.edu 12345 < single.xml
echo "done"

