# DDS Consumer Input Plugin

The DDS consumer plugin reads with DDS readers defined in a configuration for [XML App Creation](https://community.rti.com/static/documentation/connext-dds/5.3.1/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf). This plugin converts received DDS data to JSON data and adds to a Telegraf output plugin like the InfluxDB plugin. 

### Configuration:

```toml
[[inputs.dds_consumer]]
  ## XML configuration file path
  config_path = "ShapeExample.xml"

  ## Configuration name for DDS Participant from a description in XML
  participant_config = "MyParticipantLibrary::Zero"

  ## Configuration name for DDS DataReader from a description in XML
  reader_config = "MySubscriber::MySquareReader"

  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "json"

```
