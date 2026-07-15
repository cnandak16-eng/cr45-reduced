pipeline {
    agent any

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Frontend') {
            steps {
                sh '''
                docker run --rm \
                -v "$WORKSPACE/frontend:/app" \
                -w /app \
                node:22-alpine \
                sh -c "npm install && npm run build"
                '''
            }
        }

        stage('Build Backend') {
            steps {
                sh '''
                docker run --rm \
                -v "$WORKSPACE/backend:/app" \
                -w /app \
                golang:1.26-alpine \
                sh -c "go mod download && go build -o cr45-backend ./cmd/server"
                '''
            }
        }

        stage('Build Docker Images') {
            steps {
                sh 'docker compose build'
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