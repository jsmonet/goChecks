node {
      dir("$JENKINS_HOME/go"){
        if(!fileExists("/")){
          sh 'mkdir $JENKINS_HOME/go'
        }

      }
      ws("${JENKINS_HOME}/jobs/go/src/github.com/jsmonet/goChecks") {
          withEnv(["GOPATH=${JENKINS_HOME}/jobs/go"]) {
              env.Path="${JENKINS_HOME}/tools/org.jenkinsci.plugins.golang.GolangInstallation/go1.10/bin:${GOPATH}/bin:$PATH"
        dir("${JENKINS_HOME}/jobs/go/src/github.com/jsmonet/goChecks"){
            stage('clean out workspce'){
                sh 'rm -rf ./*'
                sh 'echo $PWD'
                sh 'ls -al'
            }
            stage('Checkout'){
                echo 'Check me out'
                checkout scm
            }

            stage('Pull go tools and dependencies'){
                sh 'go get -u github.com/golang/dep/cmd/dep'
                sh 'go get -u github.com/golang/lint/golint'
                sh 'go get github.com/tebeka/go2xunit'
                sh 'dep ensure'
            }
            dir("${JENKINS_HOME}/jobs/go/src/github.com/jsmonet/goChecks/grab"){
                stage('Test grab'){
                    sh 'go test'
                }
            }
            dir("${JENKINS_HOME}/jobs/go/src/github.com/jsmonet/goChecks/validify"){
                stage('Test validify'){
                    sh 'go test'
                }
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
}
