FROM gridworkz/alpine:3.9
ARG RELEASE_DESC
COPY kato-worker /run/kato-worker
COPY entrypoint.sh /run/entrypoint.sh
RUN chmod 655 /run/kato-worker

ENV EX_DOMAIN=ali-sh.gridworkz.com:10080
ENV RELEASE_DESC=${RELEASE_DESC}

ENTRYPOINT ["/run/entrypoint.sh"]
