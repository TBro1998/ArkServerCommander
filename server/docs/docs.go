// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API支持",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/check-init": {
            "get": {
                "description": "检查系统是否已经初始化过用户",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证"
                ],
                "summary": "检查系统初始化状态",
                "responses": {
                    "200": {
                        "description": "初始化状态",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "boolean"
                            }
                        }
                    }
                }
            }
        },
        "/auth/init": {
            "post": {
                "description": "创建第一个管理员用户，只能在系统未初始化时调用",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证"
                ],
                "summary": "初始化系统用户",
                "parameters": [
                    {
                        "description": "用户信息",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "初始化成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "使用用户名和密码登录系统",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "登录凭据",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "认证失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取当前登录用户的基本信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "获取当前用户信息",
                "responses": {
                    "200": {
                        "description": "用户信息",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/models.UserResponse"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/servers": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取当前用户的所有服务器列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "获取服务器列表",
                "responses": {
                    "200": {
                        "description": "服务器列表",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/models.ServerResponse"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "创建一个新的ARK服务器配置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "创建新服务器",
                "parameters": [
                    {
                        "description": "服务器配置",
                        "name": "server",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ServerRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "创建成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/models.ServerResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/servers/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "根据ID获取指定服务器的详细信息（包括配置文件内容）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "获取服务器详情",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "服务器ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "服务器信息（包含配置文件）",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/models.ServerResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "服务器不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "更新指定服务器的配置信息（包括配置文件）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "更新服务器配置",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "服务器ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "更新的服务器配置（可包含配置文件内容）",
                        "name": "server",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ServerUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "更新成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/models.ServerResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "服务器不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除指定的服务器配置（仅允许删除已停止的服务器）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "删除服务器",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "服务器ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "删除成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "服务器不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/servers/{id}/folder": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取指定服务器的文件夹路径和大小（仅供内部管理使用）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "获取服务器文件夹信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "服务器ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文件夹信息",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "服务器不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/servers/{id}/rcon": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取指定服务器的RCON连接信息（包括密码）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "获取服务器RCON信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "服务器ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "RCON信息",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "服务器不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/servers/{id}/start": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "启动指定的ARK服务器",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "启动服务器",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "服务器ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "启动成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "服务器不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/servers/{id}/stop": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "停止指定的ARK服务器",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务器管理"
                ],
                "summary": "停止服务器",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "服务器ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "停止成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "请求错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "服务器不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ServerRequest": {
            "type": "object",
            "required": [
                "admin_password",
                "max_players",
                "name",
                "port",
                "query_port",
                "rcon_port"
            ],
            "properties": {
                "admin_password": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "map": {
                    "type": "string"
                },
                "max_players": {
                    "type": "integer",
                    "maximum": 200,
                    "minimum": 1
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1
                },
                "query_port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1
                },
                "rcon_port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1
                }
            }
        },
        "models.ServerResponse": {
            "type": "object",
            "properties": {
                "admin_password": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "game_ini": {
                    "description": "Game.ini 文件内容",
                    "type": "string"
                },
                "game_user_settings": {
                    "description": "配置文件内容",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "map": {
                    "type": "string"
                },
                "max_players": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "query_port": {
                    "type": "integer"
                },
                "rcon_port": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.ServerUpdateRequest": {
            "type": "object",
            "properties": {
                "admin_password": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "game_ini": {
                    "description": "Game.ini 文件内容",
                    "type": "string"
                },
                "game_user_settings": {
                    "description": "配置文件内容（可选）",
                    "type": "string"
                },
                "map": {
                    "type": "string"
                },
                "max_players": {
                    "type": "integer",
                    "maximum": 200,
                    "minimum": 1
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1
                },
                "query_port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1
                },
                "rcon_port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1
                }
            }
        },
        "models.UserRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.UserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "JWT token in the format \"Bearer {token}\"",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "ARK服务器管理器 API",
	Description:      "基于Gin+Gorm的ARK服务器管理系统API文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
