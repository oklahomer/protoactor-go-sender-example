FROM alpine:latest

ENV CONSUL_VERSION 1.5.3
ENV HASHICORP_RELEASES=https://releases.hashicorp.com

# Setup
RUN apk add --update ca-certificates
RUN apk add curl

# Install consul
RUN wget https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip
RUN unzip consul_${CONSUL_VERSION}_linux_amd64.zip -d /bin/

EXPOSE 8500

CMD ["/bin/consul", "agent", "-server", "-dev", "-bootstrap",  "-client", "0.0.0.0"]
