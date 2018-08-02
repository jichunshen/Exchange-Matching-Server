#!/bin/bash
int=0
num=$1
file=$2
while(($int<$num))
do
    nc vcm-3855.vm.duke.edu 12345 < $file
    let "int++"
done
echo "done"
