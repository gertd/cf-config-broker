FROM alpine

ADD cf-config-broker /
ADD cf-config-broker.json /

CMD ["/cf-config-broker"]

EXPOSE 3000