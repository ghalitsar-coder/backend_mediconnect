pipeline {
    agent any

    environment {
        BACKEND_IMAGE_NAME = 'ghalitsar/backend-app:v2'
        
        // ===== KONFIGURASI CREDENTIALS JENKINS =====
        // Pastikan Anda membuat ID relasi credentials berikut di Jenkins (Manage Jenkins -> Credentials)
        
        // 1. DockerHub Credentials (tipe Username with password)
        DOCKERHUB_CREDS_ID = 'dockerhub-credentials'
        
        // 2. SSH Private Key Credentials (tipe SSH Username with private key)
        // Kunci ini harus bisa mengakses frontend (sebagai proxy) dan backend
        SSH_CREDS_ID = 'ssh-credentials' 
        
        // 3. Secret text credentials untuk Variabel Lingkungan
        BACKEND_PRIVATE_IP   = credentials('backend-private-ip')
        BACKEND_USERNAME     = credentials('backend-username')
        FRONTEND_HOST        = credentials('frontend-host')
        FRONTEND_USERNAME    = credentials('frontend-username')
        AZURE_STORAGE_STRING = credentials('azure-storage-connection-string')
        DB_URL               = credentials('db-url')
    }

    options {
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build & Push Image') {
            steps {
                script {
                    dir('backend_mediconnect') {
                        docker.withRegistry('https://index.docker.io/v1/', "${DOCKERHUB_CREDS_ID}") {
                            def backendImage = docker.build("${BACKEND_IMAGE_NAME}")
                            backendImage.push()
                        }
                    }
                }
            }
        }

        stage('Prepare Image Tarball') {
            steps {
                dir('backend_mediconnect') {
                    sh """
                        docker pull ${BACKEND_IMAGE_NAME}
                        echo "Menyimpan image menjadi file .tar..."
                        docker save -o backend-image.tar ${BACKEND_IMAGE_NAME}
                    """
                }
            }
        }

        stage('Deploy via SSH ProxyJump') {
            steps {
                dir('backend_mediconnect') {
                    // Menggunakan plugin sshagent untuk otomatis memuat private key
                    sshagent(credentials: ["${SSH_CREDS_ID}"]) {
                        sh """
                            echo "Mengkonfigurasi SSH Proxyjump..."
                            mkdir -p ~/.ssh
                            cat <<EOF > ~/.ssh/config
Host frontend-proxy
  HostName ${FRONTEND_HOST}
  User ${FRONTEND_USERNAME}
  StrictHostKeyChecking no

Host backend-target
  HostName ${BACKEND_PRIVATE_IP}
  User ${BACKEND_USERNAME}
  ProxyJump frontend-proxy
  StrictHostKeyChecking no
EOF
                            chmod 600 ~/.ssh/config

                            echo "Menyalin file docker-compose, schema, dan image.tar ke Backend..."
                            scp backend-image.tar docker-compose.yml mediconnect_id_schema.sql backend-target:/home/${BACKEND_USERNAME}/

                            echo "Menjalankan deployment di server backend..."
                            ssh backend-target << 'EOF'
                                echo "Memuat image ke Docker dari file .tar lokal..."
                                docker load -i /home/${BACKEND_USERNAME}/backend-image.tar
                                
                                cd /home/${BACKEND_USERNAME} || true
                                
                                echo "Membuat file .env..."
                                echo "AZURE_STORAGE_CONNECTION_STRING='${AZURE_STORAGE_STRING}'" > .env
                                echo "DB_URL='${DB_URL}'" >> .env
                                
                                docker-compose up -d
                                
                                echo "Membersihkan file .tar..."
                                rm -f /home/${BACKEND_USERNAME}/backend-image.tar
                                exit
EOF
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}
