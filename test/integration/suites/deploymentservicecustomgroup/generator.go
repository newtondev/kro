// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
package deploymentservicecustomgroup_test

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	krov1alpha1 "github.com/awslabs/kro/api/v1alpha1"
	"github.com/awslabs/kro/pkg/testutil/generator"
)

const (
	customGroup      = "mycompany.it"
	customApiVersion = "v1alpha2"
)

// deploymentServiceCustomGroup creates a ResourceGroup for testing deployment+service combinations with custom group
func deploymentServiceCustomGroup(
	namespace, name string,
) (
	*krov1alpha1.ResourceGroup,
	func(namespace, name string, port int) *unstructured.Unstructured,
) {
	resourcegroup := generator.NewResourceGroup(name,
		generator.WithNamespace(namespace),
		generator.WithSchemaAndGroup(
			customGroup, "DeploymentServiceCustomGroup", customApiVersion,
			map[string]interface{}{
				"name": "string",
				"port": "integer | default=80",
			},
			map[string]interface{}{
				"deploymentConditions": "${deployment.status.conditions}",
				"availableReplicas":    "${deployment.status.availableReplicas}",
			},
		),
		generator.WithResource("deployment", deploymentDef(), nil, nil),
		generator.WithResource("service", serviceDef(), nil, nil),
	)
	instanceGenerator := func(namespace, name string, port int) *unstructured.Unstructured {
		return &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": fmt.Sprintf("%s/%s", customGroup, customApiVersion),
				"kind":       "DeploymentServiceCustomGroup",
				"metadata": map[string]interface{}{
					"name":      name,
					"namespace": namespace,
				},
				"spec": map[string]interface{}{
					"name": name,
					"port": port,
				},
			},
		}
	}
	return resourcegroup, instanceGenerator
}

func deploymentDef() map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": "${schema.spec.name}",
		},
		"spec": map[string]interface{}{
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app": "deployment",
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app": "deployment",
					},
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "${schema.spec.name}-deployment",
							"image": "nginx",
							"ports": []interface{}{
								map[string]interface{}{
									"containerPort": "${schema.spec.port}",
								},
							},
						},
					},
				},
			},
		},
	}
}

func serviceDef() map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"name": "${schema.spec.name}",
		},
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"app": "deployment",
			},
			"ports": []interface{}{
				map[string]interface{}{
					"port":       "${schema.spec.port}",
					"targetPort": "${schema.spec.port}",
				},
			},
		},
	}
}
