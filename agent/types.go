package agent

import (

)

type Error struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

type ImageSummaryReturn struct {
	Code int `json:"code"`
}

