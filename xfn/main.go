package main

import (
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/util/json"

	"github.com/crossplane/crossplane-runtime/pkg/password"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composed"

	"github.com/crossplane/crossplane-runtime/pkg/errors"

	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composite"
	"github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"
	"sigs.k8s.io/yaml"
)

func main() {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, errors.Wrap(err, "failed to read stdin"))
		os.Exit(1)
	}
	result, err := Run(in)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	_, _ = fmt.Fprint(os.Stdout, result)
}

func Run(in []byte) (string, error) {
	obj := &v1alpha1.FunctionIO{}
	if err := yaml.Unmarshal(in, obj); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal stdin")
	}
	xnodeGroup := composite.New()
	if err := yaml.Unmarshal(obj.Observed.Composite.Resource.Raw, &xnodeGroup.Unstructured); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal observed composite")
	}
	observedResources := make(map[string]*composed.Unstructured, len(obj.Observed.Resources))
	for _, entry := range obj.Observed.Resources {
		u := composed.New()
		if err := yaml.Unmarshal(entry.Resource.Raw, &u.Unstructured); err != nil {
			return "", errors.Wrap(err, "failed to unmarshal observed resource")
		}
		observedResources[entry.Name] = u
	}

	parameters, err := GetResourceEntries(fieldpath.Pave(xnodeGroup.Object))
	if err != nil {
		return "", errors.Wrap(err, "failed to get generic parameters from observed composite")
	}

	nodegroups, err := GetNodeGroupEntries(fieldpath.Pave(xnodeGroup.Object))
	if err != nil {
		return "", errors.Wrap(err, "failed to get nodegroups from observed composite")
	}

	resources, err := GetResources(nodegroups, *parameters)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate resources")
	}

	rawResources := make([]v1alpha1.DesiredResource, len(resources))
	for i, r := range resources {
		if o, ok := observedResources[r.Name]; ok {
			resources[i].Resource.SetName(o.GetName())
		} else {
			resources[i].Resource.SetName(generateObjectName(xnodeGroup.GetName()))
		}
		raw, err := json.Marshal(r.Resource)
		if err != nil {
			return "", errors.Wrapf(err, "failed to marshal generated resource %s", r.Name)
		}
		rawResources[i] = v1alpha1.DesiredResource{
			Name:     r.Name,
			Resource: runtime.RawExtension{Raw: raw},
		}
	}

	obj.Desired.Resources = append(obj.Desired.Resources, rawResources...)
	result, err := yaml.Marshal(obj)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal resulting functionio")
	}
	return string(result), nil
}

func generateObjectName(prefix string) string {
	suf, _ := password.Settings{
		CharacterSet: "abcdefghijklmnopqrstuvwxyz0123456789",
		Length:       5,
	}.Generate()
	return prefix + "-" + suf
}
