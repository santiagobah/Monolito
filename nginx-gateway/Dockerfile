FROM nginx:alpine
RUN apk add --no-cache curl bash unzip

ENV CONSUL_TEMPLATE_VERSION=0.39.0
RUN curl -L https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip -o consul-template.zip && \
    unzip consul-template.zip -d /usr/local/bin && \
    rm consul-template.zip

COPY templates/nginx.conf.ctmpl /etc/nginx/templates/nginx.conf.ctmpl
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
