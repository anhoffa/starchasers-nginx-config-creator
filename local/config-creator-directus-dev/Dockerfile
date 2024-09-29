FROM golang:1.22

RUN apt update && apt install nginx libnginx-mod-stream openssh-server -y

RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN echo 'root:superSecret123' | chpasswd

RUN rm /etc/nginx/nginx.conf
ADD start.sh start.sh
RUN chmod +x ./start.sh
ENTRYPOINT ./start.sh
