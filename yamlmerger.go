package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"clickandboat.com/simpleyaml"
)

var (
	inputFlag         = flag.String("i", "", "Input YAML files. e.g: \"file1.yaml file2.yaml [...]\"")
	outputFlag        = flag.String("o", "", "[optional] Output YAML file")
	deletionTokenFlag = flag.String("del-tk", "", "[optional] Deletion token to identify which node(s) to delete")
	delimPerListFlag  = flag.String("dpl", "", "[optional] Delimiter per list to identify key and value. e.g: \"keyname1:delim1,keyname2:delim2[,...]\"")
	outForceFlag      = flag.Bool("of", false, "[optional] Overwrite output file if exists")
)

func main() {
	flag.Parse()

	if *inputFlag == "" {
		flag.Usage()
		os.Exit(2)
	}

	var inputFiles []*os.File
	processInputFlag(&inputFiles)
	for i := 0; i < len(inputFiles); i++ {
		defer inputFiles[i].Close()
	}

	deletionToken := strings.TrimSpace(*deletionTokenFlag)

	delimiterPerList := strings.TrimSpace(*delimPerListFlag)
	delimPerListMap := simpleyaml.RawDelimPerListToMap(delimiterPerList)

	var yamls []*simpleyaml.YamlNode

	parseErr := parseYamls(inputFiles, &yamls)
	if parseErr != nil {
		fmt.Println(parseErr)
		os.Exit(1)
	}

	merger := simpleyaml.NewMerger(yamls, deletionToken, delimPerListMap, true)

	mergedYaml, mergeErr := merger.Merge()
	if mergeErr != nil {
		fmt.Println(mergeErr)
		os.Exit(1)
	}

	if *outputFlag != "" {
		writeMergedFile(mergedYaml)
	}

	fmt.Println("Merge successful.")
}

func writeMergedFile(mergedYaml *simpleyaml.YamlNode) {
	outputFilePath := strings.TrimSpace(*outputFlag)

	if !*outForceFlag {
		_, statErr := os.Stat(outputFilePath)
		if !os.IsNotExist(statErr) {
			fmt.Println("Output file already exists!")
			os.Exit(1)
		}
	}

	outputFile, createErr := os.Create(outputFilePath)
	if createErr != nil {
		fmt.Println(createErr)
		os.Exit(1)
	}
	defer outputFile.Close()

	writer := simpleyaml.NewWriter(outputFile)
	writer.Write(mergedYaml)
}

func processInputFlag(inputFiles *[]*os.File) {
	inputFilesPaths := strings.Split(*inputFlag, " ")

	if len(inputFilesPaths) < 2 {
		fmt.Println("You must specify at least 2 input files")
		os.Exit(2)
	}

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
