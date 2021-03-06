package extra

import (
	"github.com/heskandari/jsoner"
	"strings"
	"unicode"
)

// SupportPrivateFields include private fields when encoding/decoding
func SupportPrivateFields(api jsoner.API) {
	api.RegisterExtension(&privateFieldsExtension{})
}

type privateFieldsExtension struct {
	jsoner.DummyExtension
}

func (extension *privateFieldsExtension) UpdateStructDescriptor(structDescriptor *jsoner.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		isPrivate := unicode.IsLower(rune(binding.Field.Name()[0]))
		if isPrivate {
			tag, hastag := binding.Field.Tag().Lookup("json")
			if !hastag {
				binding.FromNames = []string{binding.Field.Name()}
				binding.ToNames = []string{binding.Field.Name()}
				continue
			}
			tagParts := strings.Split(tag, ",")
			names := calcFieldNames(binding.Field.Name(), tagParts[0], tag)
			binding.FromNames = names
			binding.ToNames = names
		}
	}
}

func calcFieldNames(originalFieldName string, tagProvidedFieldName string, wholeTag string) []string {
	// ignore?
	if wholeTag == "-" {
		return []string{}
	}
	// rename?
	var fieldNames []string
	if tagProvidedFieldName == "" {
		fieldNames = []string{originalFieldName}
	} else {
		fieldNames = []string{tagProvidedFieldName}
	}
	// private?
	isNotExported := unicode.IsLower(rune(originalFieldName[0]))
	if isNotExported {
		fieldNames = []string{}
	}
	return fieldNames
}