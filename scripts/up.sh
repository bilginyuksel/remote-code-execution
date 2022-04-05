container_id=$(docker run -dit --privileged -p 8888:8888 rce)
echo "container= $container_id"
docker wait $container_id
docker exec -it $container_id sh -c "docker pull codigician/all-in-one-ubuntu:v1-amd64"
docker exec -it $container_id sh -c "APP_ENV=prod ./rce-engine serve"