# DDS Monitor Plugin

This input pluin monitors metrics for DDS applications. It requires `monitoring.xml` and `dds_rtf2_dcps.xml` to be located in your current path for running Telegraf. 

### Configuration:

```toml
[[inputs.dds_monitor]]
  # DDS Domain ID
  domain_id = "204"

  # Interval of polling DDS data
  interval = 1

  ## Data format to consume.
  data_format = "json"
```
