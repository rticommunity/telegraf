/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

package dds_consumer_lp

import (
	"errors"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
	"github.com/rticommunity/rticonnextdds-connector-go"
	"path"
	"runtime"
	"time"
)

type DDSConsumer struct {
	// RTI Connext Connector entities
	connector *rti.Connector
	reader    *rti.Input

	// Telegraf entities
	parser parsers.Parser
	acc    telegraf.Accumulator
}

// Default configurations
var sampleConfig = `
`

func (d *DDSConsumer) SampleConfig() string {
	return sampleConfig
}

func (d *DDSConsumer) Description() string {
	return "Read metrics from DDS"
}

func (d *DDSConsumer) SetParser(parser parsers.Parser) {
	d.parser = parser
}

func (d *DDSConsumer) Start(acc telegraf.Accumulator) (err error) {
	// Find the file path to the XML configuration
	_, cur_path, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("cannot get the path for XML config file")
	}
	filepath := path.Join(path.Dir(cur_path), "./dds_consumer_lp.xml")

	// Keep the Telegraf accumulator internally
	d.acc = acc

	// Create a Connector entity
	d.connector, err = rti.NewConnector("MyParticipantLibrary::Zero", filepath)
	if err != nil {
		return err
	}

	// Get a DDS reader
	d.reader, err = d.connector.GetInput("MySubscriber::MyReader")
	if err != nil {
		return err
	}

	// Start a thread for ingesting DDS
	go d.process()

	return nil
}

func (d *DDSConsumer) Stop() {
	d.connector.Delete()
}

// Take DDS samples from the DataReader and ingest them to Telegraf outputs
func (d *DDSConsumer) process() {
	for {
		d.connector.Wait(-1)
		d.reader.Take()
		numOfSamples := d.reader.Samples.GetLength()

		for i := 0; i < numOfSamples; i++ {
			if d.reader.Infos.IsValid(i) {
				name := d.reader.Samples.GetString(i, "name")

				tags := make(map[string]string)
				tag_length := d.reader.Samples.GetInt(i, "tags#")
				for j := 0; j < tag_length; j++ {
					key := fmt.Sprintf("tags[%d].key", j+1)
					value := fmt.Sprintf("tags[%d].value", j+1)
					tags[d.reader.Samples.GetString(i, key)] = d.reader.Samples.GetString(i, value)
				}

				fields := make(map[string]interface{})
				field_length := d.reader.Samples.GetInt(i, "fields#")
				for j := 0; j < field_length; j++ {
					key := fmt.Sprintf("fields[%d].key", j+1)
					value := fmt.Sprintf("fields[%d].value", j+1)
					fields[d.reader.Samples.GetString(i, key)] = d.reader.Samples.GetFloat64(i, value)
				}

				timestamp := d.reader.Samples.GetInt64(i, "timestamp")

				d.acc.AddFields(name, fields, tags, time.Unix(0, timestamp))
				/*
					go func(json []byte) {
						//log.Println(string(json))

						// Parse the JSON object to metrics
						metrics, err := d.parser.Parse(json)
						checkError(err)

						// Iterate the metrics
						for _, metric := range metrics {
							// Add a metric to an output plugin

						}
					}(json)
				*/
			}
		}
	}
}

func (d *DDSConsumer) Gather(acc telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("dds_consumer_lp", func() telegraf.Input {
		return &DDSConsumer{}
	})
}
