package configoration

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Section provides functionalities to access
// sections and values in a Section by a key.
//
// A key can be a value or section key itself
// like "webserver" or it can span over sections
// like "general:webserver". In this case, the
// value after the last delimiter is selected
// section or value.
type Section interface {
	// GetSection returns a section by key.
	// If the desired section is not existent,
	// the returned value will be nil.
	GetSection(key string) Section

	// GetValue returns an interface value by
	// key. If the desired value could not be
	// found, nil and ErrNil is returned.
	GetValue(key string) (interface{}, error)

	// GetString is shorthand for GetValue and
	// returns a string or an ErrNil if the
	// key was not found.
	//
	// If the value selected is not a string,
	// ErrInvalidType will be returned.
	GetString(key string) (string, error)

	// GetInt is shorthand for GetValue and
	// returns an int or an ErrNil if the
	// key was not found.
	//
	// If the value selected is not an int,
	// ErrInvalidType will be returned.
	GetInt(key string) (int, error)

	// GetBool is shorthand for GetValue and
	// returns a bool or an ErrNil if the
	// key was not found.
	//
	// If the value selected is not a bool,
	// ErrInvalidType will be returned.
	GetBool(key string) (bool, error)

	// GetFloat64 is shorthand for GetValue and
	// returns a float64 or an ErrNil if the
	// key was not found.
	//
	// If the value selected is not a float64,
	// ErrInvalidType will be returned.
	GetFloat64(key string) (float64, error)

	// GetValueOrDef returns an interface value
	// by key. If the desired value could not be
	// found, def will be returned.
	GetValueOrDef(key string, def interface{}) interface{}

	// GetStringOrDef is shorthand for GetValueOrDef
	// and returns a string which is eather the
	// found value or the vlaue of def.
	GetStringOrDef(key string, def string) string

	// GetIntOrDef is shorthand for GetValueOrDef
	// and returns an int which is eather the
	// found value or the vlaue of def.
	GetIntOrDef(key string, def int) int

	// GetBoolOrDef is shorthand for GetValueOrDef
	// and returns a bool which is eather the
	// found value or the vlaue of def.
	GetBoolOrDef(key string, def bool) interface{}

	// GetFloat64OrDef is shorthand for GetValueOrDef
	// and returns a float64 which is eather the
	// found value or the vlaue of def.
	GetFloat64OrDef(key string, def float64) float64

	// IsNil returns true if the current section
	// instance is nil.
	IsNil() bool
}

// section is the default implementation of
// the Section interface.
type section struct {
	mtx sync.Mutex
	m   ConfigMap
}

func (s *section) GetSection(key string) Section {
	for _, nextSelector := range splitSections(key) {
		if s == nil {
			return nil
		}
		s = s.getSection(nextSelector)
	}
	return s
}

func (s *section) GetValue(key string) (interface{}, error) {
	if s == nil {
		return nil, ErrNil
	}

	selectors := splitSections(key)
	lenSelectors := len(selectors)
	if lenSelectors > 1 {
		for i := 0; i < lenSelectors-1; i++ {
			s = s.getSection(selectors[i])
			if s == nil {
				return nil, ErrNil
			}
		}
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	v, ok := s.m[selectors[lenSelectors-1]]
	if !ok {
		return nil, ErrNil
	}

	return v, nil
}

func (s *section) GetString(key string) (string, error) {
	v, err := s.GetValue(key)
	if err != nil {
		return "", err
	}

	vt, ok := v.(string)
	if !ok {
		vt = valToString(v)
	}

	return vt, nil
}

func (s *section) GetInt(key string) (int, error) {
	v, err := s.GetValue(key)
	if err != nil {
		return 0, err
	}

	vt, ok := v.(int)
	if !ok {
		vt, err = strconv.Atoi(valToString(v))
	}

	return vt, err
}

func (s *section) GetBool(key string) (bool, error) {
	v, err := s.GetValue(key)
	if err != nil {
		return false, err
	}

	vt, ok := v.(bool)
	if !ok {
		vt, err = strconv.ParseBool(valToString(v))
	}

	return vt, err
}

func (s *section) GetFloat64(key string) (float64, error) {
	v, err := s.GetValue(key)
	if err != nil {
		return 0, err
	}

	vt, ok := v.(float64)
	if !ok {
		vt, err = strconv.ParseFloat(valToString(v), 64)
	}

	return vt, err
}

func (s *section) GetValueOrDef(key string, def interface{}) interface{} {
	v, err := s.GetValue(key)
	if err != nil {
		v = def
	}
	return v
}

func (s *section) GetStringOrDef(key string, def string) string {
	v, err := s.GetString(key)
	if err != nil {
		v = def
	}
	return v
}

func (s *section) GetIntOrDef(key string, def int) int {
	v, err := s.GetInt(key)
	if err != nil {
		v = def
	}
	return v
}

func (s *section) GetBoolOrDef(key string, def bool) interface{} {
	v, err := s.GetBool(key)
	if err != nil {
		v = def
	}
	return v
}

func (s *section) GetFloat64OrDef(key string, def float64) float64 {
	v, err := s.GetFloat64(key)
	if err != nil {
		v = def
	}
	return v
}

func (s *section) IsNil() bool {
	return s == nil
}

// getSection returns the desired section
// or nil, if not found.
func (s *section) getSection(sec string) *section {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	v := s.m[sec]
	vc, ok := v.(ConfigMap)
	if !ok {
		return nil
	}

	return &section{
		mtx: sync.Mutex{},
		m:   vc,
	}
}

// splitSections splits the passed key by
// the Delimiter and returns the resulting
// array of strings.
func splitSections(key string) []string {
	return strings.Split(key, Delimiter)
}

// valToString returns the passed interface
// as a string using fmt.Sprintf("%v", v) as
// converter.
func valToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
