package csgopoolweb

import (
  "strconv"
  "fmt"
)

func ParseBool(str string) bool {
	if str == "true" {
		return true
	}
	return false
}

func ParseFloat(str string) float32 {
    _nb, _ := strconv.ParseFloat(str, 32)
    return float32(_nb)
}

func BoolToString(val bool) string {
  if val {
   return "true" 
  }
  return "false"
}

func FloatToString(val float32) string {
  return fmt.Sprintf("%.2f", val)
}