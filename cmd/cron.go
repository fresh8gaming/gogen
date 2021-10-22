package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fresh8gaming/gogen/internal/util"

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

	argoApplicationFilePath := filepath.Join(absPath, "deploy", "argocd", "application.yaml")
	if _, err := os.Stat(argoApplicationFilePath); os.IsNotExist(err) {
		log.Fatalf("argocd application file expected at %s", argoApplicationFilePath)
	}

	service := Service{
		Name:                  getName(Name, absPath),
		Org:                   Org,
		ServiceName:           ServiceName,
		ServiceNameUnderscore: strings.ReplaceAll(ServiceName, "-", "_"),
	}

	fmt.Printf("Creating %s in %s\n", blue(service.ServiceName), absPath)

	updateMetadata(absPath, service, "cron")

	copyTemplates(absPath, CronTemplates, service, cronContent, CronTemplates, func(path string) string {
		return strings.ReplaceAll(path, "service-name", service.ServiceName)
	})

	updatedArgo := updateArgoApplication(argoApplicationFilePath, service)

	fmt.Printf("Created %s!\n", green(service.ServiceName))
	fmt.Println()
	fmt.Println("Run the following commands to set up the project:")
	fmt.Println()
	fmt.Printf(blue("cd %s\n"), absPath)
	fmt.Println(blue("go fmt ./..."))
	fmt.Println(blue("go mod tidy"))
	fmt.Println(blue("go mod vendor"))
	fmt.Println()

	if updatedArgo {
		fmt.Println("It is recommended you commit and push at this point, then run the following:")
		fmt.Println()
		fmt.Println(blue("kubectl apply -f deploy/argocd/application.yaml"))
		fmt.Println()
	} else {
		fmt.Println("It is recommended you commit and push at this point.")
		fmt.Println()
	}
}
