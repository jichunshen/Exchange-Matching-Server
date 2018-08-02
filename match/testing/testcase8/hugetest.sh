#!/bin/bash
int=11
up=50
balance=1000
symbol="SYM"
amount=100
while(($int<=$up))
do
    echo "169" > create$int.xml
    echo -e "<?xml version="1.0" encoding="UTF-8"?>" >> create$int.xml
    echo -e "<create>" >> create$int.xml
    echo -e "  <account id=\"$int\" balance=\"$balance\"/>" >> create$int.xml
    echo -e "  <symbol sym=\"SYM99\">" >> create$int.xml
    echo -e "    <account id=\"$int\">100</account>" >> create$int.xml
    echo -e "  </symbol>" >> create$int.xml
    echo -e "</create>" >> create$int.xml
    let "int++"
done
int=51
up=90
while(($int<=$up))
do
    echo "169" > create$int.xml
    echo -e "<?xml version="1.0" encoding="UTF-8"?>" >> create$int.xml
    echo -e "<create>" >> create$int.xml
    echo -e "  <account id=\"$int\" balance=\"$balance\"/>" >> create$int.xml
    echo -e "  <symbol sym=\"SYM10\">" >> create$int.xml
    echo -e "    <account id=\"$int\">100</account>" >> create$int.xml
    echo -e "  </symbol>" >> create$int.xml
    echo -e "</create>" >> create$int.xml
    let "int++"
done
int=11
up=90
while(($int<=$up))
do
    nc vcm-3855.vm.duke.edu 12345 < create$int.xml 
    let int++
done
rm create*.xml
sleep 5
echo "create done"
int=11
up=50
while(($int<=$up))
do
    echo "180" > transaction$int.xml
    echo -e "<?xml version="1.0" encoding="UTF-8"?>" >> transaction$int.xml
    echo -e "<transactions account=\"$int\">" >> transaction$int.xml
    echo -e "  <order sym=\"SYM10\" amount=\"10\" limit=\"10\"/>" >> transaction$int.xml
    echo -e "  <order sym=\"SYM99\" amount=\"-10\" limit=\"20\"/>" >> transaction$int.xml
    echo -e "</transactions>" >> transaction$int.xml
    let "int++"
done
int=51
up=90
while(($int<=$up))
do
    echo "180" > transaction$int.xml
    echo -e "<?xml version="1.0" encoding="UTF-8"?>" >> transaction$int.xml
    echo -e "<transactions account=\"$int\">" >> transaction$int.xml
    echo -e "  <order sym=\"SYM99\" amount=\"10\" limit=\"20\"/>" >> transaction$int.xml
    echo -e "  <order sym=\"SYM10\" amount=\"-10\" limit=\"10\"/>" >> transaction$int.xml
    echo -e "</transactions>" >> transaction$int.xml
    let "int++"
done
int=11
jnt=51
up=20
while(($int<=$up))
do
    nc vcm-3855.vm.duke.edu 12345 < transaction$int.xml &
    nc vcm-3855.vm.duke.edu 12345 < transaction$jnt.xml &
    let "int++"
    let "jnt++"
done
rm transaction*.xml
sleep 5
echo "transaction done"
