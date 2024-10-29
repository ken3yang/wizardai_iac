module "wizardai_s3_bucket_module" {
    source = "../wizardai_s3_bucket_module"
    name = "example"
    region = "ap-southeast-1"
    environment = "development"
    tags = {
        application: "sass"
    }
}