package utilities

import (
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/linkedin/goavro/v2"
)

const avroSchemaSource = `{
  "type": "record",
  "name": "SQLCommandEvent",
  "fields": [
    {
      "name": "sql_command",
      "type": "string",
      "doc": "The intercepted SQL command."
    },
    {
      "name": "database",
      "type": "string",
      "doc": "The name of the database."
    },
    {
      "name": "timestamp",
      "type": "string",
      "doc": "The time when the SQL command was intercepted."
    },
    {
      "name": "user",
      "type": "string",
      "doc": "The username executing the SQL command."
    }
  ]
}`

func CreateAvroMessage(s interface{}, encoding pubsub.SchemaEncoding) ([]byte, error) {
	codec, err := goavro.NewCodec(avroSchemaSource)
	if err != nil {
		return nil, err
	}

	data, err := StructToMap(s)
	if err != nil {
		return nil, err
	}

	var binaryData []byte

	switch encoding {

	case pubsub.EncodingJSON:
		binaryData, err = codec.TextualFromNative(nil, data)
		if err != nil {
			return nil, err
		}
	case pubsub.EncodingBinary:
		binaryData, err = codec.BinaryFromNative(nil, data)
		if err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("unsupported encoding: %v", encoding)
		return nil, err
	}

	return binaryData, nil
}
