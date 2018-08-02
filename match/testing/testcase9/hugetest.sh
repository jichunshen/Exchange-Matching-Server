#!/bin/bash
int=10
up=99
balance=1000
symbol="SYM"
amount=100
while(($int<=$up))
do
    echo "164" > create$int.xml
    echo -e "<?xml version="1.0" encoding="UTF-8"?>" >> create$int.xml
    echo -e "<create>" >> create$int.xml
    echo -e "  <account id=\"$int\" balance=\"$balance\"/>" >> create$int.xml
    echo -e "  <symbol sym=\"$symbol$int\">" >> create$int.xml
    echo -e "    <account id=\"$int\">1000</account>" >> create$int.xml
    echo -e "  </symbol>" >> create$int.xml
    echo -e "</create>" >> create$int.xml
    let "int++"
done
int=10
up=99
while(($int<=$up))
do
    nc vcm-3855.vm.duke.edu 12345 < create$int.xml 
    let int++
done
rm create*.xml
sleep 5
echo "create done"
int=11
up=61
while(($int<=$up))
do
    echo "144" > transaction$int.xml
    echo -e "<?xml version="1.0" encoding="UTF-8"?>" >> transaction$int.xml
    echo -e "<transactions account=\"$int\">" >> transaction$int.xml
    echo -e "  <order sym=\"SYM10\" amount=\"10\" limit=\"10\"/>" >> transaction$int.xml
    echo -e "  <query id=\"$int\"/>" >> transaction$int.xml
    echo -e "</transactions>" >> transaction$int.xml
    let "int++"
done
int=11
up=45
while(($int<=$up))
do
    nc vcm-3855.vm.duke.edu 12345 < transaction$int.xml &
    let "int++"
done
nc vcm-3855.vm.duke.edu 12345 < test.xml
int=46
up=61
while(($int<=$up))
do
    nc vcm-3855.vm.duke.edu 12345 < transaction$int.xml &
    let "int++"
done
rm transaction*.xml 
sleep 5
echo "transaction done"