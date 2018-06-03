node("testVmSlaver") {
timestamps {
ansiColor('xterm') {

  stage('cleanup and checkout') {
    step([$class: 'WsCleanup'])
    checkout scm
    sh "sudo mkdir -p ~/.gocache && sudo chmod -R 777 ~/.gocache"
    sh "docker system prune -f"

  }

  stage('lint') {
    try {
      timeout(time: 30, unit: 'MINUTES') {
        sh "skipper make lint"
      }
    } catch (err) {
      notifyFailed("Test failure", err)
      throw err
    }
  }

  stage('test') {
    try {
      timeout(time: 30, unit: 'MINUTES') {
        sh "skipper make test"
      }
    } catch (err) {
      notifyFailed("Test failure", err)
      throw err
    }
  }

  stage('build') {
    try {
      timeout(time: 30, unit: 'MINUTES') {
        sh "skipper make build"
      }
    } catch (err) {
      notifyFailed("Build failure", err)
      throw err
    }
  }

  stage('subsystem') {
    try {
      timeout(time: 30, unit: 'MINUTES') {
        sh "skipper make subsystem"
      }
    } catch (err) {
      notifyFailed("Subsystem failure", err)
      throw err
    } finally {
      junit testResults: 'build/subsystem.xml'
      archiveArtifacts artifacts: 'subsystem/logs/*.log'
    }
  }

  notifySuccess("Successful deployed")

}}}

def notifySuccess(message) {
  if (env.BRANCH_NAME != 'master') {
    return
  }
  slackSend(
    message: "${message} ${env.BRANCH_NAME} [${env.BUILD_NUMBER}] (${env.BUILD_URL})",
    color: "#00FF00",
    channel: "#",
    token: "pP9r6dWVWCXqljdWJI1h2ItF",
  )
}

def notifyFailed(message, err) {
  if (env.BRANCH_NAME != 'master') {
    return
  }
  slackSend(
    message: "FAILED: ${message} '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})",
    color: "#FF0000",
    channel: "#",
    token: "pP9r6dWVWCXqljdWJI1h2ItF",
  )
}
