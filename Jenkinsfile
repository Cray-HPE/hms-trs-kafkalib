@Library('dst-shared@master') _

dockerBuildPipeline {
        repository = "cray"
        imagePrefix = "hms"
        app = "trs-kafkalib"
        name = "hms-trs-kafkalib"
        description = "Cray HMS TRS Kafka library package."
        dockerfile = "Dockerfile"
        slackNotification = ["", "", false, false, true, true]
        product = "internal"
}
