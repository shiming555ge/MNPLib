package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PasskeyRequest struct {
	Description string `json:"description"`
	Operator    string `json:"operator"`
	IsActive    bool   `json:"is_active"`
}

func main() {
	// 测试创建passkey的API
