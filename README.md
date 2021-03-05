# Stupid Finalizers

## Build

`go build`

## Run

```
# List all objects blocked by a finalizer
./finalizers

# List all objects with finalizers
./finalizers --all
```

## Example

```
./finalizers

NAMESPACE   NAME                    APIVERSION                KIND                         FINALIZERS
p-nf5gh     creator-project-owner   management.cattle.io/v3   ProjectRoleTemplateBinding   [controller.cattle.io/mgmt-auth-prtb-controller]
```
