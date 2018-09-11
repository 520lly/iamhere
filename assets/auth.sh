#!/bin/bash
# encoding: utf-8
# Name  : auth.sh
# Descp : used for 
# Author: jaycee
# Date  : 11/09/18 11:58:56 +0800
__version__=0.1

#set -x                     #print every excution log
set -e                     #exit when error hanppens

sudo cp 214987401110045.* /etc/ssl
sudo openssl genrsa -out server.key 2048
sudo openssl req -new -key server.key -out server.csr
sudo openssl x509 -req -in server.csr -CA /etc/ssl/214987401110045.pem -CAkey /etc/ssl/214987401110045.key -CAcreateserial -out server.crt -days 500
sudo mv server.* /etc/ssl/iamhere
