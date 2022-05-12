package example

import (
    "fmt"
    "testing"
)

func TestLoadConfig(t *testing.T) {
    InitConfig()
    fmt.Println(AppConfig)
    fmt.Println(AppConfig.GetString("field_name"))
}