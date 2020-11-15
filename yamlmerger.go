package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"clickandboat.com/simpleyaml"
)

// TODO: set optional flags for "i", "dpl", "o", "c"
// TODO: set usage like in https://github.com/golang/lint/blob/master/golint/golint.go

var (
	inputFlag         = flag.String("i", "", "Input YAML files. e.g: \"file1.yaml file2.yaml [...]\"")
	outputFlag        = flag.String("o", "none", "Output YAML file")
	checkFlag         = flag.Bool("c", false, "Checks if the merge passes (without writing to a file)")
	deletionTokenFlag = flag.String("del-tk", "nil", "Deletion token to identify which node(s) to delete")
	delimPerListFlag  = flag.String("dpl", "none", "Delimiter per list to identify key and value. e.g: \"keyname1:delim1,keyname2:delim2[,...]\"")
)

func main() {
	flag.Parse()

	var inputFiles []*os.File
	processInputFlag(&inputFiles)
	for i := 0; i < len(inputFiles); i++ {
		defer inputFiles[i].Close()
	}

	outputFilePath := strings.TrimSpace(*outputFlag)

	_, statErr := os.Stat(outputFilePath)
	if !os.IsNotExist(statErr) {
		fmt.Println("Output file already exists!")
		os.Exit(1)
	}

	deletionToken := strings.TrimSpace(*deletionTokenFlag)

	delimiterPerList := strings.TrimSpace(*delimPerListFlag)
	if delimiterPerList == "none" {
		delimiterPerList = ""
	}

	var yamls []*simpleyaml.YamlNode

	parseErr := parseYamls(inputFiles, &yamls)
	if parseErr != nil {
		fmt.Println(parseErr)
		os.Exit(1)
	}

	merger := simpleyaml.NewMerger(yamls, deletionToken, delimiterPerList, true)

	mergedYaml, mergeErr := merger.Merge()
	if mergeErr != nil {
		fmt.Println(mergeErr)
		os.Exit(1)
	}

	if *checkFlag {
		fmt.Println("Merge successful.")
		//os.Exit(0)
		return
	}

	outputFile, createErr := os.Create(outputFilePath)
	if createErr != nil {
		fmt.Println(createErr)
		os.Exit(1)
	}
	defer outputFile.Close()

	writer := simpleyaml.NewWriter(outputFile)
	writer.Write(mergedYaml)

	fmt.Println("Merge successful.")
}

func processInputFlag(inputFiles *[]*os.File) {
	inputFilesPaths := strings.Split(*inputFlag, " ")

	c := len(inputFilesPaths)
	*inputFiles = make([]*os.File, c)

	for i := 0; i < c; i++ {
		filePath := strings.TrimSpace(inputFilesPaths[i])

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		(*inputFiles)[i] = file
	}
}

func parseYamls(files []*os.File, yamls *[]*simpleyaml.YamlNode) error {
	c := len(files)

	for i := 0; i < c; i++ {
		parser := simpleyaml.NewParser(files[i])

		yaml, err := parser.Parse()
		if err != nil {
			return err
		}

		*yamls = append(*yamls, yaml)
	}

	return nil
}
