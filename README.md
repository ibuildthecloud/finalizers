# Stupid Finalizers

## Build

`go build`

## Run

```
# List all objects with a finalizer
./finalizers

# List objects not deleting because of finalizers
./finalizers --blocked
```

## Example

```
./finalizers --blocked

NAMESPACE   NAME                    APIVERSION                KIND                         FINALIZERS
p-nf5gh     creator-project-owner   management.cattle.io/v3   ProjectRoleTemplateBinding   [controller.cattle.io/mgmt-auth-prtb-controller]
```
