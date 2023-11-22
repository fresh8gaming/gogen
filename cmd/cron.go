package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.sportradar.ag/ads/adsstaff/gogen/internal/util"

	"github.com/spf13/cobra"
)

const CronTemplates string = "_template/service/cron"

//go:embed _template/service/cron/*
var cronContent embed.FS

func runCronCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		createCronService(args)
	}
}

func createCronService(args []string) {
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

	argoApplicationFilePath := filepath.Join(absPath, "deploy", "argocd", "production.yaml")
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
	}

	fmt.Printf("Creating %s in %s\n", blue(service.ServiceName), absPath)

	updateMetadata(absPath, service, "cron")

	copyTemplates(absPath, CronTemplates, service, cronContent, CronTemplates, func(path string) string {
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
