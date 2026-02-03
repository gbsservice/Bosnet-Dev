package jobs

import (
	"api_kino/config/app"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
)

var (
	client = resty.New()
	prefix = "http://127.0.0.1:" + app.Config().Port + "/v1"
)

type JobConfig struct {
	ClientID     string `json:"client_id"`
	ScheduleTime string `json:"schedule_time"`
}

func runJob(jobConfig JobConfig) {
	url_draftso := "/post-draftso"
	params := map[string]interface{}{
		"client_id": jobConfig.ClientID,
	}

	jobName := fmt.Sprintf("Running Job Client ID: %s at", jobConfig.ClientID)
	sch := gocron.NewScheduler(time.Local)
	_, _ = sch.Every(1).Day().At(jobConfig.ScheduleTime).Do(func() {
		fmt.Println(jobName, time.Now())
		url := prefix + url_draftso
		_, err := client.R().EnableTrace().SetBody(params).Post(url)
		if err != nil {
			fmt.Println("Error during HTTP request:", err)
		}
	})
	sch.StartAsync()
}

func HandleJobs() {
	jobConfigsEnv := `[{"client_id": "35","schedule_time": "08:05"},{"client_id": "26","schedule_time": "08:15"},{"client_id": "41","schedule_time": "08:25"},{"client_id": "42","schedule_time": "08:35"},{"client_id": "35","schedule_time": "09:05"},{"client_id": "26","schedule_time": "09:15"},{"client_id": "41","schedule_time": "09:25"},{"client_id": "42","schedule_time": "09:35"},{"client_id": "35","schedule_time": "10:05"},{"client_id": "26","schedule_time": "10:15"},{"client_id": "41","schedule_time": "10:25"},{"client_id": "42","schedule_time": "10:35"},{"client_id": "35","schedule_time": "11:05"},{"client_id": "26","schedule_time": "11:15"},{"client_id": "41","schedule_time": "11:25"},{"client_id": "42","schedule_time": "11:35"},{"client_id": "35","schedule_time": "12:05"},{"client_id": "26","schedule_time": "12:15"},{"client_id": "41","schedule_time": "12:25"},{"client_id": "42","schedule_time": "12:35"},{"client_id": "35","schedule_time": "13:05"},{"client_id": "26","schedule_time": "13:15"},{"client_id": "41","schedule_time": "13:25"},{"client_id": "42","schedule_time": "13:35"},{"client_id": "35","schedule_time": "14:05"},{"client_id": "26","schedule_time": "14:15"},{"client_id": "41","schedule_time": "14:25"},{"client_id": "42","schedule_time": "14:35"},{"client_id": "35","schedule_time": "15:05"},{"client_id": "26","schedule_time": "15:15"},{"client_id": "41","schedule_time": "15:25"},{"client_id": "42","schedule_time": "15:35"},{"client_id": "35","schedule_time": "16:05"},{"client_id": "26","schedule_time": "16:15"},{"client_id": "41","schedule_time": "16:25"},{"client_id": "42","schedule_time": "16:35"},{"client_id": "35","schedule_time": "17:05"},{"client_id": "26","schedule_time": "17:15"},{"client_id": "41","schedule_time": "17:25"},{"client_id": "42","schedule_time": "17:35"},{"client_id": "35","schedule_time": "18:05"},{"client_id": "26","schedule_time": "18:15"},{"client_id": "41","schedule_time": "18:25"},{"client_id": "42","schedule_time": "18:35"},{"client_id": "35","schedule_time": "19:05"},{"client_id": "26","schedule_time": "19:15"},{"client_id": "41","schedule_time": "19:25"},{"client_id": "42","schedule_time": "19:35"},{"client_id": "35","schedule_time": "20:05"},{"client_id": "26","schedule_time": "20:15"},{"client_id": "41","schedule_time": "20:25"},{"client_id": "42","schedule_time": "20:35"},{"client_id": "35","schedule_time": "21:05"},{"client_id": "26","schedule_time": "21:15"},{"client_id": "41","schedule_time": "21:25"},{"client_id": "42","schedule_time": "21:35"},{"client_id": "93", "schedule_time": "08:45"},{"client_id": "93", "schedule_time": "09:45"},{"client_id": "93", "schedule_time": "10:45"},{"client_id": "93", "schedule_time": "11:45"},{"client_id": "93", "schedule_time": "12:45"},{"client_id": "93", "schedule_time": "13:45"},{"client_id": "93", "schedule_time": "14:45"},{"client_id": "93", "schedule_time": "15:45"},{"client_id": "93", "schedule_time": "16:45"},{"client_id": "93", "schedule_time": "17:45"},{"client_id": "93", "schedule_time": "18:45"},{"client_id": "93", "schedule_time": "19:45"},{"client_id": "93", "schedule_time": "20:45"},{"client_id": "93", "schedule_time": "21:45"}]`
	if jobConfigsEnv == "" {
		fmt.Println("No job configurations found in the environment.")
		return
	}

	var jobConfigs []JobConfig
	if err := json.Unmarshal([]byte(jobConfigsEnv), &jobConfigs); err != nil {
		fmt.Println("Error parsing job configurations:", err)
		return
	}

	for _, jobConfig := range jobConfigs {
		runJob(jobConfig)
	}
}
