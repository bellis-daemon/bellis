docker-compose -f backend-swarm.yml build
docker-compose -f backend-swarm.yml push
echo "Build and push done!"