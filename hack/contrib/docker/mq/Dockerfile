FROM gridworkz/alpine:3.9
ARG RELEASE_DESC
COPY kato-mq /run/kato-mq
ADD entrypoint.sh /run/entrypoint.sh
RUN chmod 655 /run/kato-mq
EXPOSE 6300

ENV RELEASE_DESC=${RELEASE_DESC}

ENTRYPOINT ["/run/entrypoint.sh"]

