# Cloud Foundry Event Resource

Detects cloud foundry events within a given target.

## Source Configuration

* `api`: *Required.* The Cloud Foundry API endpoint.

* `user`: *Required.* The Cloud Foundry user.

* `password`: *Required.* The Cloud Foundry user's password.

* `org`: *Required.* The target `organization`.

* `space`: *Required.* The target `space`.

* `apps`: *Optional. List of application names for which events will be retrieved.

* `skip-ssl-validation`: *Optional. Skips verification of SSL certificates.

* `debug`: *Optional. Enables debug logging of the concourse resource type.

* `trace`: *Optional. Enable tracing HTTP request/responses to the Cloud Controller API.

### Example

``` yaml

resource_types:

- name: cf-event-type
  type: docker-image
  source:
    repository: mevansam/cf-event-type
    tag: latest

resources:

- name: cf-event
  type: cf-event-type
  source:
    api: https://api.local.pcfdev.io
    user: admin
    password: admin
    org: pcfdev-org
    space: pcfdev-space
    skip-ssl-validation: true
    apps:
    - app1
    - app2
```

## Behavior

### `check`: Check for new events.

Queries the Cloud Foundry target's Event API and retrieves events for the given applications or all apps within the space. Only events that occurred after that last cached event will be retrieved. If no events are cached then the last application push event will be returned.

### `in`: Fetch most recent events.

Creates content within the download directory that can be retrieved by jobs to determine the type of event, timestamp, etc.

Creates the following files in the download directory:

* `Env.sh` exports the values passed in within the source configuration as shell variables.

    Example:
    ``` bash
    export CF_API="https://api.local.pcfdev.io"
    export CF_USER="admin"
    export CF_PASSWORD="admin"
    export CF_ORG="pcfdev-org"
    export CF_SPACE="pcfdev-space"
    export CF_SKIP_SSL_VALIDATION="--skip-ssl-validation"
    export CF_APPS="app1 app2 "
    ```

* `version` list of applications and event data.
* `metadata` contains data around application events and when their timestamps
* `[app name].event` contains the event for the `app name`
* `[app name].timestamp` contains the timestamp for the event recorded for `app name` as a Unix timestamp integer value

#### Parameters

*Not Applicable*
