package main

import (
    "fmt"
    "os"
    "os/exec"
    "bufio"
    "strings"
)

func getImageName(line string) string{
	var splits = strings.Split(line, "/")

    var image =  strings.Split(splits[len(splits)-1], ":") 
	return image[0]

}
func getImageTag(line string) string{
	var splits = strings.Split(line, "/")

    var image =  strings.Split(splits[len(splits)-1], ":") 
	return image[1]

}


func saveImage(line string, imageName string, imageTag string, directory string){

	cmd := exec.Command("docker", "save", "--output", directory + "/" + imageName + "-" + imageTag + ".tar", line)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(stdout))
}

func tagImage(line string){

	var newImage = "nicksterling/" + line

	cmd := exec.Command("docker", "tag", line, newImage)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(stdout))
}

func pullImage(line string){
	cmd := exec.Command("docker", "pull", line)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(stdout))
}

func trim(line string) string{
	var splits = strings.Split(line, "image: ")
	if len(splits) < 2{
		return ""
	} else {
		return splits[1]
	}
}

func main() {
    var line string
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
		line = trim(scanner.Text())

		if line != "" {
			fmt.Println(line)
			// fmt.Println(getImageName(line))
	
			pullImage(line)
			// tagImage(line)
			saveImage(line, getImageName(line), getImageTag(line), "images")
		}

    }
}
