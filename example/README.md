# Application Copy Pipeline Example

This is simple demo of the `cf-event` resource type being used to trigger a copy job to copy any application pushed to one space to a different space. A more practical use case would be to copy applications pushed to one Cloud Foundry target to other targets in a active-active platform configuration.

# Running the example

# Testing the task locally

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