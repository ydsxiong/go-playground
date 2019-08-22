FROM alpine:3.4

RUN apk -U add ca-certificates

EXPOSE 8081

ADD gorestapiclient /bin/gorestapiclient
ADD k8s/config.yml.dist /etc/gorestapiclient/config.yml

CMD ["gorestapiclient", "-config", "/etc/gorestapiclient/config.yml"]
