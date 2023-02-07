package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fresh8gaming/gogen/internal/util"

	"github.com/icza/dyno"
	yaml "gopkg.in/yaml.v2"
)

func updateArgoApplication(argoApplicationFilePath, valuesPath string, service *Service) bool {
	valuesFilePath := fmt.Sprintf("%s/%s.yaml", valuesPath, service.ServiceName)

	fmt.Printf("Appending %s to %s\n", blue(valuesFilePath), argoApplicationFilePath)

	argoApplicationFileByte, err := os.ReadFile(argoApplicationFilePath)
	util.Fatal(err)

	var v interface{}
	err = yaml.Unmarshal(argoApplicationFileByte, &v)
	util.Fatal(err)

	valuesFiles, err := dyno.GetSlice(v, "spec", "source", "helm", "valueFiles")
	util.Fatal(err)

	exists := false
	for _, valuesFile := range valuesFiles {
		if valuesFile == valuesFilePath {
			exists = true
		}
	}

	if exists {
		return false
	}

	// Only edit the file if the value doesn't exist already
	err = dyno.Append(v, valuesFilePath, "spec", "source", "helm", "valueFiles")
	util.Fatal(err)

	editedArgoApplicationFileByte, err := yaml.Marshal(v)
	util.Fatal(err)

	util.Fatal(os.WriteFile(argoApplicationFilePath, editedArgoApplicationFileByte, os.ModePerm))
	return true
}

func updateMetadata(absPath string, service *Service, serviceType string) {
	fmt.Println("Appending to .metadata.yml")
	metadataFilePath := filepath.Join(absPath, ".metadata.yml")

	metadataFileByte, err := os.ReadFile(metadataFilePath)
	util.Fatal(err)

	var metadata Metadata
	err = yaml.Unmarshal(metadataFileByte, &metadata)
	util.Fatal(err)

	exists := false
	for _, existingService := range metadata.Services {
		if service.ServiceName == existingService.Name {
			exists = true
		}
	}

	if exists {
		return
	}

	newService := MetadataService{
		Name:      service.ServiceName,
		Type:      serviceType,
		CIEnabled: true,
	}

	metadata.Services = append(metadata.Services, newService)

	editedMetadataFileByte, err := yaml.Marshal(metadata)
	util.Fatal(err)

	util.Fatal(os.WriteFile(metadataFilePath, editedMetadataFileByte, os.ModePerm))
}

type Metadata struct {
	Name                   string            `yaml:"name"`
	Staging                bool              `yaml:"staging"`
	Team                   string            `yaml:"team"`
	Domain                 string            `yaml:"domain"`
	WhitesourceEnabled     bool              `yaml:"whitesourceEnabled"`
	KubescoreEnabled       bool              `yaml:"kubescoreEnabled"`
	CDEnabled              bool              `yaml:"cdEnabled"`
	Services               []MetadataService `yaml:"services"`
	ArgoAppNamesProduction string            `yaml:"argoAppNamesProduction"`
	ArgoAppNamesStaging    string            `yaml:"argoAppNamesStaging"`
	SkipInstallTools       bool              `yaml:"skipInstallTools"`

	Deploy Deploy `yaml:"deploy"`
}

type Deploy struct {
	Platform string `yaml:"platform"`
	Product  string `yaml:"product"`
}

type MetadataService struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	CIEnabled  bool   `yaml:"ciEnabled"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
}
