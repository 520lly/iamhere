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

#curl --data '{"areaid":"test","userid":"test","content":"This is a test messages","userdefaddr":"test","expirytime":123,"latitude":29.7039399637,"longitude":121.1754884604,"altitude":123}' -X POST "http://$IP:8080/messages/"
#curl --data '{"areaid":"test","userid":"test","content":"This is a test messages","userdefaddr":"test","expirytime":123,"latitude":29.7039399637,"longitude":121.1754884604,"altitude":123}' -X GET "http://$IP:8080/messages/"
curl -X GET "http://$IP:8080/messages/?debug=121.1754884604"
#curl -X DELETE "http://$IP:8080/areas/${1:-""}"
#curl --data '{"areaid":"test","userid":"test","content":"This is a test messages","userdefaddr":"test","expirytime":123,"latitude":9.7039399637,"longitude":121.1754884604,"altitude":123}' -X GET "http://$IP:8080/areas/"
#curl -X GET "http://$IP:8080/messages/?longitude=121.1754884604"
#curl -X DELETE "http://$IP:8080/messages/${1:-""}"
###curl --data '{"title":"test","options":["one","two","three"]}' -X POST http://$IP:8080/polls/?key=abc123
#curl -X GET http://$IP:8080/messages/?key=abc123
curl -X DELETE http://$IP:8080/accounts/${1:-""} | jq
curl -X GET http://$IP:8080/accounts/?debug=1 | jq
curl --data '{"associatedId":"wechat_djakdjakdj", "nickname":"mnma","password":"12345ddqw", "email":"jacking.wang.wjq@163.com","firstname":"jianqing","lastname":"wang","phonenumber":"13167016114","birthday":"19990919","gender":"male"}' -X POST "http://$IP:8080/accounts/register" | jq
#curl --data '{"password":"12345ddqw", "email":"jacking.wang.wjq@gmail.com","firstname":"jianqing","lastname":"wang","phonenumber":"13167016112","birthday":"19990919","gender":"male"}' -X POST "http://$IP:8080/accounts/register"
#curl --data '{"password":"112345ddqw", "email":"jacking.wang.wjq@gmail.com","firstname":"jianqing","lastname":"wang","phonenumber":"13167016112","birthday":"19990919","gender":"female"}' -X PUT "http://$IP:8080/accounts/${1:-""}"
#curl --data '{"password":"112345ddqw", "username":"jacking.wang.wjq@gmail.com"}' -X GET "http://$IP:8080/accounts/login"
#curl -X DELETE "http://$IP:8080/accounts/${1:-""}"
#curl -X GET "http://$IP:8080/accounts/?debug=1"
#curl --data '{"name":"jianqing雪窦山 徐凫岩瀑布","description":"喜欢瀑布下沐浴水汽的感觉","province":"Zhengjiang","city":"Ningbo","district":"Fenghua","address1":"浙江省宁波市奉化区","address2":"","category":1,"type":0,"latitude":2.7039399637,"longitude":121.1754884604,"altitude":0,"radius":50.00}' -X POST http://$IP:8080/areas/${1:-""}
#curl -H '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6MTUzNjI4MTUxNSwibmFtZSI6IiJ9.XrtLNodTmGiWfOcOdx1t_Soavdi51UUt--RUMzkUxZg"}' -X GET "http://$IP:8080/areas/?debug=1"
#curl --data '{"name":"dad雪窦山 徐凫岩瀑布","description":"喜欢瀑布下沐浴水汽的感觉","address1":"浙江省宁波市奉化区","address2":"","category":1,"type":0,"latitude":29.7039399637,"longitude":121.1754884604,"altitude":0,"radius":50.00}' -X POST http://$IP:8080/areas/?key=abc123
#curl --data '{"name":"2323雪窦山 徐凫岩瀑布","description":"喜欢瀑布下沐浴水汽的感觉","address1":"浙江省宁波市奉化区","address2":"","category":1,"type":0,"latitude":29.7039399637,"longitude":121.1754884604,"altitude":0,"radius":50.00}' -X POST http://$IP:8080/areas/?key=abc123
#curl --data '{"name":"11雪窦山 徐凫岩瀑布","description":"喜欢瀑布下沐浴水汽的感觉","address1":"浙江省宁波市奉化区","address2":"","category":1,"type":0,"latitude":29.7039399637,"longitude":121.1754884604,"altitude":0,"radius":50.00}' -X POST http://$IP:8080/areas/?key=abc123
#curl --data '{"name":"9雪窦山 徐凫岩瀑布","description":"喜欢瀑布下沐浴水汽的感觉","address1":"浙江省宁波市奉化区","address2":"","category":1,"type":0,"latitude":29.7039399637,"longitude":121.1754884604,"altitude":0,"radius":50.00}' -X POST http://$IP:8080/areas/?key=abc123 | jq

#curl --data '{"nickname":"Mhajd","associatedId":"wechat_djakdjakdj","password":"testtttt","email":"jackaing.wang.wjq@gmail.com","firstname":"jianqing","lastname":"wang","phonenumber":"13167016112","birthday":"19990919","gender":"male"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
#curl -X GET "http://$IP:8080/accounts/?key=abc123&debug=1"
#curl --data '{"nickname":"hahahah","password":"testtttt","email":"jacking.wang.wjq@gmail.com","firstname":"jianqing","lastname":"wang","phonenumber":"13167016112","birthday":"19990919","gender":"male"}' -X POST "http://$IP:8080/accounts/5b7a4a9ec2217bf5e4c3fd2a/?key=abc123"
#curl -X GET "http://$IP:8079/areas/?key=abc123&debug=1"
exit
curl --data '{"nickname":"hajd"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"associatedId":"wechat_djjakdj"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"password":"ddsdsesttttt"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"email":"jackaing.wang.wjdadq@gmail.com"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"firstname":"hdkajdka"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"lastname":"jliang"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"phonenumber":"13167722399"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"birthday":"19910919"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
curl --data '{"gender":"female"}' -X POST "http://$IP:8080/accounts/${1:-""}/?key=abc123"
#curl -X DELETE "http://$IP:8080/accounts/${1:-""}/?key=abc123" | jq
exit

DATA=`curl -X GET "http://$IP:8080/areas/?key=abc123&&debug=1" | jq ".data[0]"`
echo $DATA
#data=`echo $DATA | sed 's/"/\\\"/g'`
curl --data "$DATA" -X POST "http://$IP:8080/areas/?key=abc123"


if [ ${1:-"GET"} = "DELETE" -a ${2:-"all"} ]
then
    RET=`curl -X GET "http://$IP:8080/areas/?key=abc123&&debug=1"` 
    echo $RET > file
    index=0
    COUNT=`echo $RET | jq -r ".count"`
    echo $COUNT
    while [[ $index -lt $COUNT ]]
    do
        id=`jq -r ".data[$index].id" ./file`
        curl -X DELETE "http://$IP:8080/areas/$id?key=abc123" | jq
        index=$(($index+1))
    done
fi

#curl -X GET "http://$IP:8080/areas/?key=abc123&&type=0&&category=1&&latitude=10.12212"
#curl -X GET "http://$IP:8080/areas/5b743e9181b37309c1ef3e80/?key=abc123"
#curl -X DELETE http://$IP:8080/areas/${1:-""}?key=abc123 | jq

