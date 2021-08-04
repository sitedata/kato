FROM gridworkz/alpine:3.9
ARG RELEASE_DESC
COPY . /run
RUN chmod +x /run/kato-grctl /run/entrypoint.sh
VOLUME [ "/rootfs/root","/rootfs/path","/ssl" ]
ENV RELEASE_DESC=${RELEASE_DESC}
ENTRYPOINT ["/run/entrypoint.sh"]
