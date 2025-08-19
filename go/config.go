/*
Copyright 2023 The Vitess Authors.

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

package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type config struct {
	Github githubapp.Config

	botLogin        string
	reviewChecklist string
	address         string
	logFile         string
}

func readConfig() (*config, error) {
	// Only load .env if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	var c config
	c.Github.SetValuesFromEnv("")

	// Ensure required GitHub App env vars are set
	integrationID := os.Getenv("GITHUB_APP_INTEGRATION_ID")
	webhookSecret := os.Getenv("GITHUB_APP_WEBHOOK_SECRET")
	pathPrivateKey := os.Getenv("PRIVATE_KEY_PATH")
	pathReviewChecklist := os.Getenv("REVIEW_CHECKLIST_PATH")
	if integrationID == "" {
		return nil, errors.New("GITHUB_APP_INTEGRATION_ID environment variable is required but not set")
	}
	if webhookSecret == "" {
		return nil, errors.New("GITHUB_APP_WEBHOOK_SECRET environment variable is required but not set")
	}
	if pathPrivateKey == "" {
		return nil, errors.New("PRIVATE_KEY_PATH environment variable is required but not set")
	}
	if pathReviewChecklist == "" {
		return nil, errors.New("REVIEW_CHECKLIST_PATH environment variable is required but not set")
	}
	bytes, err := os.ReadFile(pathPrivateKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read private key file: %s", pathPrivateKey)
	}
	c.Github.App.PrivateKey = string(bytes)

	bytes, err = os.ReadFile(pathReviewChecklist)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read review checklist file: %s", pathReviewChecklist)
	}
	c.reviewChecklist = string(bytes)

	c.botLogin = os.Getenv("BOT_USER_LOGIN")

	// Get server address
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "127.0.0.1"
	}
	c.address = serverAddress

	// Get log file path
	c.logFile = os.Getenv("LOG_FILE")
	return &c, nil
}
