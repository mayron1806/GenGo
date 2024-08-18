/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mayron1806/gengo/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dir string

// initCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init gengo in current directory",
	Long:  `Init gengo in current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		// Diretório .gengo
		dirPath := filepath.Join(dir, ".gengo")
		fmt.Printf("Init gengo in %s...\n", dirPath)
		// Verifica se o diretório já existe
		if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
			// verifica se gengo.json existe e cria se não existir
			if !gengoJsonExists(dirPath) {
				fmt.Println("gengo.json not found, creating it...")
				createGengoJson(dirPath)

			}
			// verifica se templates existe e cria se não existir
			if !templatesDirExists(dirPath) {
				fmt.Println("templates directory not found, creating it...")
				createTemplatesDir(dirPath)
			}
			return
		}

		// Cria o diretório .gengo
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			fmt.Printf("Error creating .gengo directory: %v\n", err)
			return
		}

		// Cria o arquivo gengo.json
		gengoJsonCreated := createGengoJson(dirPath)

		// Cria o diretório templates
		templatesDirCreated := createTemplatesDir(dirPath)

		if !gengoJsonCreated || !templatesDirCreated {
			return
		}
		fmt.Println(".gengo directory initialized successfully")
	},
}

func init() {
	InitCmd.Flags().StringVarP(&dir, "dir", "d", "./", "Directory to create .gengo in")
	viper.BindPFlag("dir", InitCmd.Flags().Lookup("dir"))
}
func createGengoJson(path string) bool {
	fmt.Println("Creating gengo.json...")
	jsonFilePath := filepath.Join(path, "gengo.json")
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		fmt.Printf("Error creating gengo.json: %v\n", err)
		return false
	}
	defer jsonFile.Close()
	jsonFile.WriteString(constants.GetDefaultGengoJson())
	// add content to gengo.json
	fmt.Println("gengo.json created successfully")
	return true
}
func createTemplatesDir(path string) bool {
	fmt.Println("Creating templates directory...")
	templatesDirPath := filepath.Join(path, "templates")
	err := os.Mkdir(templatesDirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating templates directory: %v\n", err)
		return false
	}
	fmt.Println("templates directory created successfully")
	return true
}
func gengoJsonExists(path string) bool {
	jsonFilePath := filepath.Join(path, "gengo.json")
	_, err := os.Stat(jsonFilePath)
	return !os.IsNotExist(err)
}
func templatesDirExists(path string) bool {
	templatesDirPath := filepath.Join(path, "templates")
	_, err := os.Stat(templatesDirPath)
	return !os.IsNotExist(err)
}
