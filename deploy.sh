#!/bin/bash
TEST_IMAGE_NAME="go_mtg_vk_test_image"
PROD_IMAGE_NAME="go_mtg_vk_prod_image"
if [[ $1 = "test" ]]
then
	TEST_NETWORK_NAME="go_mtg_vk_test_network"
	sudo docker network create $TEST_NETWORK_NAME
	TEST_REDIS_NAME="go_mtg_vk_test_redis"
	sudo docker run -d --rm --name $TEST_REDIS_NAME --network $TEST_NETWORK_NAME --net-alias redis redis
	sudo docker rmi -f $TEST_IMAGE_NAME
	sudo docker build -t $TEST_IMAGE_NAME .
	TEST_CONTAINER_NAME="go_mtg_vk_test_container"
	sudo docker run --rm --name $TEST_CONTAINER_NAME --network $TEST_NETWORK_NAME -e MODE="test" -e VK_TOKEN="$VK_TOKEN" -e VK_SECRET_KEY="$VK_SECRET_KEY" -e VK_GROUP_ID="$VK_GROUP_ID" -e VK_CONFIRMATION_STRING="$VK_CONFIRMATION_STRING" $TEST_IMAGE_NAME
	EXIT_CODE=$(sudo docker inspect $TEST_CONTAINER_NAME --format "{{.State.ExitCode}}")
	sudo docker rm -f $TEST_REDIS_NAME
	exit $(("$EXIT_CODE"))
elif [[ $1 = "prod" ]]
then
	PROD_NETWORK_NAME="go_mtg_vk_prod_network"
	sudo docker network create $PROD_NETWORK_NAME
	PROD_REDIS_NAME="go_mtg_vk_prod_redis"
	sudo docker run -d --net-alias redis --restart always --network $PROD_NETWORK_NAME --name $PROD_REDIS_NAME redis
	sudo docker  rmi -f $PROD_IMAGE_NAME
	sudo docker tag $TEST_IMAGE_NAME $PROD_IMAGE_NAME
	sudo docker rmi $TEST_IMAGE_NAME
	PROD_CONTAINER_NAME="go_mtg_vk_prod_container"
	sudo docker stop $PROD_CONTAINER_NAME
	sudo docker rm $PROD_CONTAINER_NAME
	sudo docker run -v /var/log/go_mtg_vk:/go/src/go-mtg-vk/logs -d -p 80:80 --network $PROD_NETWORK_NAME --restart always --name $PROD_CONTAINER_NAME -e MODE="prod" -e VK_TOKEN="$VK_TOKEN" -e VK_SECRET_KEY="$VK_SECRET_KEY" -e VK_GROUP_ID="$VK_GROUP_ID" -e VK_CONFIRMATION_STRING="$VK_CONFIRMATION_STRING" $PROD_IMAGE_NAME
fi
