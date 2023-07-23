package main

import (
	v1beta1 "github.com/upbound/provider-aws/apis/eks/v1beta1"
)

type ResourceParameters struct {
	// Region is the region you'd like your resource to be created in.
	Region *string `json:"region"`

	// ClusterName is the Name you'd like your NodeGroup attached to specific EKS Cluster
	ClusterName *string `json:"clusterName"`

	// Identifiers of EC2 Subnets to associate with the EKS Node Group. Amazon EKS managed node groups can be launched in both public and private subnets. If you plan to deploy load balancers to a subnet, the private subnet must have tag kubernetes.io/role/internal-elb, the public subnet must have tag kubernetes.io/role/elb.
	SubnetIds []*string `json:"subnetIds"`

	// Amazon Resource Name (ARN) of the IAM Role that provides permissions for the EKS Node Group.
	NodeRoleArn *string `json:"nodeRoleArn,omitempty"`
}

type NodeGroupsEntry struct {
	// Name of the NodeGroup.
	Name string `json:"name"`

	// List of instance types associated with the EKS Node Group.
	InstanceTypes []*string `json:"instanceTypes"`

	// Configuration block with scaling settings. Detailed below.
	ScalingConfig []ScalingConfigParameters `json:"scalingConfig"`
}

type ScalingConfigParameters struct {

	// Desired number of worker nodes.
	DesiredSize *float64 `json:"desiredSize" tf:"desired_size,omitempty"`

	// Maximum number of worker nodes.
	MaxSize *float64 `json:"maxSize" tf:"max_size,omitempty"`

	// Minimum number of worker nodes.
	MinSize *float64 `json:"minSize" tf:"min_size,omitempty"`
}

type DesiredResource struct {
	Name     string
	Resource *v1beta1.NodeGroup
}

func GetResources(nodegroups []NodeGroupsEntry, parameters ResourceParameters) ([]DesiredResource, error) {
	resources := []DesiredResource{}
	for _, nodegroup := range nodegroups {
		resources = append(resources, DesiredResource{
			Name:     nodegroup.Name,
			Resource: WrapForAWS(parameters, nodegroup.InstanceTypes, nodegroup.ScalingConfig),
		})
	}
	return resources, nil
}
