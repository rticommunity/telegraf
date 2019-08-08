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
	// DDS Domain ID Configuration
	DomainId string `toml:"domain_id"`

	connector *rti.Connector
	writer    *rti.Output

	serializer serializers.Serializer
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	FIELD_DOUBLE = 0
	FIELD_INT    = 1
	FIELD_UINT   = 2
	FIELD_STRING = 3
	FIELD_BOOL   = 4
)

type FieldValueUnion struct {
	D *float64 `json:"d,omitempty"`
	I *int64   `json:"i,omitempty"`
	U *uint64  `json:"u,omitempty"`
	S *string  `json:"s,omitempty"`
	B *bool    `json:"b,omitempty"`
}

type Field struct {
	Key   string          `json:"key"`
	Kind  int             `json:"kind"`
	Value FieldValueUnion `json:"value"`
}

type Metric struct {
	Name      string  `json:"name"`
	Tags      []Tag   `json:"tags"`
	Fields    []Field `json:"fields"`
	Timestamp int64   `json:"timestamp"`
}

// Default configurations
var sampleConfig = `
  ## DDS Domain ID configuration
  domain_id = "0"
`

func (d *DDSProducer) SetSerializer(serializer serializers.Serializer) {
	d.serializer = serializer
}

func (d *DDSProducer) Connect() (err error) {

	var xmlString = `
	str://"<dds>
	<qos_library name="QosLibrary">
	<qos_profile name="DefaultProfile" base_name="BuiltinQosLibExp::Generic.StrictReliable" is_default_qos="true">
	</qos_profile>
	</qos_library>
	<types xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="file:////home/kyounghoan/rti_connext_dds-6.0.0/bin/../resource/app/app_support/rtiddsgen/schema/rti_dds_topic_types.xsd">
	<struct name= "Tag">
	  <member name="key" stringMaxLength="255" type="string"/>
	  <member name="value" stringMaxLength="255" type="string"/>
	</struct>
	<enum name="FieldKind">
	  <enumerator name="FIELD_DOUBLE"/>
	  <enumerator name="FIELD_INT"/>
	  <enumerator name="FIELD_UINT"/>
	  <enumerator name="FIELD_STRING"/>
	  <enumerator name="FIELD_BOOL"/>
	</enum>
	<union name="FieldValue">
	<discriminator type="nonBasic" nonBasicTypeName="FieldKind"/>
	<case>
	  <caseDiscriminator value="(FIELD_DOUBLE)"/>
	<member name="d" type="float64"/>
	</case>
	<case>
	  <caseDiscriminator value="(FIELD_INT)"/>
	<member name="i" type="int64"/>
	</case>
	<case>
	  <caseDiscriminator value="(FIELD_UINT)"/>
	<member name="u" type="uint64"/>
	</case>
	<case>
	  <caseDiscriminator value="(FIELD_STRING)"/>
	<member name="s" stringMaxLength="255" type="string"/>
	</case>
	<case>
	  <caseDiscriminator value="(FIELD_BOOL)"/>
	<member name="b" type="boolean"/>
	</case>
	</union>
	<struct name= "Field">
	  <member name="key" stringMaxLength="255" type="string"/>
	  <member name="kind" type="nonBasic"  nonBasicTypeName= "FieldKind"/>
	  <member name="value" type="nonBasic"  nonBasicTypeName= "FieldValue"/>
	</struct>
	<const name="MAX_TAGS" type="int32" value="32"/>
	<const name="MAX_FIELDS" type="int32" value="128"/>
	<struct name= "Metric">
	  <member name="name" stringMaxLength="255" type="string" key="true"/>
	  <member name="tags" sequenceMaxLength="MAX_TAGS" type="nonBasic"  nonBasicTypeName= "Tag"/>
	  <member name="fields" sequenceMaxLength="MAX_FIELDS" type="nonBasic"  nonBasicTypeName= "Field"/>
	  <member name="timestamp" type="int64"/>
	</struct>
	</types>

	<domain_library name="DomainLib">
	<domain name="Telegraf" domain_id="
	` + d.DomainId +
		`">
	<register_type name="Metric" type_ref="Metric"/>
	<topic name="Telegraf" register_type_ref="Metric"/>
	</domain>
	</domain_library>
	<domain_participant_library name="ParticipantLib">
	<domain_participant name="TelegrafParticipant" domain_ref="DomainLib::Telegraf">
	<publisher name="TelegrafPublisher">
	<data_writer name="TelegrafWriter" topic_ref="Telegraf"/>
	</publisher>
	</domain_participant>
	</domain_participant_library>
	</dds>"
	`

	// Create a Connector
	d.connector, err = rti.NewConnector("ParticipantLib::TelegrafParticipant", xmlString)
	if err != nil {
		return err
	}

	// Get a DDS Writer
	d.writer, err = d.connector.GetOutput("TelegrafPublisher::TelegrafWriter")
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
		for _, field := range metric.FieldList() {
			var f Field
			f.Key = field.Key
			switch field.Value.(type) {
			case float64:
				f.Kind = FIELD_DOUBLE
				value := field.Value.(float64)
				f.Value.D = &value
			case int64:
				f.Kind = FIELD_INT
				value := field.Value.(int64)
				f.Value.I = &value
			case uint64:
				f.Kind = FIELD_UINT
				value := field.Value.(uint64)
				f.Value.U = &value
			case string:
				f.Kind = FIELD_STRING
				value := field.Value.(string)
				f.Value.S = &value
			case bool:
				f.Kind = FIELD_BOOL
				value := field.Value.(bool)
				f.Value.B = &value
			default:
			}

			m.Fields = append(m.Fields, f)
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

func init() {
	outputs.Add("dds_producer_lp", func() telegraf.Output {
		return &DDSProducer{}
	})
}
