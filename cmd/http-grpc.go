package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fresh8gaming/gogen/internal/util"

	"github.com/icza/dyno"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const HTTPGRPCTemplates string = "_template/service/http-grpc"

//go:embed _template/service/http-grpc/*
var httpGRPCContent embed.FS

type Service struct {
	Name                  string
	Org                   string
	ServiceName           string
	ServiceNameUnderscore string
	ServiceNameProto      string
}

func runHTTPCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		createHTTPGRPCService(args)
	}
}

func runGRPCCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		createHTTPGRPCService(args)
	}
}

func createHTTPGRPCService(args []string) {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	absPath, err := filepath.Abs(target)
	util.Fatal(err)

	if fileStat, err := os.Stat(absPath); os.IsNotExist(err) {
		err := os.MkdirAll(absPath, os.ModePerm)
		util.Fatal(err)
	} else if !fileStat.IsDir() {
		log.Fatalf("%s is not a directory", absPath)
	}

	argoApplicationFilePath := filepath.Join(absPath, "deploy", "argocd", "application.yaml")
	if _, err := os.Stat(argoApplicationFilePath); os.IsNotExist(err) {
		log.Fatalf("argocd application file expected at %s", argoApplicationFilePath)
	}
	stagingExists := true
	stagingArgoApplicationFilePath := filepath.Join(absPath, "deploy", "argocd", "staging.yaml")
	if _, err := os.Stat(stagingArgoApplicationFilePath); os.IsNotExist(err) {
		stagingExists = false
		log.Printf("staging argocd application file expected at %s", argoApplicationFilePath)
	}

	service := &Service{
		Name:                  getName(Name, absPath),
		Org:                   Org,
		ServiceName:           ServiceName,
		ServiceNameUnderscore: strings.ReplaceAll(ServiceName, "-", "_"),
		ServiceNameProto:      strings.ReplaceAll(strings.ReplaceAll(ServiceName, "_", ""), "-", ""),
	}

	fmt.Printf("Creating %s in %s\n", blue(service.ServiceName), absPath)

	updateMetadata(absPath, service, "http-grpc")

	copyTemplates(absPath, HTTPGRPCTemplates, service, httpGRPCContent, HTTPGRPCTemplates, func(path string) string {
		replaced := strings.ReplaceAll(path, "service-name", service.ServiceName)
		replaced = strings.ReplaceAll(replaced, "service_name", service.ServiceNameUnderscore)
		return replaced
	})

	updatedArgo := updateArgoApplication(argoApplicationFilePath, "values-production", service)

	var updatedStagingArgo bool
	if stagingExists {
		updatedStagingArgo = updateArgoApplication(stagingArgoApplicationFilePath, "values-staging", service)
	}

	fmt.Printf("Created %s!\n", green(service.ServiceName))
	fmt.Println()
	fmt.Println("Run the following commands to set up the project:")
	fmt.Println()
	fmt.Printf(blue("cd %s\n"), absPath)
	fmt.Println(blue("go fmt ./..."))
	fmt.Println(blue("make buf-mod-update"))
	fmt.Println(blue("make buf-generate"))
	fmt.Println(blue("go mod tidy"))
	fmt.Println(blue("go mod vendor"))
	fmt.Println()

	if updatedArgo || updatedStagingArgo {
		fmt.Println("It is recommended you commit and push at this point, then run the following:")
		fmt.Println()
		if updatedArgo {
			fmt.Println(blue("kubectl apply -f deploy/argocd/production.yaml"))
		}
		if updatedStagingArgo {
			fmt.Println(blue("kubectl apply -f deploy/argocd/staging.yaml"))
		}
		fmt.Println()
	} else {
		fmt.Println("It is recommended you commit and push at this point.")
		fmt.Println()
	}
}

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
	Name     string            `yaml:"name"`
	Team     string            `yaml:"team"`
	Domain   string            `yaml:"domain"`
	Services []MetadataService `yaml:"services"`
}

type MetadataService struct {
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	CIEnabled bool   `yaml:"ciEnabled"`
}
