FROM nginx:latest

RUN apt-get update -y && apt-get install -y software-properties-common
RUN apt-get install -y vim nano
RUN apt-get install -y net-tools
RUN apt-get install -y iputils-ping
RUN apt-get install wget -y
RUN apt-get install curl -y
# Nginx will listen on this port
EXPOSE 80
