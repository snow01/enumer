package main

// Arguments to format are:
//	[1]: type name
const stringNameToValueMethod = `// %[1]sString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func %[1]sString(s string) (val %[1]s, err error) {
	var ok = false
	if val, ok = _%[1]sNameToValueMap[s]; !ok {
		err = fmt.Errorf("%%s does not belong to %[1]s values", s)
	}
	
	return
}
`

// Arguments to format are:
//	[1]: type name
const stringNameToValueMethodWithUnknownSupport = `// %[1]sString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func %[1]sString(s string) (val %[1]s, err error) {
	var ok = false
	if val, ok = _%[1]sNameToValueMap[s]; ok {
		return
	}
    
    if val, ok = _%[1]sNameToValueMap[%[2]q]; !ok {
		err = fmt.Errorf("%%s does not belong to %[1]s values", s)
	}

	return
}
`

// Arguments to format are:
//	[1]: type name
const stringValuesMethod = `// %[1]sValues returns all values of the enum
func %[1]sValues() []%[1]s {
	return _%[1]sValues
}
`

// Arguments to format are:
//	[1]: type name
const stringBelongsMethodLoop = `// IsA%[1]s returns "true" if the value is listed in the enum definition. "false" otherwise
func (i %[1]s) IsA%[1]s() bool {
	for _, v := range _%[1]sValues {
		if i == v {
			return true
		}
	}
	return false
}
`

// Arguments to format are:
//	[1]: type name
const stringBelongsMethodSet = `// IsA%[1]s returns "true" if the value is listed in the enum definition. "false" otherwise
func (i %[1]s) IsA%[1]s() bool {
	_, ok := _%[1]sMap[i] 
	return ok
}
`

func (g *Generator) buildBasicExtras(runs [][]Value, values []Value, typeName string, runsThreshold int, unknown string) {
	// At this moment, either "g.declareIndexAndNameVars()" or "g.declareNameVars()" has been called

	// Print the slice of values
	g.Printf("\nvar _%sValues = []%s{", typeName, typeName)
	for _, value := range values {
		g.Printf("\t%s, ", value.str)
	}
	g.Printf("}\n\n")

	// Print the map between name and value
	g.Printf("\nvar _%sNameToValueMap = map[string]%s{\n", typeName, typeName)
	for _, value := range values {
		for _, d := range value.decodes {
			g.Printf("\t\"%s\": %s,\n", d, &value)
		}
	}
	g.Printf("}\n\n")

	// Print the basic extra methods
	if unknown == "" {
		g.Printf(stringNameToValueMethod, typeName)
	} else {
		g.Printf(stringNameToValueMethodWithUnknownSupport, typeName, unknown)
	}

	g.Printf(stringValuesMethod, typeName)
	if len(runs) < runsThreshold {
		g.Printf(stringBelongsMethodLoop, typeName)
	} else { // There is a map of values, the code is simpler then
		g.Printf(stringBelongsMethodSet, typeName)
	}
}

// Arguments to format are:
//	[1]: type name
const jsonMethods = `
// MarshalJSON implements the json.Marshaler interface for %[1]s
func (i %[1]s) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for %[1]s
func (i *%[1]s) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("%[1]s should be a string, got %%s", data)
	}

	var err error
	*i, err = %[1]sString(s)
	return err
}
`

func (g *Generator) buildJSONMethods(typeName string) {
	g.Printf(jsonMethods, typeName)
}

// Arguments to format are:
//	[1]: type name
const textMethods = `
// MarshalText implements the encoding.TextMarshaler interface for %[1]s
func (i %[1]s) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for %[1]s
func (i *%[1]s) UnmarshalText(text []byte) error {
	var err error
	*i, err = %[1]sString(string(text))
	return err
}
`

func (g *Generator) buildTextMethods(typeName string) {
	g.Printf(textMethods, typeName)
}

// Arguments to format are:
//	[1]: type name
const yamlMethods = `
// MarshalYAML implements a YAML Marshaler for %[1]s
func (i %[1]s) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for %[1]s
func (i *%[1]s) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = %[1]sString(s)
	return err
}
`

func (g *Generator) buildYAMLMethods(typeName string) {
	g.Printf(yamlMethods, typeName)
}
