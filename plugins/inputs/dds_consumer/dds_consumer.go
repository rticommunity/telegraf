//go:build dds

/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

package dds_consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers/json"

	// NOTE: This import requires the RTI Connext DDS Connector for Go
	// To install: go get github.com/rticommunity/rticonnextdds-connector-go
	rti "github.com/rticommunity/rticonnextdds-connector-go"
)

type DDSConsumer struct {
	// XML configuration file path
	ConfigFilePath string `toml:"config_path"`
	// XML configuration name for DDS Participant
	ParticipantConfig string `toml:"participant_config"`
	// XML configuration names for DDS Readers
	ReaderConfig string   `toml:"reader_config"`
	TagKeys      []string `toml:"tag_keys"`

	// RTI Connext Connector entities
	connector *rti.Connector
	reader    *rti.Input

	// Telegraf entities
	parser *json.Parser
	acc    telegraf.Accumulator

	// Shutdown coordination
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Default configurations
var sampleConfig = `
  ## XML configuration file path
  config_path = "example_configs/ShapeExample.xml"

  ## Configuration name for DDS Participant from a description in XML
  participant_config = "MyParticipantLibrary::Zero"

  ## Configuration name for DDS DataReader from a description in XML
  reader_config = "MySubscriber::MySquareReader"

  ## Tag key is an array of keys that should be added as tags.
  tag_keys = ["color"]

  ## Override the base name of the measurement
  name_override = "shapes"
`

func checkFatalError(err error) {
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}

func checkError(err error) {
	if err != nil {
		log.Println("ERROR:", err)
	}
}

func (d *DDSConsumer) SampleConfig() string {
	return sampleConfig
}

func (d *DDSConsumer) Description() string {
	return "Read metrics from DDS"
}

func (d *DDSConsumer) Start(acc telegraf.Accumulator) error {
	// Keep the Telegraf accumulator internally
	d.acc = acc

	// Initialize shutdown coordination
	d.ctx, d.cancel = context.WithCancel(context.Background())

	var err error

	// Create a Connector entity
	d.connector, err = rti.NewConnector(d.ParticipantConfig, d.ConfigFilePath)
	if err != nil {
		return err
	}

	// Get a DDS reader
	d.reader, err = d.connector.GetInput(d.ReaderConfig)
	if err != nil {
		d.connector.Delete()
		return err
	}

	// Initialize JSON parser
	d.parser = &json.Parser{
		MetricName: "dds",
		TagKeys:    d.TagKeys,
	}
	err = d.parser.Init()
	if err != nil {
		d.connector.Delete()
		return err
	}

	// Start a thread for ingesting DDS
	d.wg.Add(1)
	go d.process()

	return nil
}

func (d *DDSConsumer) Stop() {
	// Signal the process goroutine to stop
	if d.cancel != nil {
		d.cancel()
	}

	// Wait for the process goroutine to finish
	d.wg.Wait()

	// Now safely delete the connector
	if d.connector != nil {
		d.connector.Delete()
		d.connector = nil
	}
}

// Take DDS samples from the DataReader and ingest them to Telegraf outputs
func (d *DDSConsumer) process() {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			// Shutdown signal received
			log.Println("DDS Consumer: Stopping processing loop")
			return
		default:
			// Continue processing
		}

		// Use a timeout for Wait to avoid blocking indefinitely
		waitTimeout := 1000 // 1 second timeout in milliseconds
		d.connector.Wait(waitTimeout)

		err := d.reader.Take()
		if err != nil {
			checkError(err)
			continue
		}

		numOfSamples, err := d.reader.Samples.GetLength()
		checkError(err)
		if err != nil {
			continue
		}

		for i := 0; i < numOfSamples; i++ {
			valid, err := d.reader.Infos.IsValid(i)
			checkError(err)
			if err != nil {
				continue
			}
			if valid {
				json, err := d.reader.Samples.GetJSON(i)
				checkError(err)
				if err != nil {
					continue
				}
				ts, err := d.reader.Infos.GetSourceTimestamp(i)
				checkError(err)
				if err != nil {
					continue
				}

				// Process synchronously to avoid goroutine leaks on shutdown
				d.processMessage(json, ts)
			}
		}
	}
}

// Helper function to process individual messages
func (d *DDSConsumer) processMessage(jsonStr string, ts int64) {
	// Parse the JSON object to metrics
	metrics, err := d.parser.Parse([]byte(jsonStr))
	checkError(err)
	if err != nil {
		return
	}

	// Iterate the metrics
	for _, metric := range metrics {
		// Add a metric to an output plugin
		d.acc.AddFields(metric.Name(), metric.Fields(), metric.Tags(), time.Unix(0, ts))
	}
}

func (d *DDSConsumer) Gather(acc telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("dds_consumer", func() telegraf.Input {
		return &DDSConsumer{}
	})
}
