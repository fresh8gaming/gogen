package cmd

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gitlab.sportradar.ag/ads/adsstaff/gogen/internal/util"

	"github.com/spf13/cobra"
)

const RepoTemplates string = "_template/repo"

type Repo struct {
	Name   string
	Team   string
	Domain string
	Org    string
}

//go:embed _template/repo/*
//go:embed _template/repo/deploy/chart/templates/_helpers.tpl
//go:embed _template/repo/.ci/blacklist.txt
//go:embed _template/repo/.ci/config.yml
var repoContent embed.FS

var (
	Name   string
	Domain string
	Team   string
	Org    string
)

func GetRepoCmd() (*cobra.Command, error) {
	var cmdRepo = &cobra.Command{
		Use:   "repo",
		Short: "",
		Long:  "",
		Run:   runRepoCmd(),
	}

	var err error

	cmdRepo.Flags().StringVarP(&Name, "name", "n", "", "Name of the monorepo (default to name of target directory)")

	cmdRepo.Flags().StringVarP(&Domain, "domain", "d", "", "Domain the monorepo relates to, also used for K8s namespace (required)")
	err = cmdRepo.MarkFlagRequired("domain")
	if err != nil {
		return nil, err
	}

	cmdRepo.Flags().StringVarP(&Team, "team", "t", "", "Team responsible for the monorepo (required), "+
		"must be one of adspp|adssd|adssocial|adssearch|adsstaff")
	err = cmdRepo.MarkFlagRequired("team")
	if err != nil {
		return nil, err
	}

	return cmdRepo, nil
}

func runRepoCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
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

		repo := Repo{
			Name:   getName(Name, absPath),
			Domain: Domain,
			Team:   Team,
		}

		fmt.Printf("Creating %s at %s\n", blue(repo.Name), absPath)

		copyTemplates(absPath, RepoTemplates, repo, repoContent, RepoTemplates, func(path string) string { return path })

		fmt.Printf("Created %s!\n", green(repo.Name))
		fmt.Println()
		fmt.Println("Run the following commands to set up the project:")
		fmt.Println()
		fmt.Printf(blue("cd %s\n"), absPath)
		modPath := fmt.Sprintf("gitlab.sportradar.ag/ads/%s/%s", repo.Team, repo.Name)
		fmt.Printf(blue("go mod init %s\n"), modPath)
		fmt.Println(blue("go mod tidy"))
		fmt.Println(blue("go mod vendor"))
		fmt.Println(blue("make install-tools"))
		fmt.Println()
	}
}

func copyTemplates(
	targetDir, dir string,
	templateData interface{},
	content embed.FS,
	templatePath string,
	pathManipulation func(string) string,
) {
	testDir, err := content.ReadDir(dir)
	util.Fatal(err)

	for _, testDirContent := range testDir {
		pathing := filepath.Join(dir, testDirContent.Name())
		localPathing, err := filepath.Abs(filepath.Join(targetDir, strings.TrimPrefix(pathing, templatePath)))
		util.Fatal(err)

		if testDirContent.Type().IsDir() {
			err := os.MkdirAll(pathManipulation(localPathing), os.ModePerm)
			util.Fatal(err)

			copyTemplates(targetDir, pathing, templateData, content, templatePath, pathManipulation)
		} else {
			f, err := os.Create(pathManipulation(localPathing))
			util.Fatal(err)

			byteBuf := bytes.NewBuffer([]byte{})

			t, err := template.ParseFS(content, pathing)
			util.Fatal(err)

			w := bufio.NewWriter(byteBuf)
			err = t.Execute(w, templateData)
			util.Fatal(err)

			err = w.Flush()
			util.Fatal(err)

			var fmtByte []byte
			if filepath.Ext(localPathing) == "go" {
				fmtByte, err = format.Source(byteBuf.Bytes())
				util.Fatal(err)
			} else {
				fmtByte = byteBuf.Bytes()
			}

			_, err = f.Write(fmtByte)
			util.Fatal(err)
		}
	}
}

func getName(nameFlag, absPath string) string {
	if nameFlag != "" {
		return nameFlag
	}

	return filepath.Base(absPath)
}
