# DDS Consumer Input Plugin

The DDS consumer plugin reads metrics over DDS by creating readers defined in [XML App Creation](https://community.rti.com/static/documentation/connext-dds/5.3.1/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf) configurations. This plugin converts received DDS data to JSON data and adds to a Telegraf output plugin. 

### Configuration:

```toml
[[inputs.dds_consumer]]
  ## XML configuration file path
  config_path = "example_configs/ShapeExample.xml"

  ## Configuration name for DDS Participant from a description in XML
  participant_config = "MyParticipantLibrary::Zero"

  ## Configuration name for DDS DataReader from a description in XML
  reader_config = "MySubscriber::MySquareReader"

  # Tag key is an array of keys that should be added as tags.
  tag_keys = ["color"]

  # Override the base name of the measurement
  name_override = "shapes"

  ## Data format to consume.
  data_format = "json"
```
