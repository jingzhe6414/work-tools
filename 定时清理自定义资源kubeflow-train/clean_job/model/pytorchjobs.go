package trainjobs

import (
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PyTorchJob struct {
	// Standard Kubernetes type metadata.
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired state of the PyTorchJob.
	Spec PyTorchJobSpec `json:"spec,omitempty"`

	// Most recently observed status of the PyTorchJob.
	// Read-only (modified by the system).
	Status commonv1.JobStatus `json:"status,omitempty"`
}

// PyTorchJobSpec is a desired state description of the PyTorchJob.
type PyTorchJobSpec struct {
	// RunPolicy encapsulates various runtime policies of the distributed training
	// job, for example how to clean up resources and how long the job can stay
	// active.
	//+kubebuilder:validation:Optional
	RunPolicy commonv1.RunPolicy `json:"runPolicy"`

	ElasticPolicy *ElasticPolicy `json:"elasticPolicy,omitempty"`

	// A map of PyTorchReplicaType (type) to ReplicaSpec (value). Specifies the PyTorch cluster configuration.
	// For example,
	//   {
	//     "Master": PyTorchReplicaSpec,
	//     "Worker": PyTorchReplicaSpec,
	//   }
	PyTorchReplicaSpecs map[commonv1.ReplicaType]*commonv1.ReplicaSpec `json:"pytorchReplicaSpecs"`
}

type ElasticPolicy struct {
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	MaxReplicas *int32 `json:"maxReplicas,omitempty"`

	RDZVBackend *RDZVBackend `json:"rdzvBackend,omitempty"`
	RDZVPort    *int32       `json:"rdzvPort,omitempty"`
	RDZVHost    *string      `json:"rdzvHost,omitempty"`
	RDZVID      *string      `json:"rdzvId,omitempty"`
	// RDZVConf contains additional rendezvous configuration (<key1>=<value1>,<key2>=<value2>,...).
	RDZVConf []RDZVConf `json:"rdzvConf,omitempty"`
	// Start a local standalone rendezvous backend that is represented by a C10d TCP store
	// on port 29400. Useful when launching single-node, multi-worker job. If specified
	// --rdzv_backend, --rdzv_endpoint, --rdzv_id are auto-assigned; any explicitly set values
	// are ignored.
	Standalone *bool `json:"standalone,omitempty"`
	// Number of workers per node; supported values: [auto, cpu, gpu, int].
	NProcPerNode *int32 `json:"nProcPerNode,omitempty"`

	MaxRestarts *int32 `json:"maxRestarts,omitempty"`

	// Metrics contains the specifications which are used to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).  The desired replica count is calculated with multiplying the
	// ratio between the target value and the current value by the current
	// number of pods. Ergo, metrics used must decrease as the pod count is
	// increased, and vice-versa.  See the individual metric source types for
	// more information about how each type of metric must respond.
	// If not set, the HPA will not be created.
	// +optional
	Metrics []autoscalingv2.MetricSpec `json:"metrics,omitempty"`
}

type RDZVConf struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type RDZVBackend string

// PyTorchJobList is a list of PyTorchJobs.
type PyTorchJobList struct {
	// Standard type metadata.
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of PyTorchJobs.
	Items []PyTorchJob `json:"items"`
}
