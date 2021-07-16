FROM gridworkzapps/alpine:3.4
ARG RELEASE_DESC
COPY . /run/
RUN chmod 655 /run/kato-init-probe /run/entrypoint.sh
ENV RELEASE_DESC=${RELEASE_DESC}
ENTRYPOINT [ "/run/entrypoint.sh" ]
CMD ["probe"]

