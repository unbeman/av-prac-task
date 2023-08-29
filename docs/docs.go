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
        "/api/v1/segment": {
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
                        "name": "slug",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.CreateSegment"
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
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Deletes segment with given slug",
                "parameters": [
                    {
                        "description": "Segment input",
                        "name": "slug",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.SegmentInput"
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
        "/api/v1/user/segments": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get user's active segments",
                "parameters": [
                    {
                        "description": "User",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.UserInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.Segment"
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
        }
    },
    "definitions": {
        "github_com_unbeman_av-prac-task_internal_model.CreateSegment": {
            "type": "object",
            "properties": {
                "selection": {
                    "type": "number"
                },
                "slug": {
                    "type": "string"
                }
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.OutputError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.Segment": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "selection": {
                    "description": "0 \u003c user selection \u003c= 1",
                    "type": "number"
                },
                "slug": {
                    "type": "string"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.User"
                    }
                }
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.SegmentInput": {
            "type": "object",
            "properties": {
                "slug": {
                    "type": "string"
                }
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.Segments": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.Segment"
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.User": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "segments": {
                    "$ref": "#/definitions/github_com_unbeman_av-prac-task_internal_model.Segments"
                }
            }
        },
        "github_com_unbeman_av-prac-task_internal_model.UserInput": {
            "type": "object",
            "properties": {
                "user_id": {
                    "type": "integer"
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
                    }
                },
                "segments_to_delete": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "user_id": {
                    "type": "integer"
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
