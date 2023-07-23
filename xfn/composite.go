package main

import (
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
)

func GetNodeGroupEntries(cp *fieldpath.Paved) ([]NodeGroupsEntry, error) {
	nodeGroup := []NodeGroupsEntry{}
	if err := cp.GetValueInto("spec.nodeGroups", &nodeGroup); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.nodeGroups from observed composite")
	}

	return nodeGroup, nil
}

func GetResourceEntries(cp *fieldpath.Paved) (*ResourceParameters, error) {
	paramters := ResourceParameters{}
	if err := cp.GetValueInto("spec.region", &paramters.Region); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.region from observed composite")
	}

	if err := cp.GetValueInto("spec.clusterName", &paramters.ClusterName); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.clusterName from observed composite")
	}

	if err := cp.GetValueInto("spec.subnetIds", &paramters.SubnetIds); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.subnetIds from observed composite")
	}

	if err := cp.GetValueInto("spec.nodeRoleArn", &paramters.NodeRoleArn); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.nodeRoleArn from observed composite")
	}

	return &paramters, nil
}
