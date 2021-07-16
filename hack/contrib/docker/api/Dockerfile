FROM gridworkzapps/alpine:3.4
ARG RELEASE_DESC
LABEL author=gdevs@gridworkz.com
ADD kato-api /run/kato-api
ADD entrypoint.sh /run/entrypoint.sh
WORKDIR /run
ENV RELEASE_DESC=${RELEASE_DESC}
VOLUME [ "/etc/gridworkz" ]
ENTRYPOINT ["/run/entrypoint.sh"]
