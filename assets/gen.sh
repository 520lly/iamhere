#!/bin/bash
# encoding: utf-8
# Name  : gen.sh
# Descp : used for 
# Author: jaycee
# Date  : 11/08/18 17:25:56 +0800
__version__=0.1

#set -x                     #print every excution log
set -e                     #exit when error hanppens

openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=iamhere.com" -days 5000 -out ca.crt
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=localhost" -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000
