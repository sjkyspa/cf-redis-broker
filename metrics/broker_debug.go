package metrics

type DebugResponse struct {
	Pool      Resource `yaml:"pool"`
	Allocated Resource `yaml:"allocated"`
}

type Resource struct {
	Count int `yaml:"count"`
}
