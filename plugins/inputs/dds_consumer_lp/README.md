# DDS Input Plugin with the Line Protocol Data Model

```dds_consumer_lp``` plugin reads DDS data in a data model based on the Line Protocol format. This plugin creates a DDS reader with a DDS topic named `Telegraf`. The data type used by the reader is described in `line_protocol.idl`. 

### Configuration:

```toml
[[inputs.dds_consumer_lp]]
# DDS Domain ID configuration
domain_id = "0"
```
