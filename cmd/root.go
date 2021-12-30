package cmd

import (
	"os"

	"encoding/json"
	"net/http"

	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const STUDENT_QUERY = "https://api.sciers.alanjin.me:2334/stu_query/"

type Student struct {
	Psid            string
	Sh_house        string
	Student_pingyin string
	Student_num     int
	Sh_grade        string
	Student_ename   string
	Myclass         string
	Student_name    string
}

func studentInfo(query string) ([]Student, error) {
	url := STUDENT_QUERY + query
	log.Info("Querying: ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return []Student{}, err
	}
	defer resp.Body.Close()

	log.Trace("Parsing into json")
	var student []Student
	err = json.NewDecoder(resp.Body).Decode(&student)
	if err != nil {
		log.Error(err)
		return []Student{}, err
	}
	return student, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stu [query]",
	Short: "SCIE student info lookup",
	Args:  cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set log according to Verbose
		if Verbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetOutput(ioutil.Discard)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		query := ""
		for _, arg := range args {
			query += arg + " "
		}
		query = query[:len(query)-1]
		if query == "" {
			fmt.Println("No query specified, use `stu -h` for help")
			return
		}
		log.Info("Query: ", query)
		students, err := studentInfo(query)
		if err != nil {
			log.Error(err)
			fmt.Println("Error: ", err)
			return
		}
		if Json {
			log.Info("Output in json format")
			b, err := json.MarshalIndent(students, "", "    ")
			if err != nil {
				log.Error(err)
				return
			}
			fmt.Println(string(b))
			return
		} else {
			log.Info("Output in text format")
			for _, student := range students {
				fmt.Println(student.Psid)
				fmt.Println("  - " + student.Student_ename)
				fmt.Println("  - " + student.Student_name)
				fmt.Println("  - " + student.Myclass)
				fmt.Println("  - " + student.Sh_house)
			}
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var Verbose bool
var Json bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolVarP(&Json, "json", "j", false, "output in json format")
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
