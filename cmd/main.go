package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := route53domains.NewFromConfig(cfg, func(o *route53domains.Options) {
		o.Region = "us-east-1"
	})

	domain := "awaissaeed.com"
	output, err := client.CheckDomainAvailability(context.TODO(), &route53domains.CheckDomainAvailabilityInput{
		DomainName: &domain,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", output.Availability)
}
