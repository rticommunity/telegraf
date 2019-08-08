/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

package dds_producer

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/rticommunity/rticonnextdds-connector-go"
	"time"
)

type DDSProducer struct {
	// XML configuration file path
	ConfigFilePath    string `toml:"config_path"`
	ParticipantConfig string `toml:"participant_config"`
	WriterConfig      string `toml:"writer_config"`

	connector *rti.Connector
	writer    *rti.Output

	serializer serializers.Serializer
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	CPU_USAGE = 0
	MEM       = 1
)

type CpuUsageType struct {
	UsageUser      float64 `json:"usage_user"`
	UsageSystem    float64 `json:"usage_system"`
	UsageIdle      float64 `json:"usage_idle"`
	UsageActive    float64 `json:"usage_active"`
	UsageNice      float64 `json:"usage_nice"`
	UsageIoWait    float64 `json:"usage_iowait"`
	UsageIrq       float64 `json:"usage_irq"`
	UsageSoftirq   float64 `json:"usage_softirq"`
	UsageSteal     float64 `json:"usage_steal"`
	UsageGuest     float64 `json:"usage_guest"`
	UsageGuestNice float64 `json:"usage_guest_nice"`
}

type MemType struct {
	Active    uint64 `json:"active"`
	Available uint64 `json:"available"`
}

type FieldValueUnion struct {
	CpuUsage *CpuUsageType `json:"cpu_usage,omitempty"`
	Mem      *MemType      `json:"mem,omitempty"`
}

type Field struct {
	Kind  int         `json:"kind"`
	Value interface{} `json:"value"`
}

type Metric struct {
	Name      string `json:"name"`
	Tags      []Tag  `json:"tags"`
	Fields    Field  `json:"fields"`
	Timestamp int64  `json:"timestamp"`
}

var sampleConfig = `
  ## XML configuration file path
  config_path = "dds_producer.xml"
   ## Configuration name for DDS Participant from a description in XML
  participant_config = "MyParticipantLibrary::Zero"
   ## Configuration name for DDS DataWriter from a description in XML
  writer_config = "MyPublisher::MyWriter"
`

func (d *DDSProducer) SetSerializer(serializer serializers.Serializer) {
	d.serializer = serializer
}

func (d *DDSProducer) Connect() (err error) {
	// Create a Connector
	d.connector, err = rti.NewConnector(d.ParticipantConfig, d.ConfigFilePath)
	if err != nil {
		return err
	}

	// Get a DDS Writer
	d.writer, err = d.connector.GetOutput(d.WriterConfig)
	if err != nil {
		return err
	}

	return nil
}

func (d *DDSProducer) Close() error {
	d.connector.Delete()
	return nil
}

func (d *DDSProducer) SampleConfig() string {
	return sampleConfig
}

func (d *DDSProducer) Description() string {
	return "Send metrics over DDS"
}

func (d *DDSProducer) Write(metrics []telegraf.Metric) (err error) {
	if len(metrics) == 0 {
		return nil
	}

	for _, metric := range metrics {
		var m Metric
		m.Name = metric.Name()

		for _, tag := range metric.TagList() {
			var t Tag
			t.Key = tag.Key
			t.Value = tag.Value
			m.Tags = append(m.Tags, t)
		}
		switch metric.Name() {
		case "cpu":
			m.Fields.Kind = CPU_USAGE
			m.Fields.Value = FieldValueUnion{
				CpuUsage: &CpuUsageType{
					UsageUser:   metric.Fields()["usage_user"].(float64),
					UsageSystem: metric.Fields()["usage_system"].(float64),
					UsageIdle:   metric.Fields()["usage_idle"].(float64),
					//UsageActive: metric.Fields()["usage_active"].(float64),
					UsageNice:      metric.Fields()["usage_nice"].(float64),
					UsageIoWait:    metric.Fields()["usage_iowait"].(float64),
					UsageIrq:       metric.Fields()["usage_irq"].(float64),
					UsageSoftirq:   metric.Fields()["usage_softirq"].(float64),
					UsageSteal:     metric.Fields()["usage_steal"].(float64),
					UsageGuest:     metric.Fields()["usage_guest"].(float64),
					UsageGuestNice: metric.Fields()["usage_guest_nice"].(float64),
				},
			}
		case "mem":
			m.Fields.Kind = MEM
			m.Fields.Value = FieldValueUnion{
				Mem: &MemType{
					Active:    metric.Fields()["active"].(uint64),
					Available: metric.Fields()["available"].(uint64),
				},
			}
		}

		m.Timestamp = time.Now().UTC().UnixNano()

		d.writer.Instance.Set(&m)

		err = d.writer.Write()
		if err != nil {
			return err
		}

		err = d.writer.ClearMembers()
		if err != nil {
			return err
		}
	}
	return nil
}

// Convert field to a supported type or nil if unconvertible
func convertField(v interface{}) float64 {
	switch v := v.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
		//case string:
		//  return v
		//case bool:
		//  return float64(v)
	case int:
		return float64(v)
	case uint:
		return float64(v)
	case uint64:
		return float64(v)
		//case []byte:
		//  return float64(v)
	case int32:
		return float64(v)
	case int16:
		return float64(v)
	case int8:
		return float64(v)
	case uint32:
		return float64(v)
	case uint16:
		return float64(v)
	case uint8:
		return float64(v)
	case float32:
		return float64(v)
	default:
		return 0
	}
}

func init() {
	outputs.Add("dds_producer_union", func() telegraf.Output {
		return &DDSProducer{}
	})
}
