pipeline {
	agent any
	stages {
		stage('Pull') {
			steps {
				sh 'git pull'
			}
		}
		stage('Test') {
			steps {
				sh './deploy.sh test'
			}
		}
		stage('Deploy') {
			steps {
				sh './deploy.sh prod'
			}
		}
	}
}
