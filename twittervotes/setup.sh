#!/bin/bash
# encoding: utf-8
# Name  : setup.sh
# Descp : used for 
# Author: jaycee
# Date  : 04/08/18 11:33:35 +0800
__version__=0.1

#set -x                     #print every excution log
set -e                     #exit when error hanppens#!/bin/bash

export SP_TWITTER_KEY=yC2EDnaNrEhN5fd33g
export SP_TWITTER_SECRET=6n0rToIpskCo1ob
export SP_TWITTER_ACCESSTOKEN=2427-13677
export SP_TWITTER_ACCESSSECRET=SpnZf336u
go build -o twittervotes && ./twittervotes
