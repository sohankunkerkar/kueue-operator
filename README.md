# Kueue Operator

Kueue Operator provides the ability to deploy kueue using different configurations

## Releases

| ko version   | ocp version         |kueue version  | k8s version | golang |
| ------------ | ------------------- |---------------| ----------- | ------ |
| 1.0.0        | 4.19 - 4.20         |0.11.z         | 1.32        | 1.23   |

Kueue releases around 6 times a year.
For the latest Openshift version, we will take the latest version that was build with that underlying
Kubernetes version.

See [Kueue Release](https://github.com/kubernetes-sigs/kueue/blob/main/RELEASE.md) for more details
on the Kueue release policy.

## Deploy the Operator

### Quick Development

1. Build and push the operator image to a registry:

   ```sh
   export QUAY_USER=${your_quay_user_id}
   export IMAGE_TAG=${your_image_tag}
   podman build -t quay.io/${QUAY_USER}/kueue-operator:${IMAGE_TAG} .
   podman login quay.io -u ${QUAY_USER}
   podman push quay.io/${QUAY_USER}/kueue-operator:${IMAGE_TAG}
   ```

1. Update the image spec under `.spec.template.spec.containers[0].image` field in the `deploy/08_deployment.yaml` Deployment to point to the newly built image

1. Update the `.spec.image` field under `deploy/09_kueue_crd.yaml` CR to point to a kueue image

1. Apply the manifests from `deploy` directory:

   ```sh
   oc apply -f deploy/
   ```

## Sample CR

```yaml
apiVersion: operator.openshift.io/v1alpha1
kind: Kueue
metadata:
  labels:
    app.kubernetes.io/name: kueue-operator
    app.kubernetes.io/managed-by: kustomize
  name: cluster
  namespace: openshift-kueue-operator
spec:
  image: "RHEL-Released-Image"
  config:
    integrations:
      frameworks:
      - "batch/job" 
```