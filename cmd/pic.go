package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const PIC_LOCATION = "https://icon.scie.com.cn/user/sicon/s%s_o.jpg"

func downloadFile(URL, fileName string) error {

	log.Info("Downloading file from: ", URL)
	// add referer header to custom request
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	req.Header["Referer"] = []string{"https://www.alevel.com.cn/"}
	// do request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Error("Error downloading file: ", response.Status)
		return errors.New("received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Trace("Writing file")

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getStudentPic(sid string) error {
	url := fmt.Sprintf(PIC_LOCATION, sid)
	log.Info("Downloading pic from: ", url)
	err := downloadFile(url, sid+".jpg")
	if err != nil {
		return err
	}
	return nil
}

var picCmd = &cobra.Command{
	Use:   "pic [sid]",
	Short: "Get the corresponding profile pic of a student",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sid := args[0]
		err := getStudentPic(sid)
		if err != nil {
			log.Error(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(picCmd)
}
