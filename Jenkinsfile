pipeline {
    agent any

    environment {
        IMAGE_NAME = "mediconnect-backend"
        REGISTRY   = "registry.example.com"   // replace with your registry
        GIT_SHA    = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
        IMAGE_TAG  = "${REGISTRY}/${IMAGE_NAME}:${GIT_SHA}"
    }

    options {
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {

        // ── S1: Checkout & Lint ────────────────────────────────────────────
        stage('Checkout & Lint') {
            steps {
                checkout scm
                sh 'go version'
                sh 'golangci-lint run ./... --timeout 5m'
            }
        }

        // ── S2: Unit Tests ────────────────────────────────────────────────
        stage('Unit Tests') {
            steps {
                sh '''
                    go test -v -race -count=1 \
                        -coverprofile=coverage.out \
                        ./...
                    go tool cover -func=coverage.out | tail -1
                '''
            }
            post {
                always {
                    archiveArtifacts artifacts: 'coverage.out', allowEmptyArchive: true
                }
            }
        }

        // ── S3: Integration Tests ─────────────────────────────────────────
        stage('Integration Tests') {
            steps {
                sh 'docker compose up -d db redis rabbitmq'
                sh 'sleep 10'  // wait for healthchecks
                sh 'go test -v -tags=integration ./...'
            }
            post {
                always {
                    sh 'docker compose down'
                }
            }
        }

        // ── S4: Build & Push Docker Image ─────────────────────────────────
        stage('Build & Push Image') {
            steps {
                sh "docker build -t ${IMAGE_TAG} ."
                sh "trivy image --exit-code 1 --severity CRITICAL ${IMAGE_TAG}"
                withCredentials([usernamePassword(
                    credentialsId: 'registry-credentials',
                    usernameVariable: 'REG_USER',
                    passwordVariable: 'REG_PASS'
                )]) {
                    sh "docker login ${REGISTRY} -u ${REG_USER} -p ${REG_PASS}"
                    sh "docker push ${IMAGE_TAG}"
                }
            }
        }

        // ── S5: Deploy to Staging ─────────────────────────────────────────
        stage('Deploy Staging') {
            steps {
                sh """
                    kubectl set image deployment/mediconnect-backend \
                        mediconnect-backend=${IMAGE_TAG} \
                        -n staging
                    kubectl rollout status deployment/mediconnect-backend \
                        -n staging --timeout=3m
                """
            }
        }

        // ── S6: Smoke Test & Manual Gate ──────────────────────────────────
        stage('Smoke Test') {
            steps {
                sh 'sleep 5'
                sh 'curl -sf https://staging.mediconnect.id/api/v1/health'
            }
        }

        stage('Approval Gate') {
            steps {
                input message: 'Deploy to production?', ok: 'Ship it 🚀'
            }
        }

        // ── S7: Deploy to Production ──────────────────────────────────────
        stage('Deploy Production') {
            steps {
                sh """
                    kubectl set image deployment/mediconnect-backend \
                        mediconnect-backend=${IMAGE_TAG} \
                        -n production
                    kubectl rollout status deployment/mediconnect-backend \
                        -n production --timeout=5m
                """
            }
        }
    }

    post {
        failure {
            echo "Pipeline failed — triggering rollback on production namespace"
            sh 'kubectl rollout undo deployment/mediconnect-backend -n production || true'
        }
        always {
            cleanWs()
        }
    }
}
