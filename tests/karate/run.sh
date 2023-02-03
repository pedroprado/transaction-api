#!/bin/sh

# If pubsub host is not defined use default pubsub
if [ -z "$PUBSUB_HOST" ]
then
  echo "Defining default pubsub host"
  export PUBSUB_HOST=pubsub
fi


/karate/wait-for -t 30 ${PUBSUB_HOST}:8682 -- java -Dlogback.configurationFile=log-config.xml -cp karate.jar:./cases:KarateUtils.jar com.intuit.karate.Main "$@"