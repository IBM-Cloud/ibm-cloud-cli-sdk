package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
)

// PluginConfigInvalidTypeError describes the value for a given key in plugin configuration
// was not appropriate for a value of a specific type.
type PluginConfigInvalidTypeError struct {
	Key          string
	ExpectedType string
	Value        interface{}
}

func (e PluginConfigInvalidTypeError) Error() string {
	return fmt.Sprintf("plugin config: invalid type of value for key '%s', expected type %s, but got %T", e.Key, e.ExpectedType, e.Value)
}

// PluginConfig defines methods to access plug-in's private configuration stored
// in a JSON format.
type PluginConfig interface {
	// Get returns the value for a given key.
	// The value can be float64, bool, string, []interface{},
	// map[string]interface or nil if the key not exist.
	Get(key string) interface{}

	// GetWithDefault returns the value for a given key, or defaultVal if the
	// key not exist.
	GetWithDefault(key string, defaultVal interface{}) interface{}

	// GetString returns string value for a given key.
	GetString(key string) (string, error)

	// GetStringWithDefault returns string value for a given key or defaultVal
	// if the key not exist.
	GetStringWithDefault(key string, defaultVal string) (string, error)

	// GetBool returns boolean value for a given key.
	// If the value is a string, it attempts to convert it to bool.
	GetBool(key string) (bool, error)

	// GetBoolWithDefault returns boolean value for a given key or defaultVal if
	// the key not exist.
	GetBoolWithDefault(key string, defaultVal bool) (bool, error)

	// GetInt returns int value for a given key.
	// If the value is float or string, attempts to convert it to int.
	GetInt(key string) (int, error)

	// GetIntWithDefault returns int value for a given key or defaultVal if the
	// key not exist.
	// If the value is float or string, it attempts to convert it to int.
	GetIntWithDefault(key string, defaultVal int) (int, error)

	// GetFloat returns float64 value for a given key.
	// If the value is int or string, it attempts to convert it to float64.
	GetFloat(key string) (float64, error)

	// GetFloatWithDefault returns float64 value for a given key or defaultVal
	// if the key not exist.
	// If the value is int or string, it attempts to convert it to float64.
	GetFloatWithDefault(key string, defaultVal float64) (float64, error)

	// GetStringSlice return string slice for a given key.
	// If the key not exists, return empty string slice.
	GetStringSlice(key string) ([]string, error)

	// GetIntSlice return string slice for a given key.
	GetIntSlice(key string) ([]int, error)

	// GetFloatSlice return string slice for a given key.
	GetFloatSlice(key string) ([]float64, error)

	// GetStringMap return map[string]interface{} for a given key.
	GetStringMap(key string) (map[string]interface{}, error)

	// GetStringMap return map[string]string for a given key.
	GetStringMapString(key string) (map[string]string, error)

	// Exists checks whether the value for a given key exists or not.
	Exists(key string) bool

	// Set sets the value for a given key.
	Set(string, interface{}) error

	// Erase delete a given key.
	Erase(key string) error
}

type pd map[string]interface{}

func (data pd) Marshal() ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func (data pd) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, &data)
}

type pluginConfig struct {
	initOnce  *sync.Once
	lock      sync.RWMutex
	data      pd
	persistor configuration.Persistor
}

func loadPluginConfigFromPath(path string) PluginConfig {
	return &pluginConfig{
		initOnce:  new(sync.Once),
		data:      make(map[string]interface{}),
		persistor: configuration.NewDiskPersistor(path, runtime.GOOS),
	}
}

func (c *pluginConfig) init() {
	c.initOnce.Do(func() {
		err := c.persistor.Load(c.data)
		if err != nil {
			panic(err)
		}
	})
}

func (c *pluginConfig) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c.init()

	return c.data[key]
}

func (c *pluginConfig) GetWithDefault(key string, defaultVal interface{}) interface{} {
	v := c.Get(key)
	if v == nil {
		return defaultVal
	}
	return v
}

func (c *pluginConfig) GetString(key string) (string, error) {
	return c.GetStringWithDefault(key, "")
}

func (c *pluginConfig) GetStringWithDefault(key string, defaultVal string) (string, error) {
	v := c.Get(key)
	if v == nil {
		return defaultVal, nil
	}
	if s, ok := toString(v); ok {
		return s, nil
	}
	return defaultVal, PluginConfigInvalidTypeError{Key: key, ExpectedType: "string", Value: v}
}

func toString(v interface{}) (string, bool) {
	switch v.(type) {
	case bool:
		return strconv.FormatBool(v.(bool)), true
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64), true
	case string:
		return v.(string), true
	case nil:
		return "", true
	}
	return "", false
}

func (c *pluginConfig) GetBool(key string) (bool, error) {
	return c.GetBoolWithDefault(key, false)
}

func (c *pluginConfig) GetBoolWithDefault(key string, defaultVal bool) (bool, error) {
	v := c.Get(key)
	if v == nil {
		return defaultVal, nil
	}
	if b, ok := toBool(v); ok {
		return b, nil
	}
	return defaultVal, PluginConfigInvalidTypeError{Key: key, ExpectedType: "bool", Value: v}
}

func toBool(v interface{}) (bool, bool) {
	switch v.(type) {
	case bool:
		return v.(bool), true
	case string:
		b, err := strconv.ParseBool(v.(string))
		if err == nil {
			return b, true
		}
	}
	return false, false
}

func (c *pluginConfig) GetInt(key string) (int, error) {
	return c.GetIntWithDefault(key, 0)
}

