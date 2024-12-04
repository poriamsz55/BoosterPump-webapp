package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Get by parameter name.
// Parameter could be sent by POST and PUT body parameters,
// or URL query string values.
func Uint32(c echo.Context, paramName string) (uint32, error) {
	valStr := c.FormValue(paramName)

	// If key is not present, FormValue returns the empty string.
	if len(valStr) == 0 {
		return 0, errors.New("expected parameter is not set: " + paramName)
	}

	val, err := strconv.ParseUint(valStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parsing uint32 failed: %w, value: %s", err, valStr)
	}

	return uint32(val), nil
}

// Get by parameter name.
// Parameter could be sent by POST and PUT body parameters,
// or URL query string values.
func Uint64(c echo.Context, paramName string) (uint64, error) {
	valStr := c.FormValue(paramName)

	// If key is not present, FormValue returns the empty string.
	if len(valStr) == 0 {
		return 0, errors.New("expected parameter is not set: " + paramName)
	}

	val, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing uint32 failed: %w, value: %s", err, valStr)
	}

	return val, nil
}

// Get by parameter name.
// Parameter could be sent by POST and PUT body parameters,
// or URL query string values.
func Int(c echo.Context, paramName string) (int, error) {
	valStr := c.FormValue(paramName)

	// If key is not present, FormValue returns the empty string.
	if len(valStr) == 0 {
		return 0, errors.New("expected parameter is not set: " + paramName)
	}

	val, err := strconv.ParseInt(valStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

func Float32(c echo.Context, paramName string) (float32, error) {
	valStr := c.FormValue(paramName)

	// If key is not present, FormValue returns the empty string.
	if len(valStr) == 0 {
		return 0, errors.New("expected parameter is not set: " + paramName + " : " + valStr)
	}

	val, err := strconv.ParseFloat(valStr, 32)
	if err != nil {
		return 0, err
	}

	return float32(val), nil
}

func Float64(c echo.Context, paramName string) (float64, error) {
	valStr := c.FormValue(paramName)

	// If key is not present, FormValue returns the empty string.
	if len(valStr) == 0 {
		return 0, errors.New("expected parameter is not set: " + paramName + " : " + valStr)
	}

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0, err
	}

	return float64(val), nil
}


func AsBytes(c echo.Context, param string) ([]byte, error) {
	f, err := c.FormFile(param)
	if err != nil {
		return []byte{}, fmt.Errorf("reading upload as bytes: %w", err)
	}
	src, err := f.Open()
	if err != nil {
		return []byte{}, fmt.Errorf("open upload as bytes: %w", err)
	}
	defer src.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, src)
	return buf.Bytes(), err
}

func AsJSONFile(c echo.Context, param string, output interface{}) error {
	d, err := AsBytes(c, param)
	if err != nil {
		return fmt.Errorf("parsing json file error: %w", err)
	}
	err = json.Unmarshal(d, output)
	if err != nil {
		return fmt.Errorf("unmarshal json file error: %w", err)
	}
	return nil
}

func AsJSON(c echo.Context, param string, output interface{}) error {
	d := c.FormValue(param)
	err := json.Unmarshal([]byte(d), output)
	if err != nil {
		return fmt.Errorf("unmarshal json file error: %w", err)
	}
	return nil
}
