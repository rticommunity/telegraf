# DDS Consumer Input Plugin

The DDS consumer plugin reads with DDS readers defined in a configuration for [XML App Creation](https://community.rti.com/static/documentation/connext-dds/5.3.1/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf). This plugin addes received DDS data in the Line Protocol format for Telegraf output plugins such as the InfluxDB plugin. 

### Configuration:

```toml
[[inputs.dds_consumer_lp]]
  ## XML configuration file path
  config_path = "dds_consumer_lp.xml"

  ## Configuration name for DDS Participant from a description in XML
  participant_config = "MyParticipantLibrary::Zero"

  ## Configuration name for DDS DataReader from a description in XML
  reader_config = "MySubscriber::MyReader"

  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  #data_format = "json"

```
