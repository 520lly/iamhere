#!/bin/bash
# encoding: utf-8
# Name  : test.sh
# Descp : used for 
# Author: jaycee
# Date  : 04/08/18 11:42:30 +0800
__version__=0.1

#set -x                     #print every excution log
set -e                     #exit when error hanppens

#IP=35.187.154.122
IP=localhost

#curl --data '{"areaid":"test","userid":"test","content":"This is a test messages","timestamp":123,"userdefaddr":"test","expirytime":123,"latitude":123,"longitude":123,"altitude":123}' -X POST http://$IP:8080/messages/?key=abc123
###curl --data '{"title":"test","options":["one","two","three"]}' -X POST http://$IP:8080/polls/?key=abc123
#curl -X DELETE http://$IP:8080/messages/${1:-""}?key=abc123
#curl -X GET http://$IP:8080/messages/?key=abc123
#curl --data '{"name":"雪窦山 徐凫岩瀑布","description":"喜欢瀑布下沐浴水汽的感觉","address1":"浙江省宁波市奉化区","address2":"","category":1,"type":0,"latitude":29.7039399637,"longitude":121.1754884604,"altitude":0,"radius":50.00}' -X POST http://$IP:8080/areas/?key=abc123
#curl -X GET http://$IP:8080/areas/?key=abc123
curl -X GET "http://$IP:8080/areas/?key=abc123&&type=0&&category=1&&latitude=10.12212"
#curl -X GET "http://$IP:8080/areas/5b743e9181b37309c1ef3e80/?key=abc123"

