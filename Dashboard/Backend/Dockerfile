# syntax=docker/dockerfile:1
FROM maven:3.9-eclipse-temurin-21 AS build
WORKDIR /app

# Arguments Maven en plus, ex. "-s .mvn/settings-cn.xml" pour passer par le miroir.
ARG MAVEN_ARGS=""

COPY pom.xml ./
COPY .mvn .mvn
COPY src src
# Le cache ~/.m2 survit aux builds : les dépendances ne sont téléchargées qu'une fois.
RUN --mount=type=cache,target=/root/.m2 \
    mvn -B -ntp ${MAVEN_ARGS} package -DskipTests \
    && cp target/backend-*.jar app.jar

FROM eclipse-temurin:21-jre
WORKDIR /app
RUN useradd --system --uid 10001 autostrike
COPY --from=build /app/app.jar app.jar
USER autostrike
EXPOSE 3001
ENTRYPOINT ["java", "-jar", "app.jar"]
