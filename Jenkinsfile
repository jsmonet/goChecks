node {
      ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
          withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
              env.Path="${GOPATH}/bin:$PATH"

          stage('Checkout'){
              echo 'Check me out'
              checkout scm
          }

          stage('Pull go tools and dependencies'){
              sh 'go get -u github.com/golang/dep/cmd/dep'
              sh 'go get -u github.com/golang/lint/golint'
              sh 'go get github.com/tebeka/go2xunit'
              sh 'go get ./...'
          }

          stage('Build'){
              sh 'go build -ldflags "-X main.version=$BUILDVER" -o check scheck.go'
          }

          stage('Validate version'){
              sh 'expr `echo $BUILDVER` = `./check -version`'
           }
       }
   }
}
