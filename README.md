# cfseeker

Tool to find where apps on your Cloud Foundry are.
Very much WIP at the moment.

Right now, the local standalone version can tell you your app locations when
connected to a Cloud Foundry and BOSH.

## Goals

Support local standalone version, server with web UI, and cli to interact with server.

## Local configuration:

It's YAML!
```
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
```
