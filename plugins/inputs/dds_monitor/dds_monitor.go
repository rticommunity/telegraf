/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

package dds_monitor

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
	"github.com/rticommunity/rticonnextdds-connector-go"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	defaultDomainId = "0"
	defaultInterval = 1
)

type DDSConsumer struct {
	// DDS Domain ID
	DomainId string `toml:"domain_id"`

	// Interval of polling DDS data in second
	Interval float64 `toml:"interval"`

	// RTI Connext Connector entities
	connector *rti.Connector
	readers   map[string]*rti.Input

	// Telegraf entities
	parser parsers.Parser
	acc    telegraf.Accumulator
}

// Default configurations
var sampleConfig = `
  ## DDS Domain ID
  domain_id = "0"

  ## Interval of polling DDS data in second
  interval = 1

  ## Data format to consume.
  data_format = "json"
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

func (d *DDSConsumer) SetParser(parser parsers.Parser) {
	d.parser = parser
}

func (d *DDSConsumer) Start(acc telegraf.Accumulator) error {
	// Keep the Telegraf accumulator internally
	d.acc = acc

	var xmlString = `
    str://"<dds>
    <qos_library name="QosLibrary">
    <qos_profile name="DefaultProfile" is_default_qos="true">
    </qos_profile>

    </qos_library>
    <types xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="file:////home/kyounghoan/rti_connext_dds-6.0.0/bin/../resource/app/app_support/rtiddsgen/schema/rti_dds_topic_types.xsd">
	<include file="monitoring.xml"/>
    </types>
    <domain_library name="DomainLib">
    <domain name="DDSMonitor" domain_id="
    ` + d.DomainId +
		`">
            <register_type name="rti::dds::monitoring::DomainParticipantDescription"
                           type_ref="rti::dds::monitoring::DomainParticipantDescription" />
            <register_type name="rti::dds::monitoring::PublisherDescription"
                            type_ref="rti::dds::monitoring::PublisherDescription" />
            <register_type name="rti::dds::monitoring::SubscriberDescription"
                            type_ref="rti::dds::monitoring::SubscriberDescription" />
            <register_type name="rti::dds::monitoring::TopicDescription"
                            type_ref="rti::dds::monitoring::TopicDescription" />
            <register_type name="rti::dds::monitoring::DataWriterDescription"
                            type_ref="rti::dds::monitoring::DataWriterDescription" />
            <register_type name="rti::dds::monitoring::DataReaderDescription"
                            type_ref="rti::dds::monitoring::DataReaderDescription" />
            <register_type name="rti::dds::monitoring::DomainParticipantEntityStatistics"
                            type_ref="rti::dds::monitoring::DomainParticipantEntityStatistics" />
            <register_type name="rti::dds::monitoring::DataReaderEntityStatistics"
                            type_ref="rti::dds::monitoring::DataReaderEntityStatistics" />
            <register_type name="rti::dds::monitoring::DataWriterEntityStatistics"
                            type_ref="rti::dds::monitoring::DataWriterEntityStatistics" />
            <register_type name="rti::dds::monitoring::TopicEntityStatistics"
                            type_ref="rti::dds::monitoring::TopicEntityStatistics" />
            <register_type name="rti::dds::monitoring::DataReaderEntityMatchedPublicationStatistics"
                            type_ref="rti::dds::monitoring::DataReaderEntityMatchedPublicationStatistics" />

            <topic name="rti/dds/monitoring/domainParticipantDescription"
                   register_type_ref="rti::dds::monitoring::DomainParticipantDescription"/>

            <topic name="rti/dds/monitoring/publisherDescription"
                   register_type_ref="rti::dds::monitoring::PublisherDescription"/>

            <topic name="rti/dds/monitoring/subscriberDescription"
                   register_type_ref="rti::dds::monitoring::SubscriberDescription"/>

            <topic name="rti/dds/monitoring/topicDescription"
                   register_type_ref="rti::dds::monitoring::TopicDescription"/>

            <topic name="rti/dds/monitoring/dataWriterDescription"
                   register_type_ref="rti::dds::monitoring::DataWriterDescription"/>

            <topic name="rti/dds/monitoring/dataReaderDescription"
                   register_type_ref="rti::dds::monitoring::DataReaderDescription"/>

            <topic name="rti/dds/monitoring/domainParticipantEntityStatistics"
                   register_type_ref="rti::dds::monitoring::DomainParticipantEntityStatistics"/>

            <topic name="rti/dds/monitoring/dataReaderEntityStatistics"
                   register_type_ref="rti::dds::monitoring::DataReaderEntityStatistics"/>

            <topic name="rti/dds/monitoring/dataWriterEntityStatistics"
                   register_type_ref="rti::dds::monitoring::DataWriterEntityStatistics"/>

            <topic name="rti/dds/monitoring/topicEntityStatistics"
                   register_type_ref="rti::dds::monitoring::TopicEntityStatistics"/>

            <topic name="rti/dds/monitoring/dataReaderEntityMatchedPublicationStatistics"
                   register_type_ref="rti::dds::monitoring::DataReaderEntityMatchedPublicationStatistics"/>

    </domain>
    </domain_library>
    <domain_participant_library name="ParticipantLib">
    <domain_participant name="DDSMonitorParticipant" domain_ref="DomainLib::DDSMonitor">
    <subscriber name="DDSMonitorSubscriber">

                <data_reader name="DomainParticipantDescriptionReader"
                             topic_ref="rti/dds/monitoring/domainParticipantDescription"/>

                <data_reader name="PublisherDescriptionReader"
                             topic_ref="rti/dds/monitoring/publisherDescription"/>

                <data_reader name="SubscriberDescriptionReader"
                             topic_ref="rti/dds/monitoring/subscriberDescription"/>

                <data_reader name="TopicDescriptionReader"
                             topic_ref="rti/dds/monitoring/topicDescription"/>

                <data_reader name="DataWriterDescriptionReader"
                             topic_ref="rti/dds/monitoring/dataWriterDescription"/>

                <data_reader name="DataReaderDescriptionReader"
                             topic_ref="rti/dds/monitoring/dataReaderDescription"/>

                <data_reader name="DomainParticipantEntityStatisticsReader"
                             topic_ref="rti/dds/monitoring/domainParticipantEntityStatistics"/>

                <data_reader name="DataReaderEntityStatisticsReader"
                             topic_ref="rti/dds/monitoring/dataReaderEntityStatistics"/>

                <data_reader name="DataWriterEntityStatisticsReader"
                             topic_ref="rti/dds/monitoring/dataWriterEntityStatistics"/>

                <data_reader name="TopicEntityStatisticsReader"
                             topic_ref="rti/dds/monitoring/topicEntityStatistics"/>

                <data_reader name="DataReaderEntityMatchedPublicationStatisticsReader"
                             topic_ref="rti/dds/monitoring/dataReaderEntityMatchedPublicationStatistics"/>

    </subscriber>
    </domain_participant>
    </domain_participant_library>
    </dds>"
    `
	var err error
	d.readers = make(map[string]*rti.Input)

	if d.DomainId == "" {
		d.DomainId = defaultDomainId
	}
	if d.Interval == 0 {
		d.Interval = defaultInterval
	}

	// Create a Connector entity
	d.connector, err = rti.NewConnector("ParticipantLib::DDSMonitorParticipant", xmlString)
	checkFatalError(err)

	// Get a DDS reader
	d.readers["ParticipantStats"], err = d.connector.GetInput("DDSMonitorSubscriber::DomainParticipantEntityStatisticsReader")
	d.readers["WriterStats"], err = d.connector.GetInput("DDSMonitorSubscriber::DataWriterEntityStatisticsReader")
	d.readers["ReaderStats"], err = d.connector.GetInput("DDSMonitorSubscriber::DataReaderEntityStatisticsReader")
	checkFatalError(err)

	// Start a thread for reading and processing DDS metrics
	go d.read()

	return nil
}

func (d *DDSConsumer) Stop() {
	d.connector.Delete()
}

func (d *DDSConsumer) process(key string, json []byte) {
	// Parse the JSON object to metrics
	metrics, err := d.parser.Parse(json)
	checkError(err)

	// Iterate the metrics
	for _, metric := range metrics {

		var metricName string

		switch key {
		case "ParticipantStats":
			metricName = "dds_participant_stats"

			// Make Domain ID as a tag
			value, _ := metric.GetField("domain_id")
			strValue := strconv.FormatFloat(value.(float64), 'f', -1, 64)
			metric.AddTag("domain_id", strValue)
			metric.RemoveField("domain_id")

			// Make Process ID as a tag
			value, _ = metric.GetField("process_id")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			metric.AddTag("process_id", strValue)
			metric.RemoveField("process_id")

			// Make Participant ID as a tag
			value, _ = metric.GetField("participant_key_value_0")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId := strValue
			metric.RemoveField("participant_key_value_0")
			value, _ = metric.GetField("participant_key_value_1")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_1")
			value, _ = metric.GetField("participant_key_value_2")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_2")
			value, _ = metric.GetField("participant_key_value_3")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_3")
			metric.AddTag("participant_id", participantId)

			// Remove fields not needed
			metric.RemoveField("period_sec")
			metric.RemoveField("period_nanosec")
			metric.RemoveField("host_id")

		case "WriterStats":
			metricName = "dds_writer_stats"
			// Make Domain ID as a tag
			value, _ := metric.GetField("domain_id")
			strValue := strconv.FormatFloat(value.(float64), 'f', -1, 64)
			metric.AddTag("domain_id", strValue)
			metric.RemoveField("domain_id")

			// Make Process ID as a tag
			value, _ = metric.GetField("process_id")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			metric.AddTag("process_id", strValue)
			metric.RemoveField("process_id")

			// Make Participant ID as a tag
			value, _ = metric.GetField("participant_key_value_0")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId := strValue
			metric.RemoveField("participant_key_value_0")
			value, _ = metric.GetField("participant_key_value_1")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_1")
			value, _ = metric.GetField("participant_key_value_2")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_2")
			value, _ = metric.GetField("participant_key_value_3")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_3")
			metric.AddTag("participant_id", participantId)

			// Make Writer ID as a tag
			value, _ = metric.GetField("datawriter_key_value_0")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			writerId := strValue
			metric.RemoveField("datawriter_key_value_0")
			value, _ = metric.GetField("datawriter_key_value_1")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			writerId += "." + strValue
			metric.RemoveField("datawriter_key_value_1")
			value, _ = metric.GetField("datawriter_key_value_2")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			writerId += "." + strValue
			metric.RemoveField("datawriter_key_value_2")
			value, _ = metric.GetField("datawriter_key_value_3")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			writerId += "." + strValue
			metric.RemoveField("datawriter_key_value_3")
			metric.AddTag("datawriter_id", writerId)

			// Remove fields not needed
			metric.RemoveField("period_sec")
			metric.RemoveField("period_nanosec")
			metric.RemoveField("host_id")
			metric.RemoveField("publisher_key_value_0")
			metric.RemoveField("publisher_key_value_1")
			metric.RemoveField("publisher_key_value_2")
			metric.RemoveField("publisher_key_value_3")
			metric.RemoveField("topic_key_value_0")
			metric.RemoveField("topic_key_value_1")
			metric.RemoveField("topic_key_value_2")
			metric.RemoveField("topic_key_value_3")
		case "ReaderStats":
			metricName = "dds_reader_stats"

			// Make Domain ID as a tag
			value, _ := metric.GetField("domain_id")
			strValue := strconv.FormatFloat(value.(float64), 'f', -1, 64)
			metric.AddTag("domain_id", strValue)
			metric.RemoveField("domain_id")

			// Make Process ID as a tag
			value, _ = metric.GetField("process_id")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			metric.AddTag("process_id", strValue)
			metric.RemoveField("process_id")

			// Make Participant ID as a tag
			value, _ = metric.GetField("participant_key_value_0")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId := strValue
			metric.RemoveField("participant_key_value_0")
			value, _ = metric.GetField("participant_key_value_1")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_1")
			value, _ = metric.GetField("participant_key_value_2")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_2")
			value, _ = metric.GetField("participant_key_value_3")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			participantId += "." + strValue
			metric.RemoveField("participant_key_value_3")
			metric.AddTag("participant_id", participantId)

			// Make Reader ID as a tag
			value, _ = metric.GetField("datareader_key_value_0")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			readerId := strValue
			metric.RemoveField("datareader_key_value_0")
			value, _ = metric.GetField("datareader_key_value_1")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			readerId += "." + strValue
			metric.RemoveField("datareader_key_value_1")
			value, _ = metric.GetField("datareader_key_value_2")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			readerId += "." + strValue
			metric.RemoveField("datareader_key_value_2")
			value, _ = metric.GetField("datareader_key_value_3")
			strValue = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			readerId += "." + strValue
			metric.RemoveField("datareader_key_value_3")
			metric.AddTag("datareader_id", readerId)

			// Remove fields not needed
			metric.RemoveField("period_sec")
			metric.RemoveField("period_nanosec")
			metric.RemoveField("host_id")
			metric.RemoveField("subscriber_key_value_0")
			metric.RemoveField("subscriber_key_value_1")
			metric.RemoveField("subscriber_key_value_2")
			metric.RemoveField("subscriber_key_value_3")
			metric.RemoveField("topic_key_value_0")
			metric.RemoveField("topic_key_value_1")
			metric.RemoveField("topic_key_value_2")
			metric.RemoveField("topic_key_value_3")
			metric.RemoveField("datareader_protocol_status_status_first_available_sample_sequence_number_high")
			metric.RemoveField("datareader_protocol_status_status_first_available_sample_sequence_number_low")
			metric.RemoveField("datareader_protocol_status_status_last_available_sample_sequence_number_low")
			metric.RemoveField("datareader_protocol_status_status_last_available_sample_sequence_number_high")
			metric.RemoveField("datareader_protocol_status_status_last_committed_sample_sequence_number_high")
			metric.RemoveField("datareader_protocol_status_status_last_committed_sample_sequence_number_low")
		default:
		}

		// Remove filed including "change" or "handle_value"
		metric.SetName(metricName)
		for k, _ := range metric.Fields() {
			if strings.Contains(k, "change") || strings.Contains(k, "_handle_") {
				metric.RemoveField(k)
			}
		}

		// Add a metric to an output plugin
		d.acc.AddFields(metric.Name(), metric.Fields(), metric.Tags(), metric.Time())
	}
}

// Take DDS samples from the DataReader and ingest them to Telegraf outputs
func (d *DDSConsumer) read() {
	for {
		for key, reader := range d.readers {
			reader.Take()
			numOfSamples := reader.Samples.GetLength()

			for i := 0; i < numOfSamples; i++ {
				if reader.Infos.IsValid(i) {
					json, err := reader.Samples.GetJSON(i)
					checkError(err)
					go d.process(key, json)
				}
			}
		}
		time.Sleep(time.Duration(d.Interval) * time.Second)
	}
}

func (d *DDSConsumer) Gather(acc telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("dds_monitor", func() telegraf.Input {
		return &DDSConsumer{}
	})
}
