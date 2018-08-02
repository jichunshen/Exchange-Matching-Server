#!/bin/bash
host=$1
nc $1 12345 < step1.xml
