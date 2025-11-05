# DDS Plugin Installation Guide

This guide provides step-by-step instructions for installing and configuring the DDS Consumer input plugin for Telegraf.

## Prerequisites

- **Go Environment**: Ensure Go is properly installed and configured
  - Go 1.18 or later is recommended
- **RTI Connector for Go**: This plugin uses RTI Connector libraries instead of requiring a full RTI Connext DDS installation
- **Network connectivity**: Access to download RTI Connector libraries from GitHub

## Installation Steps

### 1. Download RTI Connector Libraries

Download the RTI Connector libraries using the Go command:

```bash
go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest
```

### 2. Set Library Path (for Runtime)

Configure the library path based on your platform:

```bash
# macOS (Apple Silicon/ARM64)
export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-arm64:$DYLD_LIBRARY_PATH

# macOS (Intel/x86_64)  
export DYLD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/osx-x64:$DYLD_LIBRARY_PATH

# Linux  
export LD_LIBRARY_PATH=$(pwd)/rticonnextdds-connector/lib/linux-x64:$LD_LIBRARY_PATH

# Windows (PowerShell)
$env:PATH = "$(pwd)\rticonnextdds-connector\lib\win-x64;$env:PATH"
```

### 3. Build Telegraf with DDS Plugin

```bash
cd /path/to/telegraf
go build ./cmd/telegraf
```

### 4. Verify Installation

Test that the plugin is available:

```bash
./telegraf --input-list | grep dds_consumer
```

## Configuration

### 1. Create DDS XML Configuration

Create an XML file describing your DDS configuration (see `example_config.xml`):

```xml
<?xml version="1.0"?>
<dds>
    <types>
        <!-- Define your data types -->
    </types>
    <domain_participant_library name="MyParticipantLibrary">
        <!-- Define participants and readers -->
    </domain_participant_library>
</dds>
```

### 2. Configure Telegraf

Add the DDS consumer plugin to your Telegraf configuration:

```toml
[[inputs.dds_consumer]]
  config_path = "/path/to/your/dds_config.xml"
  participant_config = "MyParticipantLibrary::MyParticipant"
  reader_config = "MySubscriber::MyReader"
  tag_keys = ["field1", "field2"]
  name_override = "my_metrics"
  data_format = "json"
```

### 3. Test Configuration

```bash
./telegraf --config /path/to/your/telegraf.conf --test
```

## Troubleshooting

### Common Issues

1. **Library Not Found**: Verify library path is correctly set for your platform
2. **XML Parse Error**: Validate XML syntax and participant/reader names
3. **Network Issues**: Check DDS domain ID and network connectivity
4. **Architecture Mismatch**: Ensure you're using the correct library for your platform (ARM64 vs x86_64)