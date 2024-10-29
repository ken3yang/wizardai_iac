## Description

The module will create a s3 bucket wil following rules:

- All data is encrypted at rest.
- All data is encrypted in transit.
- The bucket naming convention wizardai-<name>-<environment>.
- Bucket versioning enabled.
- All public access blocked.


## Usage

```
module "wizardai_s3_bucket_module" {
    source = "./wizardai_s3_bucket_module"
    name = "example"
    region = "ap-southeast-1"
    environment = "development"
    tags = {
        application: "sass"
    }
}
```

or refer examples/wizardai_s3_bucket_module.tf

## decisions
- Bucket versioning enabled and all public access blocked to enforce sensible security defaults.
- tests/wizardai_s3_bucket_module_test.go was inplemened for testing

## guidelines

### test
     `cd tests && go get . &&  go test -v`