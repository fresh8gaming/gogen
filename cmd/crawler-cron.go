package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fresh8gaming/gogen/internal/util"
)

const CrawlerCronTemplates string = "_template/service/crawler-cron"

//go:embed _template/service/crawler-cron/*
var crawlerCronContent embed.FS

var (
	isInplay bool
)

func runCrawlerCronCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		createCrawlerCronService(args)
	}
}

func createCrawlerCronService(args []string) {
	if len(args) < 1 {
		fmt.Println("insifficient arguments, requires at least name being set")
		return
	}
	target := args[0]
	if len(args) > 1 && args[1] == "true" {
		isInplay = true
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

	service := &Service{
		Name:                  getName(Name, absPath),
		Org:                   Org,
		ServiceName:           ServiceName,
		ServiceNameUnderscore: strings.ReplaceAll(ServiceName, "-", "_"),
		ServiceInplay:         strconv.FormatBool(isInplay),
	}

	fmt.Printf("Creating %s in %s\n", blue(service.ServiceName), absPath)

	updateMetadata(absPath, service, "cron")

	copyTemplates(absPath, CrawlerCronTemplates, service, crawlerCronContent, CrawlerCronTemplates, func(path string) string {
		replaced := strings.ReplaceAll(path, "service-name", service.ServiceName)
		replaced = strings.ReplaceAll(replaced, "service_name", service.ServiceNameUnderscore)
		return replaced
	})

	updatedArgo := updateArgoApplication(argoApplicationFilePath, "values-production", service)

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
		if updatedArgo {
			fmt.Println(blue("kubectl apply -f deploy/argocd/production.yaml"))
		}
		fmt.Println()
	} else {
		fmt.Println("It is recommended you commit and push at this point.")
		fmt.Println()
	}
}
