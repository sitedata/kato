FROM gridworkz/openresty:1.15.8.2
ARG RELEASE_DESC
RUN apk add --no-cache bash net-tools curl tzdata && \
        cp /usr/share/zoneinfo/America/Toronto /etc/localtime && \
        echo "America/Toronto" >  /etc/timezone && \
        date && apk del --no-cache tzdata
ADD . /run
ENV NGINX_CONFIG_TMPL=/run/nginxtmp
ENV NGINX_CUSTOM_CONFIG=/run/nginx/conf
ENV RELEASE_DESC=${RELEASE_DESC}
ENV OPENRESTY_HOME=/usr/local/openresty
ENV PATH="${PATH}:${OPENRESTY_HOME}/nginx/sbin"
EXPOSE 8080

ENTRYPOINT ["/run/entrypoint.sh"]
