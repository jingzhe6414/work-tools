package trainjobs

import (
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TFJob struct {
	// Standard Kubernetes type metadata.
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired state of the TFJob.
	// +optional
	Spec TFJobSpec `json:"spec,omitempty"`

	// Most recently observed status of the TFJob.
	// Populated by the system.
	// Read-only.
	// +optional
	Status commonv1.JobStatus `json:"status,omitempty"`
}

// SuccessPolicy is the success policy.
type SuccessPolicy string

// TFJobSpec is a desired state description of the TFJob.
type TFJobSpec struct {
	// RunPolicy encapsulates various runtime policies of the distributed training
	// job, for example how to clean up resources and how long the job can stay
	// active.
	//+kubebuilder:validation:Optional
	RunPolicy commonv1.RunPolicy `json:"runPolicy"`

	// SuccessPolicy defines the policy to mark the TFJob as succeeded.
	// Default to "", using the default rules.
	// +optional
	SuccessPolicy *SuccessPolicy `json:"successPolicy,omitempty"`

	// A map of TFReplicaType (type) to ReplicaSpec (value). Specifies the TF cluster configuration.
	// For example,
	//   {
	//     "PS": ReplicaSpec,
	//     "Worker": ReplicaSpec,
	//   }
	TFReplicaSpecs map[commonv1.ReplicaType]*commonv1.ReplicaSpec `json:"tfReplicaSpecs"`

	// A switch to enable dynamic worker
	EnableDynamicWorker bool `json:"enableDynamicWorker,omitempty"`
}

type TFJobList struct {
	// Standard type metadata.
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of TFJobs.
	Items []TFJob `json:"items"`
}
