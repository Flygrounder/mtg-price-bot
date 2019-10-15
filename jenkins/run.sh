#!/bin/bash
IMAGE_NAME="mtg_vk_go_jenkins_image"
CONTAINER_NAME="mtg_vk_go_jenkins_container"
sudo docker build -t $IMAGE_NAME .
sudo docker run -d \
	--restart always\
	-p 8081:8080\
	--privileged\
	-v /var/run/docker.sock:/var/run/docker.sock\
	-v jenkins_home_mtg_go:/var/jenkins_home\
	--name $CONTAINER_NAME\
	-e VK_TOKEN=$VK_TOKEN\
	-e VK_SECRET_KEY=$VK_SECRET_KEY\
	-e VK_GROUP_ID=$VK_GROUP_ID\
	-e VK_CONFIRMATION_STRING=$VK_CONFIRMATION_STRING\
	--name $CONTAINER_NAME\
   	$IMAGE_NAME
sudo docker exec -it $CONTAINER_NAME sudo service docker start 
