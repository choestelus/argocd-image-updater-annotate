# argocd-image-updater-annotate
when kustomize has too many images... ResidentSleeper

## Usage

with following kustomization manifest as `kustomization.yaml`
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: foo-example

images:
- name: foo-example-service-image
  newName: example.com/path/to/foo-example-service
  newTag: v0.1.0
- name: bar-example-service-image
  newName: example.com/path/to/bar-example-service
  newTag: v1.1.0
```

Run following command

```sh
go run main.go kustomization.yaml > manifest.yaml
```

will produces Argo CD application manifest as following
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  annotations:
    argocd-image-updater.argoproj.io/foo-example-service.kustomize.image-name: 'foo-example-service-image'
    argocd-image-updater.argoproj.io/foo-example-service.update-strategy: 'latest'
    argocd-image-updater.argoproj.io/bar-example-service.kustomize.image-name: 'bar-example-service-image'
    argocd-image-updater.argoproj.io/bar-example-service.update-strategy: 'latest'
    argocd-image-updater.argoproj.io/image-list: 'foo-example-service=example.com/path/to/foo-example-service,bar-example-service=example.com/path/to/bar-example-service'
```

## Caveats
- malformed kustomization, particulary `images` section without complete `name` `newName` and `newTag` presence might cause undefined behavior
use this only if you have linting/vetting process in-place to ensure kustomization format
- no warning when generated service referencing annotations exceed 63 characters limit as described in RFC-1035/RFC-1123

## License
This project is licensed under Apache License 2.0

