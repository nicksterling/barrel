/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")

		//First let's pull off the command line flags we need to consume...

		/**
		 * Where do we want to store the output of the docker images
		 */
		outputDir, _ := cmd.Flags().GetString("output-dir")
		if outputDir == "" {
			outputDir = "."
		}
		fmt.Println("Output Directory: " + outputDir)

		// Are we integrating directly to Helm
		helm, _ := cmd.Flags().GetString("helm")

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "World"
		}

		/**
		 *   If we're not using HELM then we're using STDIN...
		 */
		if helm == "" {
			// f, _ := os.Open("harbor.txt")
			var line string
			scanner := bufio.NewScanner(os.Stdin)
			// scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line = trim(scanner.Text())
				if line != "" {

					pullImage(line)
					saveImage(line, getImageName(line), getImageTag(line), "images")
				}
			}

		} else {
			// Let's grab the manifest of the helm chart and save it to a temp file

			tmpFile, err := ioutil.TempFile(os.TempDir(), "barrel-")
			if err != nil {
				fmt.Println("Cannot create temporary file")
				log.Fatal("Cannot create temporary file", err)

			}

			// Remember to clean up the file afterwards
			defer os.Remove(tmpFile.Name())

			// fmt.Println("Created File: " + tmpFile.Name())

			// Example writing to the file
			cmd := exec.Command("helm", "template", helm)
			stdout, err := cmd.Output()

			if err != nil {
				fmt.Println(err.Error())
			}

			text := []byte(string(stdout))
			if _, err = tmpFile.Write(text); err != nil {
				fmt.Println("Failed to write to temporary file")
				log.Fatal("Failed to write to temporary file", err)

			}

			//Ensure the output directory exists
			createDirIfNotExist(outputDir)

			saveHelmChart(helm, outputDir)

			f, _ := os.Open(tmpFile.Name())
			var line string
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line = trim(scanner.Text())
				if line != "" {
					pullImage(line)
					saveImage(line, getImageName(line), getImageTag(line), outputDir)
				}
			}

			// Close the file
			if err := tmpFile.Close(); err != nil {
				log.Fatal(err)
			}

		}

	},
}

func saveHelmChart(chart string, outputDir string) {
	cmd := exec.Command("helm", "pull", chart, "-d", outputDir)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(stdout))
}

func createDirIfNotExist(dir string) {

	// if _, err := os.Stat(dir); os.IsNotExist(err) {
	// 	os.Mkdir(dir, os.ModeDir)
	// }

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func saveImage(line string, imageName string, imageTag string, directory string) {
	fmt.Println("-----------------------------------------------------")
	fmt.Println("--  SAVING IMAGE")
	fmt.Println("--  IMAGE: " + getImageName(line))
	fmt.Println("--    TAG: " + getImageTag(line))
	fmt.Println("--   PATH: " + directory + "/" + imageName + "-" + imageTag + ".tar")
	fmt.Println("-----------------------------------------------------")

	cmd := exec.Command("docker", "save", "--output", directory+"/"+imageName+"-"+imageTag+".tar", line)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(stdout))
}

func getImageName(input string) string {
	var splits = strings.Split(input, "/")

	var result = strings.Split(splits[len(splits)-1], ":")
	return result[0]
}

func getImageTag(input string) string {
	var splits = strings.Split(input, "/")

	var result = strings.Split(splits[len(splits)-1], ":")
	return result[1]
}

func pullImage(line string) {
	fmt.Println("-----------------------------------------------------")
	fmt.Println("--  PULLING IMAGE")
	fmt.Println("--  IMAGE: " + getImageName(line))
	fmt.Println("--    TAG: " + getImageTag(line))
	fmt.Println("-----------------------------------------------------")

	cmd := exec.Command("docker", "pull", line)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(stdout))
}

/**
 *
 */
func trim(input string) string {
	var splits = strings.Split(input, "image: ")
	if len(splits) < 2 {
		return ""
	}

	return strings.Replace(splits[1], "\"", "", -1)
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	pullCmd.PersistentFlags().String("name", "", "A help for foo")
	pullCmd.PersistentFlags().String("output-dir", "", "A help for foo")
	pullCmd.PersistentFlags().String("helm", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
