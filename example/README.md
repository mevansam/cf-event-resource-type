# Application Copy Pipeline Example

This is simple demo of the `cf-event` resource type being used to trigger a copy job to copy any application pushed to one space to a different space. A more practical use case would be to copy applications pushed to one Cloud Foundry target to other targets in an active-active platform configuration.

# Running the example

To run the example you need to have a [PCFDev](https://github.com/pivotal-cf/pcfdev) instance running locally. Within this instance create the following two spaces.

```
cf create-space dev1 -o pcfdev-org
cf create-space dev2 -o pcfdev-org
```

Applications deployed to the `dev1` space will be copied over to the `dev2` space via the pipeline set-up as shown below. You will need to have [Vagrant](https://www.vagrantup.com) and [Virtualbox](https://www.virtualbox.org) installed in order to run the [Concourse CI](http://concourse.ci/) environment locally.

```
vagrant up

fly --target local login --concourse-url http://192.168.100.4:8080
fly --target local sync

fly --target local set-pipeline --pipeline app-copy-pipeline --config app-copy-pipeline.yml
```

# Testing the task locally

You can also run the task locally by running the `check.go` and `in.go` command code via the [Go](https://golang.org/) as shown below. This can be useful for debugging the Concourse resource type code.

> To run the examples below you will need the [`jq`](https://stedolan.github.io/jq/) tool available within the system to path to be able to parse the JSON results

```
SOURCE=$(echo '{
    "api": "https://api.local.pcfdev.io",
    "user": "admin",
    "password": "admin",
    "org": "pcfdev-org",
    "space": "dev1",
    "skip-ssl-validation": true
}')

VERSION=$(echo "{ 
    \"source\": $SOURCE
}" | go run ../cmd/check/check.go | jq .[0])

TMPDIR=$(mktemp -d)
echo "{ 
    \"source\": $SOURCE,
    \"version\": $VERSION
}" | go run ../cmd/in/in.go $TMPDIR | jq .

fly --target local \
    execute --config copy-apps.yml \
    --input ci-pipeline-source=. \
    --input cf-event=$TMPDIR
```