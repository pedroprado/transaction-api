if [ -z "$1" ]
then
  echo "Should inform Karate Docker container ID as first parameter"
  echo "Example: sh getLog.sh <CONTAINER_ID>"
  exit 1
fi

docker cp $1:/karate/target/karate.log ./karate.log
if [ $? -ne 0 ]; then
  echo "Error getting log from container"
else
  echo "Log successfully obtained. Could be found with name karate.log in the current folder"
fi