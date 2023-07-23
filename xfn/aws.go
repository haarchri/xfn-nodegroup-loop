package main

import (
	v1beta1 "github.com/upbound/provider-aws/apis/eks/v1beta1"
)

func WrapForAWS(parameters ResourceParameters, instanceType []*string, scalingConfigs []ScalingConfigParameters) *v1beta1.NodeGroup {
	o := &v1beta1.NodeGroup{
		Spec: v1beta1.NodeGroupSpec{
			ForProvider: v1beta1.NodeGroupParameters{
				InstanceTypes: instanceType,
				Region:        parameters.Region,
				ClusterName:   parameters.ClusterName,
				NodeRoleArn:   parameters.NodeRoleArn,
			},
			// ResourceSpec: v1.ResourceSpec{
			// 	ProviderConfigReference: &v1.Reference{
			// 		Name: providerConfigName,
			// 	},
			// },
		},
	}

	for _, subnetId := range parameters.SubnetIds {
		o.Spec.ForProvider.SubnetIds = append(o.Spec.ForProvider.SubnetIds, subnetId)
	}
	for _, scalingConfig := range scalingConfigs {
		o.Spec.ForProvider.ScalingConfig = append(o.Spec.ForProvider.ScalingConfig, v1beta1.ScalingConfigParameters{
			DesiredSize: scalingConfig.DesiredSize,
			MinSize:     scalingConfig.MinSize,
			MaxSize:     scalingConfig.MaxSize,
		})
	}
	o.SetGroupVersionKind(v1beta1.NodeGroup_GroupVersionKind)
	return o
}
