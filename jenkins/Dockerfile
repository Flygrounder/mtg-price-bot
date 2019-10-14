FROM jenkins/jenkins:latest
USER root
RUN apt update
RUN apt install -y sudo
RUN echo "jenkins ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
USER jenkins
RUN sudo apt update
RUN sudo apt install -y \
	apt-transport-https \
	ca-certificates \
	curl \
	gnupg-agent \
	software-properties-common
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
RUN sudo add-apt-repository \
	"deb [arch=amd64] https://download.docker.com/linux/debian \
	$(lsb_release -cs) \
	stable"
RUN sudo apt update
RUN sudo apt install -y docker-ce docker-ce-cli containerd.io
