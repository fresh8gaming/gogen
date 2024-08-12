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

const HTTPGRPCTemplates string = "_template/service/http-grpc"

//go:embed _template/service/http-grpc/*
var httpGRPCContent embed.FS

type Service struct {
	Name                  string
	Team                  string
	ServiceName           string
	ServiceNameUnderscore string
	ServiceNameProto      string
	ServiceInplay         string
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

	service := &Service{
		Name:                  getName(Name, absPath),
		Team:                  Team,
		ServiceName:           ServiceName,
		ServiceNameUnderscore: strings.ReplaceAll(ServiceName, "-", "_"),
		ServiceNameProto:      strings.ReplaceAll(strings.ReplaceAll(ServiceName, "_", ""), "-", ""),
	}

	fmt.Printf("Creating %s in %s\n", blue(service.ServiceName), absPath)

	updateMetadata(absPath, service, "http-grpc")

	copyTemplates(absPath, HTTPGRPCTemplates, service, httpGRPCContent, HTTPGRPCTemplates, func(path string) string {
		replaced := strings.ReplaceAll(path, "service-name", service.ServiceName)
		replaced = strings.ReplaceAll(replaced, "service_name", service.ServiceNameUnderscore)
		replaced = strings.ReplaceAll(replaced, "servicename", service.ServiceNameProto)
		return replaced
	})

	fmt.Printf("Created %s!\n", green(service.ServiceName))
	fmt.Println()
	fmt.Println("Run the following commands to set up the project:")
	fmt.Println()
	fmt.Printf(blue("cd %s\n"), absPath)
	fmt.Println(blue("go fmt ./..."))
	fmt.Println(blue("go mod tidy"))
	fmt.Println(blue("go mod vendor"))
	fmt.Println()

	fmt.Println("It is recommended you commit and push at this point.")
	fmt.Println()
}
