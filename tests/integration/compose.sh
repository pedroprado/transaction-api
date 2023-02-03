docker-compose down

cp -r cases ../karate/cases
cp karate-config.js ../karate
docker build -t karate ../karate/.
rm -rf ../karate/cases
rm ../karate/karate-config.js

docker build -t transaction-api:local ../../.
if [ $? -ne 0 ]; then
    echo "Error building Temis Registration"
    exit 1
fi

# If docker-compose is not installed, download it
docker-compose ps
if [ $? -ne 0 ]; then

  curl -L "https://github.com/docker/compose/releases/download/1.27.3/docker-compose-$(uname -s)-$(uname -m)" -o docker-compose
  if [ $? -ne 0 ]; then
      echo "Error getting docker compose"
      exit 1
  fi
  export PATH=$PATH:$PWD
  chmod +x docker-compose

fi

docker-compose up --force-recreate -d --build