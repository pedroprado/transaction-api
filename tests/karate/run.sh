#!/bin/sh

# If pubsub host is not defined use default pubsub
if [ -z "$PUBSUB_HOST" ]
then
  echo "Defining default pubsub host"
  export PUBSUB_HOST=pubsub
fi

java -Dlogback.configurationFile=log-config.xml -cp karate.jar:./cases:KarateUtils.jar com.intuit.karate.Main "$@"