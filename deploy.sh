#!/bin/bash
if [[ $1 = "test" ]]
then
	TEST_NETWORK_NAME="go_mtg_vk_test_network"
	sudo docker network create $TEST_NETWORK_NAME || true
	TEST_REDIS_NAME="go_mtg_vk_test_redis"
	sudo docker run -d --rm --name $TEST_REDIS_NAME --network $TEST_NETWORK_NAME --net-alias redis redis
	TEST_IMAGE_NAME="go_mtg_vk_test_image"
	sudo docker build -t $TEST_IMAGE_NAME .
	TEST_CONTAINER_NAME="go_mtg_vk_test_container"
	sudo docker run --rm --name $TEST_CONTAINER_NAME --network $TEST_NETWORK_NAME -e MODE="test" -e VK_TOKEN=$VK_TOKEN -e VK_SECRET_KEY=$VK_SECRET_KEY -e VK_GROUP_ID=$VK_GROUP_ID -e VK_CONFIRMATION_STRING=$VK_CONFIRMATION_STRING $TEST_IMAGE_NAME
	EXIT_CODE=$(sudo docker inspect $TEST_CONTAINER_NAME --format "{{.State.ExitCode}}")
	sudo docker rm -f $TEST_REDIS_NAME
	sudo docker image rm -f $TEST_IMAGE_NAME
	exit $EXIT_CODE
elif [[ $1 = "prod" ]]
then
	PROD_NETWORK_NAME="go_mtg_vk_prod_network"
	sudo docker network create $PROD_NETWORK_NAME || true
	PROD_REDIS_NAME="go_mtg_vk_prod_redis"
	sudo docker run -d --restart always --network $PROD_NETWORK_NAME --name $PROD_REDIS_NAME redis
	PROD_IMAGE_NAME="go_mtg_vk_prod_image"
	sudo docker build -t $PROD_IMAGE_NAME .
	PROD_CONTAINER_NAME="go_mtg_vk_prod_container"
	sudo docker run --network $PROD_NETWORK_NAME --restart always --name $PROD_CONTAINER_NAME -e MODE="prod" -e VK_TOKEN=$VK_TOKEN -e VK_SECRET_KEY=$VK_SECRET_KEY -e VK_GROUP_ID=$VK_GROUP_ID -e VK_CONFIRMATION_STRING=$VK_CONFIRMATION_STRING $PROD_IMAGE_NAME
fi
