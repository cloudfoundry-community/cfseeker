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

## Running the Application
You can build it if you want - grab your favorite `go` distribution and build the files in the `cmd/cfseeker` directory. But let's be serious - you don't want to build it - head over to the releases page and there are binaries provided for you, free of charge.

Currently, the only supported command is `cfseeker find`. For more information on that, you can run `cfseeker help find`. You can also just run the `help` command for all the information you could ever want, or use the `--help` flag (you should also be able to use `-h`, but I made a mistake and that doesn't work in v0.1.0 _soooo_ next release).
