package dds_consumer

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
	"github.com/rticommunity/rticonnextdds-connector-go"
	"log"
)

type DDSConsumer struct {
	// XML configuration file path
	ConfigFilePath string `toml:"config_path"`
	// XML configuration name for DDS Participant
	ParticipantConfig string `toml:"participant_config"`
	// XML configuration names for DDS Readers
	ReaderConfig string `toml:"reader_config"`

	// RTI Connext Connector entities
	connector *rti.Connector
	reader    *rti.Input

	// Telegraf entities
	parser parsers.Parser
	acc    telegraf.Accumulator
}

// Default configurations
var sampleConfig = `
  ## XML configuration file path
  config_path = "USER_QOS_PROFILES.xml"
  ## Configuration name for DDS Participant from a description in XML
  participant_config = "MyParticipantLibrary::SubscriptionParticipant"
  ## Configuration name for DDS DataReader from a description in XML
  reader_config = "MySubscriber::HelloWorldReader"
  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
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
	// Disable logs
	//log.SetFlags(0)
	//log.SetOutput(ioutil.Discard)

	// Keep the Telegraf accumulator internally
	d.acc = acc

	var err error

	// Create a Connector entity
	d.connector, err = rti.NewConnector(d.ParticipantConfig, d.ConfigFilePath)
	checkFatalError(err)

	// Get a DDS reader
	d.reader, err = d.connector.GetInput(d.ReaderConfig)
	checkFatalError(err)

	// Start a thread for ingesting DDS
	go d.process()

	return nil
}

func (d *DDSConsumer) Stop() {
}

// Take DDS samples from the DataReader and ingest them to Telegraf outputs
func (d *DDSConsumer) process() {
	for {
		d.connector.Wait(-1)
		d.reader.Take()
		numOfSamples := d.reader.Samples.GetLength()

		for i := 0; i < numOfSamples; i++ {
			if d.reader.Infos.IsValid(i) {
				json, err := d.reader.Samples.GetJson(i)
				checkError(err)
				go func(json []byte) {
					//log.Println(string(json))

					// Parse the JSON object to metrics
					metrics, err := d.parser.Parse(json)
					checkError(err)

					// Iterate the metrics
					for _, metric := range metrics {
						// Add a metric to an output plugin
						d.acc.AddFields(metric.Name(), metric.Fields(), metric.Tags(), metric.Time())
					}
				}(json)
			}
		}
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
