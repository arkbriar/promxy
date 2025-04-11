package sigv4

const secretToken = "<secret>"

// Secret special type for storing secrets.
type Secret string

// MarshalSecretValue if set to true will expose Secret type
// through the marshal interfaces. Useful for outside projects
// that load and marshal the Prometheus config.
var MarshalSecretValue bool = false

// MarshalYAML implements the yaml.Marshaler interface for Secrets.
func (s Secret) MarshalYAML() (interface{}, error) {
	if MarshalSecretValue {
		return string(s), nil
	}
	if s != "" {
		return secretToken, nil
	}
	return nil, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Secrets.
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Secret
	return unmarshal((*plain)(s))
}
