# cfseeker

Tool to find where apps on your Cloud Foundry are.

Can run as a local standalone version which can tell you your app locations
when configured to connect to a Cloud Foundry and (optionally) a BOSH. It can
also run as a server with an API that will allow you to access the functionality
of the tool through an HTTP API (documented below).
There is now a Web UI that uses the HTTP API which can be accessed by navigating
your browser to the root (/) endpoint of a listening cfseeker server.
Using the tool with the `--target` (`-t`) flag set as a cfseeker URL to contact,
the tool can also run its commands as a CLI that contacts the web API for you.

## Web UI

Point your browser at the root endpoint of the url where you have the server listening.
Click on links, type into fields. If you can't figure it out from there, then the
repo probably needs an issue posted.

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
  cache_ttl: 6000 #time in seconds to hold cache entries
  port: 8892
```

## Running the Application

You can build it if you want - grab your favorite `go` distribution and build the files in the `cmd/cfseeker` directory. But let's be serious - you don't want to build it - head over to the releases page and there are binaries provided for you, free of charge.

Currently, the only supported command is `cfseeker find`. For more information on that, you can run `cfseeker help find`. You can also just run the `help` command for all the information you could ever want, or use the `--help` flag.

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

**Supported Arguments:**

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
        "count": 2,
        "guid": "12345678-9abc-def1-2345-6789abcdef12",
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

### Clear the BOSH VM Info Cache

`DELETE /v1/cache/bosh`

 If your VM mappings are out of date and you don't want to wait for the cache
 TTL, you can force a cache reset by calling this endpoint.

**Example:**

```json
$ http DELETE "admin:password@localhost:8892/v1/cache/bosh"
HTTP/1.1 200 OK
Content-Length: 62
Content-Type: application/json
Date: Thu, 04 May 2017 18:40:55 GMT

{
    "meta": {
        "message": "BOSH VM info cache successfully cleared"
    }
}
```

### Convert a GUID to Names or Vice Versa

`GET /v1/convert`

**Supported Arguments:**

When getting information by name, an app name must be qualified by a space name,
and a space name must be qualified by an organization name.
You can get information about an app by providing an org, space, and app
name. You could get info about a space by providing a name for a space
and org only. Info about an org only requires an org name.

* Option 1
  * `app_guid`: The GUID of your organization, space, or application

* Option 2
  * `org_name`: The name of the CF organization your application is pushed to
  * `space_name`: The name of the CF space your application is pushed to
  * `app_name`: The name of your CF app, as it was pushed.

* Option 3
  * `org_name`: The name of the CF organization your space is within
  * `space_name`: The name of your desired space

* Option 4
  * `org_name`: The name of your desired organization name

**Example:**

```json
$ http admin:password@localhost:8892/v1/convert?guid=01234567-89ab-cdef-0123-456789abcdef
HTTP/1.1 200 OK
Content-Length: 254
Content-Type: application/json
Date: Thu, 27 Jul 2017 16:41:47 GMT

{
    "contents": {
        "app_guid": "01234567-89ab-cdef-0123-456789abcdef",
        "app_name": "cfseeker",
        "org_guid": "3456789a-bcde-f012-3456-789abcdef012",
        "org_name": "cfseeker-org",
        "space_guid": "6789abcd-ef01-2345-6789-abcdef012345",
        "space_name": "cfseeker-space",
        "type": "app"
    }
}
```