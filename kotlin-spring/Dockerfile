ARG JAVA_VERSION=21
ARG SPRING_ENV="prod"

FROM openjdk:${JAVA_VERSION}-slim-bullseye as app

RUN apt-get update && apt-get -y upgrade

RUN mkdir /app
WORKDIR /app
ENV ENVIRONMENT="${SPRING_ENV}"

USER 1000

COPY --chown=1000:1000 build/libs/app.jar ./

EXPOSE 8080
ENTRYPOINT ["java"]

CMD ["-jar", "/app/app.jar"]