# DDS Consumer Input Plugin

The DDS consumer plugin reads with DDS readers defined in a configuration for [XML App Creation](https://community.rti.com/static/documentation/connext-dds/5.3.1/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf). This plugin addes received DDS data in the Line Protocol format for Telegraf output plugins such as the InfluxDB plugin. 

### Configuration:

```toml
[[inputs.dds_consumer_lp]]

```
