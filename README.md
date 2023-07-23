# crossplane xfn-nodegroup-loop

This Crossplane composition function relies on the CompositeResourceDefinition
of `XNodeGroup`

```
kubectl apply -f api/
kubectl apply -f example/
```

### Test locally.
```bash
cat test-input.yaml | go run .
```

### Build and push the image.
```bash
cd xfn
docker build \
  --tag docker.io/haarchri/xfn-nodegroup-loop:v0.5.0 \
  --push .
```
