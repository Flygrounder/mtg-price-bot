pipeline {
	agent any
	stages {
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
