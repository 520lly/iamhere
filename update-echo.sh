#!/bin/sh

sudo git pull
sudo service iamhere-echo stop
sudo service iamhere-echo start
echo "Done!"

