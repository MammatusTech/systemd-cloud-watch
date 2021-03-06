package cloud_watch

import (
	"io/ioutil"
	"github.com/hashicorp/hcl"
)

type Config struct {
	AWSRegion     string `hcl:"aws_region"`
	EC2InstanceId string `hcl:"ec2_instance_id"`
	LogGroupName  string `hcl:"log_group"`
	LogStreamName string `hcl:"log_stream"`
	LogPriority   string `hcl:"log_priority"`
	JournalDir    string `hcl:"journal_dir"`
	BufferSize    int    `hcl:"buffer_size"`
	Debug         bool    `hcl:"debug"`
	Local         bool    `hcl:"local"`
	AllowedFields []string `hcl:"fields"`
	OmitFields    []string `hcl:"omit_fields"`
	logPriority   int
	fields        map[string]struct{}
	omitFields    map[string]struct{}
	FieldLength   int    `hcl:"field_length"`
}

func (config *Config) GetJournalDLogPriority() (Priority) {

	logLevels := map[Priority][]string{
		EMERGENCY: {"0", "emerg"},
		ALERT:     {"1", "alert"},
		CRITICAL:  {"2", "crit"},
		ERROR:     {"3", "err"},
		WARNING:   {"4", "warning"},
		NOTICE:    {"5", "notice"},
		INFO:      {"6", "info"},
		DEBUG:     {"7", "debug"},
	}

	for i, s := range logLevels {
		if s[0] == config.LogPriority || s[1] == config.LogPriority {
			return i
		}
	}

	return DEBUG
}

func (config *Config) AllowField(fieldName string) bool {

	if len(config.AllowedFields) == 0 && len(config.OmitFields) == 0 {
		return true
	} else if len(config.AllowedFields) > 0 && len(config.OmitFields) == 0 {
		_, hasField := config.fields[fieldName]
		return hasField
	} else if len(config.AllowedFields) == 0 && len(config.OmitFields) > 0 {
		_, omitField := config.omitFields[fieldName]
		return !omitField
	} else {
		logger := NewSimpleLogger("allow-field", config)
		logger.Warning.Println("Only fields or omit_fields should be set")
		_, omitField := config.omitFields[fieldName]
		if omitField {
			return !omitField
		} else {
			_, hasField := config.fields[fieldName]
			return hasField

		}
	}
}

func arrayToMap(array []string) map[string]struct{} {
	theMap := make(map[string]struct{})
	if array != nil && len(array) > 0 {
		for _, element := range array {
			theMap[element] = struct{}{}
		}
	}
	return theMap
}

func LoadConfigFromString(data string, logger *Logger) (*Config, error) {

	if logger == nil {
		logger = NewSimpleLogger("config", nil)
	}
	config := &Config{}
	logger.Debug.Println("Loading log...")
	err := hcl.Decode(&config, data)
	if err != nil {
		return nil, err
	}
	config.fields = arrayToMap(config.AllowedFields)
	config.omitFields = arrayToMap(config.OmitFields)

	if config.BufferSize == 0 {
		logger.Debug.Println("Loading log... BufferSize not set, setting to 10")
		config.BufferSize = 10
	}


	if config.LogPriority == "" {
		logger.Debug.Println("Loading log... LogPriority not set, setting to debug")
		config.LogPriority = "debug"
	}

	return config, nil

}
func LoadConfig(filename string, logger *Logger) (*Config, error) {
	logger.Info.Printf("Loading config %s", filename)

	configBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadConfigFromString(string(configBytes), logger)
}