docker-compose down
if [ "$1" = "dev" ]; then
    echo "deploying dev..."
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
else
    echo "deploying..."
    docker-compose -f docker-compose.yml -f docker-compose.deploy.yml up -d
fi
sleep 10
docker ps -a
exit 0
