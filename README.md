# cfseeker

Tool to find where apps on your Cloud Foundry are.
Very much WIP at the moment.

## Goals

Support local standalone version, server with web UI, and cli to interact with server.

## Local configuration:

It's YAML!
```
cf:
  api_address: https://<your-host>.com
  client_id: your-client-user
  client_secret: supersecret
  skip_ssl_validation: true
```
