package config

const (
	DefaultLength    = 12
	DefaultDelimiter = "|"
	DefaultUniq      = true
)

// Format is a struct that contains the format config for the workspace name.
type Format struct {
	Length    int
	Delimiter string
	Uniq      bool
}

// DefaultFormat returns the default format config.
func DefaultFormat() *Format {
	return &Format{Length: DefaultLength, Delimiter: DefaultDelimiter, Uniq: DefaultUniq}
}
