#!/bin/bash
# encoding: utf-8
# Name  : test.sh
# Descp : used for 
# Author: jaycee
# Date  : 04/08/18 11:42:30 +0800
__version__=0.1

#set -x                     #print every excution log
set -e                     #exit when error hanppens


#curl -v --data '{"areaid":"test","userid":"test","content":"This is a test messages","timestamp":123,"userdefaddr":"test","expirytime":123,"latitude":123,"longitude":123,"altitude":123}' -X POST http://localhost:8080/messages/?key=abc123
###curl --data '{"title":"test","options":["one","two","three"]}' -X POST http://localhost:8080/polls/?key=abc123
#curl -X GET http://localhost:8080/polls/?key=abc123
curl -X DELETE http://localhost:8080/messages/${1:-""}?key=abc123
curl -X GET http://localhost:8080/messages/?key=abc123
