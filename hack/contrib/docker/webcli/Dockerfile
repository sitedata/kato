FROM gridworkz/alpine:3.9
ARG RELEASE_DESC
ENV PORT 7070

ADD kato-webcli /usr/bin/kato-webcli
ADD entrypoint.sh /entrypoint.sh
RUN mkdir /root/.kube

ENV RELEASE_DESC=${RELEASE_DESC}
ENTRYPOINT ["/entrypoint.sh"]
