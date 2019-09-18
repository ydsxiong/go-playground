FROM alpine:3.4

RUN apk -U add ca-certificates

EXPOSE 8080

ADD boxoffice /bin/boxoffice
ADD config.yml.dist /etc/boxoffice/config.yml

CMD ["boxoffice", "-config", "/etc/boxoffice/config.yml"]
