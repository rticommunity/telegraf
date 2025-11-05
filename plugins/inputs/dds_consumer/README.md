# DDS Consumer Input Plugin

The DDS consumer plugin reads metrics over DDS by creating readers defined in [XML App Creation](https://community.rti.com/static/documentation/connext-dds/current/doc/manuals/connext_dds_professional/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStartedGuide.pdf) configurations. This plugin converts received DDS data to JSON data and adds to a Telegraf output plugin.

## Installation

1. **Download RTI Connector libraries**:
   ```bash
   go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest
   ```

2. **Set library path** (choose your platform):
   ```bash
   # macOS (Apple Silicon/ARM64)
   export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-arm64:$DYLD_LIBRARY_PATH
   
   # macOS (Intel/x86_64)  
   export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-x64:$DYLD_LIBRARY_PATH
   
   # Linux  
   export LD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH
   ```

3. **Build Telegraf with DDS support**:
   ```bash
   go build -tags dds ./cmd/telegraf
   ```

4. **Test with example configuration**:
   ```bash
   ./telegraf --config plugins/inputs/dds_consumer/example_telegraf.conf
   ```

## Test with RTI Shapes Demo

1. Download [RTI Shapes Demo](https://www.rti.com/free-trial/shapes-demo)
2. Publish "Square" shapes
3. Run the command above - you should see shape data printed to console

## Configuration

```toml
[[inputs.dds_consumer]]
  # Path to your DDS XML configuration file
  config_path = "plugins/inputs/dds_consumer/example_config.xml"
  
  # DDS Participant configuration from XML
  participant_config = "MyParticipantLibrary::Zero"
  
  # DDS DataReader configuration from XML  
  reader_config = "MySubscriber::MySquareReader"
  
  # Fields to use as metric tags (instead of fields)
  tag_keys = ["color"]
  
  # Name for the measurement
  name_override = "shapes"
```

**Example files**: [`example_telegraf.conf`](./example_telegraf.conf) and [`example_config.xml`](./example_config.xml)

## Troubleshooting

- **Library not found**: Check the library path for your platform
- **Plugin not available**: Ensure you built with `-tags dds`
- **No data**: Verify DDS publisher is running and on the same domain
- **XML errors**: Validate XML syntax and participant/reader names