FROM openjdk:alpine
WORKDIR /karate
RUN mkdir /cases
ADD . /karate

RUN chmod +x /karate/run.sh

# Current directory and ./cases folder are in class path. Karate config
ENTRYPOINT ["/karate/run.sh"]