func (c *pluginConfig) GetIntWithDefault(key string, defaultVal int) (int, error) {
	v := c.Get(key)
	if v == nil {
		return defaultVal, nil
	}
	if i, ok := toInt(v); ok {
		return i, nil
	}
	return defaultVal, PluginConfigInvalidTypeError{Key: key, ExpectedType: "int", Value: v}
}

func toInt(v interface{}) (int, bool) {
	switch v.(type) {
	case float64:
		return int(v.(float64)), true
	case string:
		i, err := strconv.ParseInt(v.(string), 0, 0)
		if err == nil {
			return int(i), true
		}
	}
	return 0, false
}

func (c *pluginConfig) GetFloat(key string) (float64, error) {
	return c.GetFloatWithDefault(key, 0.00)
}

func (c *pluginConfig) GetFloatWithDefault(key string, defaultVal float64) (float64, error) {
	v := c.Get(key)
	if v == nil {
		return defaultVal, nil
	}
	if f, ok := toFloat(v); ok {
		return f, nil
	}
	return defaultVal, PluginConfigInvalidTypeError{Key: key, ExpectedType: "float", Value: v}
}

func toFloat(v interface{}) (float64, bool) {
	switch v.(type) {
	case float64:
		return v.(float64), true
	case string:
		f, err := strconv.ParseFloat(v.(string), 64)
		if err == nil {
			return f, true
		}
	}
	return 0.00, false
}

func (c *pluginConfig) GetSlice(key string) ([]interface{}, error) {
	v := c.Get(key)
	if v == nil {
		return []interface{}{}, nil
	}
	if s, ok := v.([]interface{}); ok {
		return s, nil
	}
	return []interface{}{}, PluginConfigInvalidTypeError{Key: key, ExpectedType: "[]interface{}", Value: v}
}

func (c *pluginConfig) GetStringSlice(key string) ([]string, error) {
	v := c.Get(key)
	if v == nil {
		return []string{}, nil
	}
	if ss, ok := toStringSlice(v); ok {
		return ss, nil
	}
	return []string{}, PluginConfigInvalidTypeError{Key: key, ExpectedType: "[]string", Value: v}
}

func toStringSlice(v interface{}) ([]string, bool) {
	s, ok := v.([]interface{})
	if !ok {
		return []string{}, false
	}

	ret := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		ii, ok := toString(s[i])
		if !ok {
			return []string{}, false
		}
		ret[i] = ii
	}
	return ret, true
}

func (c *pluginConfig) GetIntSlice(key string) ([]int, error) {
	v := c.Get(key)
	if v == nil {
		return []int{}, nil
	}
	if is, ok := toIntSlice(v); ok {
		return is, nil
	}
	return []int{}, PluginConfigInvalidTypeError{Key: key, ExpectedType: "[]int", Value: v}
}

func toIntSlice(v interface{}) ([]int, bool) {
	s, ok := v.([]interface{})
	if !ok {
		return []int{}, false
	}

	is := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		ii, ok := toInt(s[i])
		if !ok {
			return []int{}, false
		}
		is[i] = ii
	}
	return is, true
}

func (c *pluginConfig) GetFloatSlice(key string) ([]float64, error) {
	v := c.Get(key)
	if v == nil {
		return []float64{}, nil
	}
	if fs, ok := toFloatSlice(v); ok {
		return fs, nil
	}
	return []float64{}, PluginConfigInvalidTypeError{Key: key, ExpectedType: "[]float64", Value: v}
}

func toFloatSlice(v interface{}) ([]float64, bool) {
	s, ok := v.([]interface{})
	if !ok {
		return []float64{}, false
	}

	fs := make([]float64, len(s))
	for i := 0; i < len(s); i++ {
		f, ok := toFloat(s[i])
		if !ok {
			return []float64{}, false
		}
		fs[i] = f
	}
	return fs, true
}

func (c *pluginConfig) GetStringMap(key string) (map[string]interface{}, error) {
	v := c.Get(key)
	if v == nil {
		return map[string]interface{}{}, nil
	}
	if sm, ok := v.(map[string]interface{}); ok {
		return sm, nil
	}
	return map[string]interface{}{}, PluginConfigInvalidTypeError{Key: key, ExpectedType: "map[string]interface{}", Value: v}
}

func (c *pluginConfig) GetStringMapString(key string) (map[string]string, error) {
	v := c.Get(key)
	if v == nil {
		return map[string]string{}, nil
	}
	if sms, ok := toMapStringMapString(v); ok {
		return sms, nil
	}
	return map[string]string{}, PluginConfigInvalidTypeError{Key: key, ExpectedType: "map[string]string", Value: v}
}

func toMapStringMapString(v interface{}) (map[string]string, bool) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return map[string]string{}, false
	}

	sms := make(map[string]string)
	for k, v := range m {
		s, ok := toString(v)
		if !ok {
			return map[string]string{}, false
		}
		sms[k] = s
	}

	return sms, true
}

func (c *pluginConfig) Exists(key string) bool {
	return c.Get(key) != nil
}

func (c *pluginConfig) Set(key string, v interface{}) error {
	return c.write(func() {
		c.data[key] = v
	})
}

func (c *pluginConfig) Erase(key string) error {
	return c.write(func() {
		delete(c.data, key)
	})
}

func (c *pluginConfig) write(cb func()) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.init()

	cb()

	err := c.persistor.Save(c.data)
	if err != nil {
		err = errors.New(T("Unable to save plugin config: ") + err.Error())
	}
	return err
}
