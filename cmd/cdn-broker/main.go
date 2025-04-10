package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/brokerapi"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/cloud-gov/cf-cdn-service-broker/broker"
	"github.com/cloud-gov/cf-cdn-service-broker/config"
	"github.com/cloud-gov/cf-cdn-service-broker/healthchecks"
	"github.com/cloud-gov/cf-cdn-service-broker/models"
	"github.com/cloud-gov/cf-cdn-service-broker/utils"
)

func main() {
	logger := lager.NewLogger("cdn-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.INFO))

	rand.Seed(time.Now().UnixNano())

	settings, err := config.NewSettings()
	if err != nil {
		logger.Fatal("new-settings", err)
	}

	db, err := config.Connect(settings)
	if err != nil {
		logger.Fatal("connect", err)
	}

	cfClient, err := cfclient.NewClient(&cfclient.Config{
		ApiAddress:   settings.APIAddress,
		ClientID:     settings.ClientID,
		ClientSecret: settings.ClientSecret,
	})
	if err != nil {
		logger.Fatal("cf-client", err)
	}

	session := session.New(aws.NewConfig().WithRegion(settings.AwsDefaultRegion))

	if err := db.AutoMigrate(&models.Route{}, &models.Certificate{}, &models.UserData{}).Error; err != nil {
		logger.Fatal("migrate", err)
	}

	manager := models.NewManager(
		logger,
		&utils.Iam{settings, iam.New(session)},
		&utils.Distribution{settings, cloudfront.New(session)},
		settings,
		db,
	)
	broker := broker.New(
		&manager,
		cfClient,
		settings,
		logger,
	)
	credentials := brokerapi.BrokerCredentials{
		Username: settings.BrokerUsername,
		Password: settings.BrokerPassword,
	}

	brokerAPI := brokerapi.New(broker, logger, credentials)
	server := bindHTTPHandlers(brokerAPI, settings)
	http.ListenAndServe(fmt.Sprintf(":%s", settings.Port), server)
}

func bindHTTPHandlers(handler http.Handler, settings config.Settings) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	healthchecks.Bind(mux, settings)

	return mux
}
