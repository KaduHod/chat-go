FROM golang:latest

RUN mkdir -p /var/log/chat
RUN mkdir -p /var/log/chat
RUN touch /var/log/chat/request.log
RUN touch /var/log/chat/db.log
RUN touch /var/log/chat/errors.log
RUN touch /var/log/chat/tests.log
RUN touch /var/log/chat/debug.log

EXPOSE 3000
