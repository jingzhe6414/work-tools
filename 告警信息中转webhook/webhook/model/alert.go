package model

type Alertmanager struct {
	Receiver    string      `json:"receiver,omitempty"`
	Status      string      `json:"status,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Alerts      []Alert     `json:"alerts,omitempty"`
	GroupLabels GroupLabel  `json:"groupLabels"`
	ExternalURL string      `json:"externalURL,omitempty"`
	Version     string      `json:"version,omitempty"`
	GroupKey    string      `json:"groupKey,omitempty"`
}

type Alert struct {
	Status       string     `json:"status,omitempty"`
	Labels       Label      `json:"labels,omitempty"`
	Annotations  Annotation `json:"annotations,omitempty"`
	StartsAt     string     `json:"startsAt,omitempty"`
	EndsAt       string     `json:"endsAt,omitempty"`
	ActiveAt     string     `json:"activeAt,omitempty"`
	GeneratorURL string     `json:"generatorURL,omitempty"`
	Fingerprint  string     `json:"fingerprint,omitempty"`
}

type Label struct {
	Alertname string `json:"alertname,omitempty"`
	Fzzn      string `json:"fzzn,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Pod       string `json:"pod,omitempty"`
	Severity  string `json:"severity,omitempty"`
}
type Annotation struct {
	Message string `json:"message,omitempty"`
	Tittle  string `json:"tittle,omitempty"`
}
type GroupLabel struct {
	Alertname string `json:"alertname,omitempty"`
}
