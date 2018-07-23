# Playbook example

The following example shows a real use case of Blackbeard usage.

This example describe a technical stack composed of :

* 2 versions of an API;
* A front-end app.

## Build my-app by yourself

If you want to fully test this example, you may want to build the docker images yourself. To do that, edit the `my-app/Makefile` file and change the following env var :

```sh
DOCKER_API_V1 = "seblegall/blackbeard-example-api:v1"
DOCKER_API_V2 = "seblegall/blackbeard-example-api:v2"
DOCKER_FRONT = "seblegall/blackbeard-example-front:v1"
```

*Replace the docker images name using your own docker hub namespace.*

## Using blackbeard to deploy my-app

### Requirement

You must have `kubectl` installed and configured and a Kubernetes cluster ready.

*Tips: On MacOS and Windows, you may use the built-in Kubernetes cluster with docker-for-desktop*

You must be located in the `playbook` directory : 

```sh
cd my-app-playbook
```

### Create a namespace using the v1 api

```sh
blackbeard create -n v1
blackbeard apply -n v1
```

Those commands will create a namespace called `v1` and deploy the API (using the v1 version) and the front-end app in this namespace.

You may check that the api v1 is actually running with this command :

```sh
kubectl logs {api_pod_name} -n v1
```

### Test api v2 in a different namespace

Now, if you'd like to run the API v2 in a different namespace, you may run :

```sh
blackbeard create -n v2
```

Edit the `inventories/v2_inventory.json` file and change the version value. This file now should look like :

```json
{
    "namespace": "v2",
    "values": {
        "api": {
            "version": "v2"
        },
        "front": {
            "version": "v1"
        }
    }
}
```

Finally, apply the changes :

```sh
blackbeard apply -n v2
```

Now, you can check that v2 is actually deployed by running the following command :

```sh
kubectl logs {api_pod_name} -n v2
```
