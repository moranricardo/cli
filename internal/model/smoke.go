package model

import "fmt"

type RunCommand string

const (
	UpdateFilesCommand RunCommand = "update"
	VersionCommand     RunCommand = "version"
	RecreateCommand    RunCommand = "recreate"
	SecurityCommand    RunCommand = "security"
	UpdateGraphCommand RunCommand = "graph"
)

type OutputType string

const (
	CreatePRType        OutputType = "create_pull_request"
	ClosePRType         OutputType = "close_pull_request"
	UpdateDepListType   OutputType = "update_dependency_list"
	MarkAsProcessedType OutputType = "mark_as_processed"
	RecordEcosystemType OutputType = "record_ecosystem_versions"
	RecordUpdateJobType OutputType = "record_update_job_warning"
)

type SmokeTest struct {
	Input  Input    `yaml:"input" json:"input"`
	Output []Output `yaml:"output,omitempty" json:"output,omitempty"`
}

type Input struct {
	Job         Job          `yaml:"job" json:"job"`
	Credentials []Credential `yaml:"credentials,omitempty" json:"credentials,omitempty"`
}

type Output struct {
	Type   OutputType    `yaml:"type" json:"type"`
	Expect UpdateWrapper `yaml:"expect" json:"expect"`
}

func (s SmokeTest) Validate() error {
	if s.Input.Job.PackageManager == "" {
		return fmt.Errorf("package-manager required")
	}
	return nil
}
