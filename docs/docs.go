// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/segment": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Creates new segment with given slug",
                "parameters": [
                    {
                        "description": "Segment input",
                        "name": "segment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.CreateSegmentInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    }
                }
            }
        },
        "/segment/{slug}": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "Deletes segment with given slug",
                "parameters": [
                    {
                        "type": "string",
                        "description": "slug",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    }
                }
            }
        },
        "/segments/user/history/{filename}": {
            "get": {
                "produces": [
                    "text/csv"
                ],
                "summary": "Get user's segments history csv file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "file name",
                        "name": "filename",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    }
                }
            }
        },
        "/segments/user/{user_id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get user's active segments",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Updates user's segments",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User segments input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.UserSegmentsInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    }
                }
            }
        },
        "/segments/user/{user_id}/csv": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get user's segments history link to download",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "date",
                        "example": "\"2023-08-01\"",
                        "description": "From Date",
                        "name": "from",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "date",
                        "example": "\"2023-09-01\"",
                        "description": "To Date",
                        "name": "to",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.OutputError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_unbeman_av-prac-task_internal_model.CreateSegmentInput": {
            "type": "object",
            "properties": {
                "selection": {
                    "type": "number",
                    "example": 0.2
                },
                "slug": {
                    "type": "string",
                    "example": "AVITO_VOICE_MESSAGES"
                }
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.OutputError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "error message"
                }
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.UserSegmentsInput": {
            "type": "object",
            "properties": {
                "segments_to_add": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "PROTECTED_PHONE_NUMBER",
                        "VOICE_MSG"
                    ]
                },
                "segments_to_delete": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "PROMO_5"
                    ]
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Dynamic user segments server",
	Description:      "Avito homework.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
