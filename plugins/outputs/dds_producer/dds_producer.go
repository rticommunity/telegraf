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
	"fmt"
	"time"
	"runtime"
	"errors"
	"path"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/rticommunity/rticonnextdds-connector-go"
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

var sampleConfig = `
`

func (d *DDSProducer) SetSerializer(serializer serializers.Serializer) {
	d.serializer = serializer
}

func (d *DDSProducer) Connect() (err error) {
    // Find the file path to the XML configuration
    _, cur_path, _, ok := runtime.Caller(0)
    if !ok {
        return errors.New("cannot get the path for XML config file")
    }
    filepath := path.Join(path.Dir(cur_path), "./dds_producer.xml")

    // Create a Connector entity
    d.connector, err = rti.NewConnector("MyParticipantLibrary::Zero", filepath)
    if err != nil {
        return err
    }

    // Get a DDS reader
    d.writer, err = d.connector.GetOutput("MyPublisher::MyWriter")
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
	return "Send metrics to DDS"
}

func (d *DDSProducer) Write(metrics []telegraf.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	for _, metric := range metrics {
		err := d.writer.Instance.SetString("name", metric.Name())
		if err != nil {
			return err
		}

		for i, tag := range metric.TagList() {
			key := fmt.Sprintf("tags[%d].key", i+1)
			value := fmt.Sprintf("tags[%d].value", i+1)
			err = d.writer.Instance.SetString(key, tag.Key)
			if err != nil {
				return err
			}
			err = d.writer.Instance.SetString(value, tag.Value)
			if err != nil {
				return err
			}
		}
		for i, field := range metric.FieldList() {
			key := fmt.Sprintf("fields[%d].key", i+1)
			value := fmt.Sprintf("fields[%d].value", i+1)
			err = d.writer.Instance.SetString(key, field.Key)
			if err != nil {
				return err
			}

			v := convertField(field.Value)

			err = d.writer.Instance.SetFloat64(value, v)
			if err != nil {
				return err
			}

		}

		d.writer.Instance.SetInt64("timestamp", time.Now().UTC().UnixNano())
		if err != nil {
			return err
		}

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
	outputs.Add("dds_producer", func() telegraf.Output {
		return &DDSProducer{}
	})
}
