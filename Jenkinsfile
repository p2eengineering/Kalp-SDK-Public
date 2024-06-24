pipeline {
	agent any
	options {
        skipStagesAfterUnstable()
		timeout(time: 30, unit: 'MINUTES')
    	}
	environment {
		PROD_ECR_URL = '408153089286.dkr.ecr.ap-south-1.amazonaws.com/kalp-sdk-backend-prod'
		PROD_ENV = 'prod'
		SLACK_CHANNEL = 'pl-kalp-build-alerts'
    }
	stages {
		stage('PROD_BUILD') {
			when{
				branch 'main'
			}
			steps {
				script {
					committerEmail = sh (
      				script: 'git log -1 --pretty=format:"%an"', returnStdout: true
					).trim()
				}
				echo "Committer Email : '${committerEmail}'"
				slackSend (	color: 'good', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | Job has initiated : #${env.BUILD_NUMBER} by ${committerEmail}")
				slackSend (	color: 'warning', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | Build is started : #${env.BUILD_NUMBER}")
				echo "Step: BUILD, initiated..."
				sh "docker build -t '${PROD_ECR_URL}':'${BUILD_NUMBER}' . --no-cache"
				slackSend (	color: 'warning', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | Build has been completed : #${env.BUILD_NUMBER}")
			}
		}
		stage('PROD_ECR PUSH') {
			when{
				branch 'main'
			}
			steps {
				echo "Step: Pushing Image ..."
				sh "aws ecr get-login --no-include-email --region ap-south-1 | sh"
				sh "docker push '${PROD_ECR_URL}':${BUILD_NUMBER}"
           		slackSend (	color: 'warning', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | Build has been pushed to the ECR : #${env.BUILD_NUMBER}")
			}
		}
		stage('PROD_DEPLOY') {
			when{
				branch 'main'
			}
			steps {
				echo "Deploying into '${PROD_ENV}' environment"
				sh "aws eks --region ap-south-1 update-kubeconfig --name kalp-myipr-prod"
				sh "sed -i 's/<VERSION>/${BUILD_NUMBER}/g' deployment-'${PROD_ENV}'.yaml"
				sh "kubectl apply -f deployment-'${PROD_ENV}'.yaml"
				echo "'${PROD_ENV}' deployment completed: '${env.BUILD_ID}' on '${env.BUILD_URL}'"
           		slackSend (	color: 'warning', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | Deployment has been completed : #${env.BUILD_NUMBER}")
			}
		}
		stage('POST_CHECKS') {
			when{
				branch 'main'
			}
			steps {
				echo "POST test"
			}	
			post {
				always {
					echo "ALWAYS test1"
				}
				success {
					slackSend (	color: 'good', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | Job has succeeded : #${env.BUILD_NUMBER} in ${currentBuild.durationString.replace(' and counting', '')} \n For more info, please click (<${env.BUILD_URL}|here>)")
				}
				failure {
					slackSend (	color: 'danger', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | @channel - Job has failed #${env.BUILD_NUMBER}\nPlease check full info, (<${env.BUILD_URL}|here>)")
				}
			}
		}
	}
	post {
		always {
			echo "ALWAYS last-post check"
		}
		aborted {
			slackSend (
				color: '#AEACAC', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | Job has aborted : #${env.BUILD_NUMBER} in ${currentBuild.durationString.replace(' and counting', '')} \n For more info, please click (<${env.BUILD_URL}|here>)")
		}
		failure {
			slackSend (
				color: 'danger', channel: "${SLACK_CHANNEL}", message: "${env.JOB_NAME} | @channel - Job has failed #${env.BUILD_NUMBER}\nPlease check full info, (<${env.BUILD_URL}|here>)")
		}
	}
}
