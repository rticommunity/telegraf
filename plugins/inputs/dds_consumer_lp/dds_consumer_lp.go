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
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
	"github.com/rticommunity/rticonnextdds-connector-go"
	"time"
)

type DDSConsumer struct {
	// XML configuration file path
	DomainId string `toml:"domain_id"`

	// RTI Connext Connector entities
	connector *rti.Connector
	reader    *rti.Input

	// Telegraf entities
	parser parsers.Parser
	acc    telegraf.Accumulator
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	FIELD_DOUBLE = 0
	FIELD_INT    = 1
	FIELD_UINT   = 2
	FIELD_STRING = 3
	FIELD_BOOL   = 4
)

type FieldValueUnion struct {
	D *float64 `json:"d,omitempty"`
	I *int64   `json:"i,omitempty"`
	U *uint64  `json:"u,omitempty"`
	S *string  `json:"s,omitempty"`
	B *bool    `json:"b,omitempty"`
}

type Field struct {
	Key   string          `json:"key"`
	Kind  int             `json:"kind"`
	Value FieldValueUnion `json:"value"`
}

type Metric struct {
	Name      string  `json:"name"`
	Tags      []Tag   `json:"tags"`
	Fields    []Field `json:"fields"`
	Timestamp int64   `json:"timestamp"`
}

// Default configurations
var sampleConfig = `
domain_id = "0"
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
	d.acc = acc

	var xmlString = `
	str://"<dds>
	<qos_library name="QosLibrary">
	<qos_profile name="DefaultProfile" base_name="BuiltinQosLibExp::Generic.StrictReliable" is_default_qos="true">
	</qos_profile>
	</qos_library>
	<types>
	<include file="line_protocol.xml"/>
	</types>
	<domain_library name="DomainLib">
	<domain name="Telegraf" domain_id="
	` + d.DomainId +
		`">
	<register_type name="Metric" type_ref="Metric"/>
	<topic name="Telegraf" register_type_ref="Metric"/>
	</domain>
	</domain_library>
	<domain_participant_library name="ParticipantLib">
	<domain_participant name="TelegrafParticipant" domain_ref="DomainLib::Telegraf">
	<subscriber name="TelegrafSubscriber">
	<data_reader name="TelegrafReader" topic_ref="Telegraf"/>
	</subscriber>
	</domain_participant>
	</domain_participant_library>
	</dds>"
	`
	// Create a Connector object from the XML config
	d.connector, err = rti.NewConnector("ParticipantLib::TelegrafParticipant", xmlString)
	if err != nil {
		return err
	}

	// Get a DDS reader
	d.reader, err = d.connector.GetInput("TelegrafSubscriber::TelegrafReader")
	if err != nil {
		return err
	}

	// Start a go routine for processing DDS samples
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
				var m Metric
				d.reader.Samples.Get(i, &m)

				tags := make(map[string]string)
				for _, tag := range m.Tags {
					tags[tag.Key] = tag.Value
				}

				fields := make(map[string]interface{})
				for _, field := range m.Fields {
					switch field.Kind {
					case FIELD_DOUBLE:
						fields[field.Key] = field.Value.D
					case FIELD_INT:
						fields[field.Key] = field.Value.I
					case FIELD_UINT:
						fields[field.Key] = field.Value.U
					case FIELD_STRING:
						fields[field.Key] = field.Value.S
					case FIELD_BOOL:
						fields[field.Key] = field.Value.B
					default:
					}
				}

				d.acc.AddFields(m.Name, fields, tags, time.Unix(0, m.Timestamp))
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
