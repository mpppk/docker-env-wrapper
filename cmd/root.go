package cmd

import (
	"fmt"
	"os"

	"io/ioutil"

	"strings"

	"github.com/mpppk/docker-env/env"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var queryFlag string
var formatFlag string

const FORMAT_fLAG_DOCKER_FILE = "dockerfile"
const FORMAT_fLAG_DOCKER_COMPOSE = "compose"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "docker-env IMAGE_NAME",
	Short: "Generate Dockerfile or docker-compose.yml with host environment variables",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Image name must be specified")
			os.Exit(1)
		}

		env, err := env.New().Filter(queryFlag)
		if err != nil {
			panic(err)
		}

		imageName := args[0]

		switch formatFlag {
		case FORMAT_fLAG_DOCKER_FILE:
			out := fmt.Sprintln("FROM " + imageName)
			for k, v := range env {
				out += fmt.Sprintf("ENV %v %v\n", k, v)
			}

			fmt.Println(out)
			ioutil.WriteFile("Dockerfile", []byte(out), 0777)
		case FORMAT_fLAG_DOCKER_COMPOSE:
			containerKey := strings.Replace(imageName, ":", "", -1)

			d := map[string]interface{}{
				"version": "3",
				"services": map[string]interface{}{
					containerKey: map[string]interface{}{
						"image":       imageName,
						"environment": env,
					},
				},
			}
			o, err := yaml.Marshal(d)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(o))
			ioutil.WriteFile("docker-compose.yml", o, 0777)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.Flags().StringVarP(&queryFlag, "query", "q", ".*", "Filter host environment variables with regular expressions.")
	RootCmd.Flags().StringVarP(&formatFlag, "format", "f", FORMAT_fLAG_DOCKER_FILE, "Specify output format. [dockerfile|compose]")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
