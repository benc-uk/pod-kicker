# Pod Kicker

This is small proof of concept of a component to restart

# Getting Started

## Installing / Deploying

See [examples in deploy directory](./deploy/) for how to add the podkicker sidecar to any pod, and the service account needed.

## Running locally

The code only works running inside a pod in Kubernetes

# Configuration

Details of any configuration files, environmental variables, command line parameters, etc.

For services
| Setting / Variable | Purpose | Default |
| ------------------ | ------------------------------------------- | ------- |
| PODKICKER_WATCH | What file or directory to watch for changes, when a directory is added all files and sub-directories under it are watched recursively. **_Required_** | _None_ |
| PODKICKER_TARGET_NAME | The name of the Kubernetes deployment or stateful set to be restarted. Typically this is the same one as the sidecar is running under. **_Required_** | _None_ |
| PODKICKER_TARGET_TYPE | Either "deployment" or "statefulset" | "deployment" |

# Repository Structure

A brief description of the top-level directories of this project is as follows:

```c
/build      - Build configuration e.g. Dockerfiles
/deploy     - Deployment and infrastructure as code, inc Kubernetes
/cmd        - Source code
```

# Known Issues

All file operations will trigger the restart, there is no filter yet

# License

This project uses the MIT software license. See [full license file](./LICENSE)

