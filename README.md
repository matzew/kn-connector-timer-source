# time source poc

```
ko apply -f config/
```

```
kubectl apply -f - <<EOF
apiVersion: samples.knative.dev/v1alpha1
kind: KameletSource
metadata:
  name: kamelet-source
  namespace: knative-samples
spec:
  type: "timer-source"
  text: "Hello receiver 2"
  sink:
    ref:
      apiVersion: eventing.knative.dev/v1
      kind: Broker
      name: default-2
EOF
```

Note: requires knative and `pre-setup` manifest
