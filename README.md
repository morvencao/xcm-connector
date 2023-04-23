# xCM Connector

xCM connector is a controller that watches managed clusters and addons and synchronizes them to the xCM API server.

## Get Started

### Building

1. Build binary

```shell
make build
```

2. Build image

```shell
make image
```

_Note:_ Image name can be overriden by setting `IMAGE_NAME` environment variable, eg, run `IMAGE_NAME=quay.io/xxx/xcm-connector:latest` to build image with name `IMAGE_NAME=quay.io/foo/xcm-connector:latest`.

3. Push image

```shell
make push
```

### Deploying

1. Deploy xCM connector alongside MCE:

```shell
make deploy
```

_Note_: `xcm-server` in [deploy/deployment.yaml](deploy/deployment.yaml) may need to be changed to the correct xcm server if you're running testing against your own xCM API server.

2. Deploy xCM connector alongside multicluster-controlplane

_Note:_ Currently the addon for the multicluster-controlplane is not working as MCE, therefore if you deploy xCM connector with multicluster-controlplane, the addon sync will not work.
