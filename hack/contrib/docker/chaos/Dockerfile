FROM gridworkz/alpine:3.9
ARG RELEASE_DESC
LABEL author="gdevs@gridworkz.com"

ENV WORK_DIR=/run

RUN apk --no-cache add openssl openssh-client subversion
COPY kato-chaos entrypoint.sh /run/
COPY export-app /src/export-app

WORKDIR $WORK_DIR

ENV RELEASE_DESC=${RELEASE_DESC}

ENTRYPOINT ["/run/entrypoint.sh"]
