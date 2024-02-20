#!/bin/bash

yum install vim docker bind-utils telnet git php -y

sudo amazon-linux-extras install docker

systemctl enable docker
systemctl start docker

cd /tmp

sudo curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) -o /usr/bin/docker-compose
sudo chmod +x /usr/bin/docker-compose

git clone https://github.com/msfidelis/rinhadebackend-go-write-behind.git


sudo amazon-linux-extras enable corretto8
yum install java-17-amazon-corretto.
sudo yum install java-17-amazon-corretto-devel -y

wget https://repo1.maven.org/maven2/io/gatling/highcharts/gatling-charts-highcharts-bundle/3.10.3/gatling-charts-highcharts-bundle-3.10.3-bundle.zip -O gatling.zip
unzip gatling.zip 

sudo mv gatling-* /opt/gatling

echo 'export PATH=$PATH:/opt/gatling/bin' >> ~/.bashrc
echo 'export RESULTS_WORKSPACE=/tmp/gatling/results' >> ~/.bashrc
echo 'export GATLING_HOME=/opt/gatling' >> ~/.bashrc
echo 'export GATLING_BIN_DIR=/opt/gatling/bin' >> ~/.bashrc
echo 'export GATLING_WORKSPACE=/tmp/gatling' >> ~/.bashrc

mkdir -p /opt/gatling/user-files/simulations/rinhadebackend

cp /tmp/rinhadebackend-go-write-behind/misc/gatling/user-files/simulations/rinhadebackend/RinhaBackendCrebitosSimulation.scala /opt/gatling/user-files/simulations/rinhadebackend

mkdir -p /tmp/gatling/results


