pipeline {
	agent any

	stages {
		stage('Pull') {
			steps {
				git pull	
			}
		}
		stage('Test') {
			steps {
				./deploy.sh test
			}
		}
		stage('Deploy') {
			steps {
				./deploy.sh prod
			}
		}
	}
}
