package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	inFileFlag = flag.String("in", "", "Input file")
	outDirFlag = flag.String("out", "", "Output directory")
)

func printFlagUsageAndExit() {
	flag.Usage()
	os.Exit(2)
}

type generatedSingleEntityFiles struct {
	SchemaCreate []byte
	Entity       []byte
	Iterator     []byte
	Repository   []byte
}

func generateSingleEntityFiles(entity *GeneratorEntity, packageName string) (generated *generatedSingleEntityFiles, err error) {
	defer handleDeferAndSetError(&err)

	generated = &generatedSingleEntityFiles{
		SchemaCreate: NewAppender().AppendSchemaCreate(entity).Bytes(packageName),
		Entity:       NewAppender().AppendEntityStructs(entity).Bytes(packageName),
		Iterator:     NewAppender().AppendEntityIterators(entity).Bytes(packageName),
		Repository:   NewAppender().AppendRepoInterface(entity).Bytes(packageName),
	}
	err = nil
	return
}

type generatedMultipleEntityFiles struct {
	RepositoryFactory []byte
}

func generateMultipleEntityFiles(generatorSetup *GeneratorSetup, packageName string) (generated *generatedMultipleEntityFiles, err error) {
	defer handleDeferAndSetError(&err)

	generated = &generatedMultipleEntityFiles{
		RepositoryFactory: NewAppender().AppendRepositoryFactories(generatorSetup).Bytes(packageName),
	}
	err = nil
	return
}

func main() {
	flag.Parse()

	if len(*inFileFlag) == 0 ||
		len(*outDirFlag) == 0 {

		printFlagUsageAndExit()
	}

	generatorSetup, err := LoadGeneratorSetup(*inFileFlag)
	if err != nil {
		log.Fatal(err)
	}

	outFileDirNameOnly := filepath.Base(*outDirFlag)
	packageName := outFileDirNameOnly

	type fileToWrite struct {
		FilePath string
		Content  []byte
	}
	var filesToWrite []*fileToWrite

	for _, entity := range generatorSetup.Entities {
		generated, err := generateSingleEntityFiles(entity, packageName)
		if err != nil {
			log.Fatal(err)
		}

		filesToWrite = append(filesToWrite,
			&fileToWrite{FilePath: filepath.Join(*outDirFlag, entity.EntityName+"_schema_create.go"), Content: generated.SchemaCreate},
			&fileToWrite{FilePath: filepath.Join(*outDirFlag, entity.EntityName+"_entity.go"), Content: generated.Entity},
			&fileToWrite{FilePath: filepath.Join(*outDirFlag, entity.EntityName+"_iterator.go"), Content: generated.Iterator},
			&fileToWrite{FilePath: filepath.Join(*outDirFlag, entity.EntityName+"_repository.go"), Content: generated.Repository},
		)
	}

	generatedMultiEntity, err := generateMultipleEntityFiles(generatorSetup, packageName)
	if err != nil {
		log.Fatal(err)
	}
	filesToWrite = append(filesToWrite,
		&fileToWrite{FilePath: filepath.Join(*outDirFlag, "repository_factory.go"), Content: generatedMultiEntity.RepositoryFactory},
	)

	for _, f := range filesToWrite {
		err = ioutil.WriteFile(f.FilePath, f.Content, 0655)
		if err != nil {
			log.Fatal(err)
		}
	}
}
