pipeline {
  agent {
    kubernetes {
      yaml '''
        apiVersion: v1
        kind: Pod
        spec:
          containers:
          - name: maven
            image: maven:alpine
            command:
            - cat
            tty: true
          - name: docker
            image: docker:latest
            command:
            - cat
            tty: true
            volumeMounts:
            - mountPath: /var/run/docker.sock
              name: docker-sock
          volumes:
          - name: docker-sock
            hostPath:
              path: /var/run/docker.sock
        '''
    }
  }
  triggers {
    githubPush()
  }
  options {
    buildDiscarder(
      logRotator(
        daysToKeepStr: '1',
        numToKeepStr: '2',
      )
    )
  }
  environment {
    registry = "gbsservice"
    registryCredential = 'docker-gbsservice'
    imageName = "bumilindo-link-api"
    tag = new Date().format('yy.MM.dd-HH.mm')
  }
  stages {
    stage('Clone-Repo') {
      steps {
        sh """
          [ -d ~/.ssh ] || mkdir ~/.ssh && chmod 0700 ~/.ssh
          ssh-keyscan -t rsa,dsa github.com >> ~/.ssh/known_hosts
        """
        git branch: 'release',
            credentialsId: 'bumilindo-link-api-key',
            url: 'git@github.com:gbsservice/api-dashboard-bumilink.git'
        sh "ls -lat"
      }
    }
    stage('Build-Docker-Image') {
      steps {
        container('docker') {
          sh 'docker build -t $registry/$imageName:$tag .'
        }
      }
    }
    stage('Login-Into-Docker') {
      steps {
        container('docker') {
          withCredentials([usernamePassword(credentialsId: registryCredential, usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
            sh 'docker login -u $USERNAME -p $PASSWORD'
          }
        }
      }
    }
    stage('Push-Images-Docker-to-DockerHub') {
      steps {
        container('docker') {
          sh 'docker push $registry/$imageName:$tag'
        }
      }
    }
    stage('Trigger ManifestUpdate') {
      steps {
        echo "triggering updatemanifestjob"
        build job: 'update-manifest', parameters: [
          string(name: 'GIT_CREDENTIAL_ID', value: 'bumilindo-cd'),
          string(name: 'GIT_URL', value: 'git@github.com:dianyulius/bumilindo_cd.git'),
          string(name: 'GIT_DEPLOYMENT_PATH', value: 'api/api-link.yaml'),
          string(name: 'DOCKER_IMAGE', value: "${registry}/${imageName}"),
          string(name: 'DOCKER_TAG', value: "${tag}")
        ]
      }
    }
  }
  post {
    always {
      container('docker') {
        sh 'docker logout'
      }
    }
  }
}
