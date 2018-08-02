#!/bin/bash
host=$1
nc $host 12345 < step1.xml
nc $host 12345 < step2.xml
nc $host 12345 < step3.xml
nc $host 12345 < step4.xml
