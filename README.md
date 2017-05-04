# cfseeker

Tool to find where apps on your Cloud Foundry are.
Very much WIP at the moment.

Right now, the local standalone version can tell you your app locations when
connected to a Cloud Foundry and BOSH.

## Goals

Support local standalone version, server with web UI, and cli to interact with server.

## Local Configuration

It's YAML!

```yaml
cf:
  api_address: https://<your-cf-host>.com
  client_id: your-client-user
  client_secret: supersecret
  skip_ssl_validation: true
bosh:
  api_address: https://<your-bosh-host>:25555
  username: your-username-or-client-id
  password: your-password-or-client-secret
  skip_ssl_validation: true
  deployments:
  - deployment-name-1
  - deployment-name-2
# server and its subkeys are only necessary if you're running in server mode
server:
  # basic auth creds if you want basic auth.
  basic_auth:
    username: admin
    password: password
  #no_auth: true  <set this to true and don't give basic auth creds if you want no auth
  port: 8892
```

## Running the Application

You can build it if you want - grab your favorite `go` distribution and build the files in the `cmd/cfseeker` directory. But let's be serious - you don't want to build it - head over to the releases page and there are binaries provided for you, free of charge.

Currently, the only supported command is `cfseeker find`. For more information on that, you can run `cfseeker help find`. You can also just run the `help` command for all the information you could ever want, or use the `--help` flag (you should also be able to use `-h`, but I made a mistake and that doesn't work in v0.1.0 _soooo_ next release).

## API Reference

If a non-2xx HTTP code is returned, then there will be a meta.error in the JSON
giving information about the error.

### Get Info about the CFSeeker Server

`GET /v1/meta`

**Example:**

```json
$ http localhost:8892/v1/meta
HTTP/1.1 200 OK
Content-Length: 31
Content-Type: application/json
Date: Thu, 04 May 2017 15:32:02 GMT

{
    "contents": {
        "version": "1234"
    }
}
```

### Get Info About Your STARTED Application

`GET /v1/apps`

#### Supported Arguments

This endpoint requires that either `app_guid` or all three of `org_name`,
`space_name`, and `app_name` are set.

* Option 1
  * `app_guid`: The GUID of your target application

* Option 2
  * `org_name`: The name of the CF organization your application is pushed to
  * `space_name`: The name of the CF space your application is pushed to
  * `app_name`: The name of your CF app, as it was pushed.

**Example:**

```json
$ http "admin:password@localhost:8892/v1/apps?app_guid=12345678-9abc-def1-2345-6789abcdef12"
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 02 May 2017 17:23:18 GMT

{
    "contents": {
        "guid": "12345678-9abc-def1-2345-6789abcdef12",
        "count": 2,
        "instances": [
            {
                "host": "10.244.2.133",
                "number": 0,
                "port": 61017,
                "deployment": "your-cloudfoundry",
                "vm_name": "runner_z1/0"
            },
            {
                "host": "10.244.2.134",
                "number": 1,
                "port": 61011,
                "deployment": "your-cloudfoundry",
                "vm_name": "runner_z1/1"
            }
        ],
        "name": "your-test-app"
    }
}
```
