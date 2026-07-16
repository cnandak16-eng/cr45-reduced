pipeline {
    agent any

    environment {
        IMAGE_NAME = "cnk19/cr45"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    sh '''
                    npm install
                    npm run build
                    '''
                }
            }
        }

        stage('Build Backend') {
            steps {
                dir('backend') {
                    sh '''
                    go mod download
                    go build -o cr45-backend ./cmd/server
                    '''
                }
            }
        }
        stage('Build Docker Images') {
            steps {
        sh '''
        docker image rm -f cr45-pipeline-frontend:latest || true
        docker image rm -f cr45-pipeline-backend:latest || true

        docker compose build --no-cache
        '''
        }
    }
    stage('Verify Frontend Image') {
    steps {
        sh '''
        docker run --rm cr45-pipeline-frontend:latest \
            cat /etc/nginx/conf.d/default.conf
        '''
    }
}

        stage('Push Docker Images') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                    echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin

                    docker tag cr45-pipeline-frontend:latest cnk19/cr45-frontend:latest
                    docker tag cr45-pipeline-backend:latest cnk19/cr45-backend:latest

                    docker push cnk19/cr45-frontend:latest
                    docker push cnk19/cr45-backend:latest
                    '''
                }
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                docker compose down || true
                docker compose up -d
                '''
            }
        }
    }

    post {
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}