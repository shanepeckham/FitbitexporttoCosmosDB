// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"time"
	"github.com/emicklei/go-restful/log"
	"os"
)

var tokenMount string
var batch string

// dailyexportCmd represents the dailyexport command
var dailyexportCmd = &cobra.Command{
	Use:   "dailyexport",
	Short: "This will extract your sleep, activity and heartrate data",
	Long: `--tokenMount is where you wil store your OAUTH token e.g. --tokenMount=token
           --date is the exportdate for which you want to extract data for a single date e.g. --date="2017-11-01"
           --batch is used top extract your data from a point in time, if true then the date flag becomes the date from which date will be extracted`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dailyexport called")

		c := oauthDance(tokenMount)
		refreshToken(c)

		if batch == "true" {
			genesisDate, _ := time.Parse("2006-01-02", exportDate)
			currentDate := time.Now().AddDate(0,0,-1)
			currentDateMinusOne := currentDate.Format("2006-01-02")

			for exportDate != currentDateMinusOne {

				log.Print("Processing", exportDate)

				activities := getDailyActivitySummary(c, exportDate)
				WriteActivitiesToCosmosDB(activities)
				heartrate := getDailyHeartRateSummary(c, exportDate)
				WriteHeartrateToCosmosDB(heartrate)
				sleep := getDailySleepSummary(c, exportDate)
				WriteSleepToCosmosDB(sleep)

				genesisDate = genesisDate.AddDate(0,0,1)
				exportDate = genesisDate.Format("2006-01-02")

				if currentDate == genesisDate {
					fmt.Print("Halt")
					os.Exit(1)
				}
			}



		} else {
			activities := getDailyActivitySummary(c, exportDate)
			WriteActivitiesToCosmosDB(activities)
			heartrate := getDailyHeartRateSummary(c, exportDate)
			WriteHeartrateToCosmosDB(heartrate)
			sleep := getDailySleepSummary(c, exportDate)
			WriteSleepToCosmosDB(sleep)
		}



	},
}

func init() {
	rootCmd.AddCommand(dailyexportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	 dailyexportCmd.Flags().StringVarP(&exportDate, "date", "d", "", "Date to run format for yyyy-mm-dd")
	 dailyexportCmd.Flags().StringVarP(&tokenMount,"tokenMount", "t", "","Oauth Token mount for persistent tokens")
	 dailyexportCmd.Flags().StringVarP(&batch,"batch", "b", "","Process in batch from export date to current")
	 cachePath = tokenMount
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dailyexportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
