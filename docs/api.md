# Survey System API Documentation

## 概述

目标问卷系统提供 RESTful API 服务，支持问卷管理、题目配置、加密链接生成、问卷填答等功能。API 采用 JSON 格式进行数据交换，使用 JWT 进行身份认证。

**Base URL**: `http://localhost:8080/api/v1`

**API Version**: v1

## 认证说明

### JWT 认证

大部分管理端 API 需要 JWT 认证。在请求头中包含 `Authorization` 字段：

```
Authorization: Bearer <your_jwt_token>
```

### 获取 Token

通过登录接口获取 JWT token：

```bash
POST /api/v1/auth/login
```

Token 有效期为 24 小时，过期后需要重新登录。

## 通用响应格式

### 成功响应

```json
{
    "success": true,
    "data": { ... }
}
```

### 分页响应

```json
{
    "success": true,
    "data": [ ... ],
    "meta": {
        "page": 1,
        "page_size": 20,
        "total": 100
    }
}
```


### 错误响应

```json
{
    "success": false,
    "error": {
        "code": "ERROR_CODE",
        "message": "错误描述信息"
    }
}
```

## 错误码说明

| 错误码 | HTTP 状态码 | 说明 |
|--------|------------|------|
| `UNAUTHORIZED` | 401 | 未授权访问，需要登录 |
| `FORBIDDEN` | 403 | 禁止访问，权限不足 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `INVALID_TOKEN` | 400 | 无效的令牌 |
| `TOKEN_EXPIRED` | 403 | 令牌已过期 |
| `LINK_USED` | 403 | 链接已被使用 |
| `VALIDATION_FAILED` | 400 | 数据验证失败 |
| `SURVEY_NOT_PUBLISHED` | 400 | 问卷未发布 |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 |

## 分页参数

支持分页的接口接受以下查询参数：

- `page`: 页码，从 1 开始，默认为 1
- `page_size`: 每页数量，默认为 20，最大为 100

示例：`GET /api/v1/surveys?page=2&page_size=10`

---

# API 端点

## 1. 认证相关接口

### 1.1 用户登录

**端点**: `POST /api/v1/auth/login`

**认证**: 不需要

**描述**: 用户登录并获取 JWT token

**请求体**:

```json
{
    "username": "admin",
    "password": "password123"
}
```

**请求参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，3-50 字符 |
| password | string | 是 | 密码，至少 6 字符 |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@example.com",
            "role": "admin"
        }
    }
}
```

**错误响应**:

- 401 Unauthorized: 用户名或密码错误

```json
{
    "success": false,
    "error": {
        "code": "UNAUTHORIZED",
        "message": "用户名或密码错误"
    }
}
```

**cURL 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

---

## 2. 问卷管理接口

### 2.1 创建问卷

**端点**: `POST /api/v1/surveys`

**认证**: 需要 JWT

**描述**: 创建一个新的问卷（草稿状态）

**请求体**:

```json
{
    "title": "客户满意度调查",
    "description": "本问卷旨在了解客户对我们服务的满意度"
}
```

**请求参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 问卷标题，最多 200 字符 |
| description | string | 否 | 问卷描述，最多 5000 字符 |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "id": 1,
        "user_id": 1,
        "title": "客户满意度调查",
        "description": "本问卷旨在了解客户对我们服务的满意度",
        "status": "draft",
        "created_at": "2025-10-25T10:00:00Z",
        "updated_at": "2025-10-25T10:00:00Z"
    }
}
```

**cURL 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/surveys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "客户满意度调查",
    "description": "本问卷旨在了解客户对我们服务的满意度"
  }'
```

### 2.2 更新问卷

**端点**: `PUT /api/v1/surveys/:id`

**认证**: 需要 JWT

**描述**: 更新指定问卷的信息

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**请求体**:

```json
{
    "title": "客户满意度调查（更新版）",
    "description": "更新后的描述"
}
```

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "id": 1,
        "user_id": 1,
        "title": "客户满意度调查（更新版）",
        "description": "更新后的描述",
        "status": "draft",
        "created_at": "2025-10-25T10:00:00Z",
        "updated_at": "2025-10-25T10:30:00Z"
    }
}
```

**cURL 示例**:

```bash
curl -X PUT http://localhost:8080/api/v1/surveys/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "客户满意度调查（更新版）",
    "description": "更新后的描述"
  }'
```

### 2.3 删除问卷

**端点**: `DELETE /api/v1/surveys/:id`

**认证**: 需要 JWT

**描述**: 删除指定问卷及其所有题目和填答记录

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": null
}
```

**cURL 示例**:

```bash
curl -X DELETE http://localhost:8080/api/v1/surveys/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 2.4 查询问卷详情

**端点**: `GET /api/v1/surveys/:id`

**认证**: 需要 JWT

**描述**: 获取指定问卷的详细信息，包含所有题目

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "id": 1,
        "user_id": 1,
        "title": "客户满意度调查",
        "description": "本问卷旨在了解客户对我们服务的满意度",
        "status": "published",
        "created_at": "2025-10-25T10:00:00Z",
        "updated_at": "2025-10-25T10:30:00Z",
        "questions": [
            {
                "id": 1,
                "survey_id": 1,
                "type": "text",
                "title": "您的姓名",
                "description": "",
                "required": true,
                "order": 1,
                "config": {},
                "prefill_key": "name",
                "created_at": "2025-10-25T10:05:00Z",
                "updated_at": "2025-10-25T10:05:00Z"
            },
            {
                "id": 2,
                "survey_id": 1,
                "type": "single",
                "title": "您对我们的服务满意吗？",
                "description": "",
                "required": true,
                "order": 2,
                "config": {
                    "options": ["非常满意", "满意", "一般", "不满意"]
                },
                "prefill_key": "",
                "created_at": "2025-10-25T10:06:00Z",
                "updated_at": "2025-10-25T10:06:00Z"
            }
        ]
    }
}
```

**cURL 示例**:

```bash
curl -X GET http://localhost:8080/api/v1/surveys/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 2.5 查询问卷列表

**端点**: `GET /api/v1/surveys`

**认证**: 需要 JWT

**描述**: 获取当前用户的所有问卷列表（支持分页）

**查询参数**:

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | integer | 否 | 1 | 页码 |
| page_size | integer | 否 | 20 | 每页数量（最大 100） |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "user_id": 1,
            "title": "客户满意度调查",
            "description": "本问卷旨在了解客户对我们服务的满意度",
            "status": "published",
            "created_at": "2025-10-25T10:00:00Z",
            "updated_at": "2025-10-25T10:30:00Z"
        },
        {
            "id": 2,
            "user_id": 1,
            "title": "员工反馈问卷",
            "description": "",
            "status": "draft",
            "created_at": "2025-10-25T11:00:00Z",
            "updated_at": "2025-10-25T11:00:00Z"
        }
    ],
    "meta": {
        "page": 1,
        "page_size": 20,
        "total": 2
    }
}
```

**cURL 示例**:

```bash
curl -X GET "http://localhost:8080/api/v1/surveys?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 2.6 发布问卷

**端点**: `POST /api/v1/surveys/:id/publish`

**认证**: 需要 JWT

**描述**: 将问卷状态从草稿改为已发布

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": null
}
```

**cURL 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/surveys/1/publish \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 3. 题目管理接口

### 3.1 创建题目

**端点**: `POST /api/v1/questions`

**认证**: 需要 JWT

**描述**: 为问卷创建一个新题目

**请求体**:

```json
{
    "survey_id": 1,
    "type": "single",
    "title": "您对我们的服务满意吗？",
    "description": "请选择最符合您感受的选项",
    "required": true,
    "order": 1,
    "config": {
        "options": ["非常满意", "满意", "一般", "不满意"]
    },
    "prefill_key": ""
}
```

**请求参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| survey_id | integer | 是 | 问卷 ID |
| type | string | 是 | 题目类型：text, single, multiple, table |
| title | string | 是 | 题目标题，最多 500 字符 |
| description | string | 否 | 题目描述，最多 5000 字符 |
| required | boolean | 否 | 是否必填，默认 false |
| order | integer | 是 | 显示顺序，从 0 开始 |
| config | object | 否 | 题目配置（根据类型不同） |
| prefill_key | string | 否 | 预填字段键名，最多 100 字符 |

**题目配置说明**:

**单选题/多选题 (single/multiple)**:
```json
{
    "options": ["选项1", "选项2", "选项3"]
}
```

**表格题 (table)**:
```json
{
    "columns": [
        {
            "id": "col1",
            "type": "text",
            "label": "姓名"
        },
        {
            "id": "col2",
            "type": "number",
            "label": "年龄"
        },
        {
            "id": "col3",
            "type": "select",
            "label": "性别",
            "options": ["男", "女"]
        }
    ],
    "min_rows": 1,
    "max_rows": 10,
    "can_add_row": true
}
```

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "id": 1,
        "survey_id": 1,
        "type": "single",
        "title": "您对我们的服务满意吗？",
        "description": "请选择最符合您感受的选项",
        "required": true,
        "order": 1,
        "config": {
            "options": ["非常满意", "满意", "一般", "不满意"]
        },
        "prefill_key": "",
        "created_at": "2025-10-25T10:05:00Z",
        "updated_at": "2025-10-25T10:05:00Z"
    }
}
```

**cURL 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/questions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "survey_id": 1,
    "type": "single",
    "title": "您对我们的服务满意吗？",
    "description": "请选择最符合您感受的选项",
    "required": true,
    "order": 1,
    "config": {
        "options": ["非常满意", "满意", "一般", "不满意"]
    }
  }'
```

### 3.2 更新题目

**端点**: `PUT /api/v1/questions/:id`

**认证**: 需要 JWT

**描述**: 更新指定题目的信息

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 题目 ID |

**请求体**: 与创建题目相同（不包含 survey_id）

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "id": 1,
        "survey_id": 1,
        "type": "single",
        "title": "您对我们的服务满意吗？（更新版）",
        "description": "请选择最符合您感受的选项",
        "required": true,
        "order": 1,
        "config": {
            "options": ["非常满意", "满意", "一般", "不满意", "非常不满意"]
        },
        "prefill_key": "",
        "created_at": "2025-10-25T10:05:00Z",
        "updated_at": "2025-10-25T10:35:00Z"
    }
}
```

**cURL 示例**:

```bash
curl -X PUT http://localhost:8080/api/v1/questions/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "type": "single",
    "title": "您对我们的服务满意吗？（更新版）",
    "description": "请选择最符合您感受的选项",
    "required": true,
    "order": 1,
    "config": {
        "options": ["非常满意", "满意", "一般", "不满意", "非常不满意"]
    }
  }'
```

### 3.3 删除题目

**端点**: `DELETE /api/v1/questions/:id`

**认证**: 需要 JWT

**描述**: 删除指定题目

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 题目 ID |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": null
}
```

**cURL 示例**:

```bash
curl -X DELETE http://localhost:8080/api/v1/questions/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3.4 重新排序题目

**端点**: `PUT /api/v1/surveys/:id/questions/reorder`

**认证**: 需要 JWT

**描述**: 批量更新问卷中所有题目的显示顺序

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**请求体**:

```json
{
    "question_ids": [3, 1, 2, 4]
}
```

**请求参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question_ids | array | 是 | 题目 ID 数组，按新的顺序排列 |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": null
}
```

**cURL 示例**:

```bash
curl -X PUT http://localhost:8080/api/v1/surveys/1/questions/reorder \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "question_ids": [3, 1, 2, 4]
  }'
```

---

## 4. 分享链接接口

### 4.1 生成分享链接

**端点**: `POST /api/v1/surveys/:id/share`

**认证**: 需要 JWT

**描述**: 为问卷生成加密的一次性分享链接，可包含预填数据

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**请求体**:

```json
{
    "prefill_data": {
        "name": "张三",
        "email": "zhangsan@example.com"
    },
    "expires_at": "2025-10-26T10:00:00Z"
}
```

**请求参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| prefill_data | object | 否 | 预填数据，键为题目的 prefill_key，值为预填值 |
| expires_at | string | 否 | 过期时间（ISO 8601 格式），默认 1 小时后 |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "url": "http://localhost:3000/survey/1?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_at": "2025-10-26T10:00:00Z"
    }
}
```

**cURL 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/surveys/1/share \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "prefill_data": {
        "name": "张三",
        "email": "zhangsan@example.com"
    },
    "expires_at": "2025-10-26T10:00:00Z"
  }'
```

---

## 5. 公开访问接口

### 5.1 获取问卷（通过 Token）

**端点**: `GET /api/v1/public/surveys/:id`

**认证**: 不需要

**描述**: 填答者通过加密 token 访问问卷，获取问卷详情和预填数据

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | string | 是 | 加密的一次性访问令牌 |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "survey": {
            "id": 1,
            "title": "客户满意度调查",
            "description": "本问卷旨在了解客户对我们服务的满意度",
            "status": "published",
            "questions": [
                {
                    "id": 1,
                    "type": "text",
                    "title": "您的姓名",
                    "description": "",
                    "required": true,
                    "order": 1,
                    "config": {},
                    "prefill_key": "name"
                },
                {
                    "id": 2,
                    "type": "text",
                    "title": "您的邮箱",
                    "description": "",
                    "required": true,
                    "order": 2,
                    "config": {},
                    "prefill_key": "email"
                }
            ]
        },
        "prefill_data": {
            "name": "张三",
            "email": "zhangsan@example.com"
        }
    }
}
```

**错误响应**:

- 400 Bad Request: Token 无效
- 403 Forbidden: Token 已过期或已使用

**cURL 示例**:

```bash
curl -X GET "http://localhost:8080/api/v1/public/surveys/1?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 5.2 提交问卷填答

**端点**: `POST /api/v1/public/responses`

**认证**: 不需要

**描述**: 填答者提交问卷答案

**请求体**:

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "answers": [
        {
            "question_id": 1,
            "value": "张三"
        },
        {
            "question_id": 2,
            "value": "满意"
        },
        {
            "question_id": 3,
            "value": ["选项1", "选项2"]
        },
        {
            "question_id": 4,
            "value": [
                ["张三", "30", "男"],
                ["李四", "25", "女"]
            ]
        }
    ]
}
```

**请求参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | string | 是 | 加密的一次性访问令牌 |
| answers | array | 是 | 答案数组 |
| answers[].question_id | integer | 是 | 题目 ID |
| answers[].value | any | 是 | 答案值（根据题目类型不同） |

**答案值类型说明**:

- **填空题 (text)**: 字符串，如 `"张三"`
- **单选题 (single)**: 字符串，如 `"满意"`
- **多选题 (multiple)**: 字符串数组，如 `["选项1", "选项2"]`
- **表格题 (table)**: 二维字符串数组，如 `[["张三", "30", "男"], ["李四", "25", "女"]]`

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "message": "提交成功"
    }
}
```

**错误响应**:

- 400 Bad Request: 数据验证失败（必填项缺失、选项不在范围内等）
- 403 Forbidden: Token 已过期或已使用
- 400 Bad Request: 问卷未发布

**cURL 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/public/responses \
  -H "Content-Type: application/json" \
  -d '{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "answers": [
        {
            "question_id": 1,
            "value": "张三"
        },
        {
            "question_id": 2,
            "value": "满意"
        }
    ]
  }'
```

---

## 6. 数据管理接口

### 6.1 查询填答记录

**端点**: `GET /api/v1/surveys/:id/responses`

**认证**: 需要 JWT

**描述**: 查询指定问卷的所有填答记录（支持分页）

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**查询参数**:

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | integer | 否 | 1 | 页码 |
| page_size | integer | 否 | 20 | 每页数量（最大 100） |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "survey_id": 1,
            "data": {
                "answers": [
                    {
                        "question_id": 1,
                        "value": "张三"
                    },
                    {
                        "question_id": 2,
                        "value": "满意"
                    }
                ]
            },
            "ip_address": "192.168.1.100",
            "user_agent": "Mozilla/5.0...",
            "submitted_at": "2025-10-25T12:00:00Z",
            "created_at": "2025-10-25T12:00:00Z"
        }
    ],
    "meta": {
        "page": 1,
        "page_size": 20,
        "total": 1
    }
}
```

**cURL 示例**:

```bash
curl -X GET "http://localhost:8080/api/v1/surveys/1/responses?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6.2 查询统计信息

**端点**: `GET /api/v1/surveys/:id/statistics`

**认证**: 需要 JWT

**描述**: 获取问卷的统计信息

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**成功响应** (200 OK):

```json
{
    "success": true,
    "data": {
        "survey_id": 1,
        "total_responses": 150,
        "completion_rate": 100.0
    }
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| survey_id | integer | 问卷 ID |
| total_responses | integer | 总填答数 |
| completion_rate | float | 完成率（百分比） |

**cURL 示例**:

```bash
curl -X GET http://localhost:8080/api/v1/surveys/1/statistics \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6.3 导出填答数据

**端点**: `GET /api/v1/surveys/:id/export`

**认证**: 需要 JWT

**描述**: 导出问卷的所有填答数据为 CSV 或 Excel 文件

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 问卷 ID |

**查询参数**:

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| format | string | 否 | csv | 导出格式：csv 或 excel |

**成功响应** (200 OK):

返回文件流，响应头包含：

```
Content-Type: text/csv; charset=utf-8
Content-Disposition: attachment; filename="survey_1_responses.csv"
```

或

```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="survey_1_responses.xlsx"
```

**cURL 示例**:

```bash
# 导出为 CSV
curl -X GET "http://localhost:8080/api/v1/surveys/1/export?format=csv" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o responses.csv

# 导出为 Excel
curl -X GET "http://localhost:8080/api/v1/surveys/1/export?format=excel" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o responses.xlsx
```

---

## 7. 完整使用流程示例

### 场景：创建问卷、生成分享链接、收集填答

#### 步骤 1: 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

响应：
```json
{
    "success": true,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3Mjk5NTg0MDB9.xxx",
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@example.com",
            "role": "admin"
        }
    }
}
```

保存 token 用于后续请求。

#### 步骤 2: 创建问卷

```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST http://localhost:8080/api/v1/surveys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "客户满意度调查",
    "description": "本问卷旨在了解客户对我们服务的满意度"
  }'
```

响应：
```json
{
    "success": true,
    "data": {
        "id": 1,
        "title": "客户满意度调查",
        "status": "draft"
    }
}
```

#### 步骤 3: 添加题目

```bash
# 添加填空题（姓名）
curl -X POST http://localhost:8080/api/v1/questions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "survey_id": 1,
    "type": "text",
    "title": "您的姓名",
    "required": true,
    "order": 1,
    "prefill_key": "name"
  }'

# 添加单选题（满意度）
curl -X POST http://localhost:8080/api/v1/questions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "survey_id": 1,
    "type": "single",
    "title": "您对我们的服务满意吗？",
    "required": true,
    "order": 2,
    "config": {
        "options": ["非常满意", "满意", "一般", "不满意"]
    }
  }'

# 添加多选题（改进建议）
curl -X POST http://localhost:8080/api/v1/questions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "survey_id": 1,
    "type": "multiple",
    "title": "您希望我们在哪些方面改进？",
    "required": false,
    "order": 3,
    "config": {
        "options": ["服务态度", "响应速度", "产品质量", "价格"]
    }
  }'
```

#### 步骤 4: 发布问卷

```bash
curl -X POST http://localhost:8080/api/v1/surveys/1/publish \
  -H "Authorization: Bearer $TOKEN"
```

#### 步骤 5: 生成分享链接

```bash
curl -X POST http://localhost:8080/api/v1/surveys/1/share \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "prefill_data": {
        "name": "张三"
    },
    "expires_at": "2025-10-26T10:00:00Z"
  }'
```

响应：
```json
{
    "success": true,
    "data": {
        "url": "http://localhost:3000/survey/1?token=abc123...",
        "token": "abc123...",
        "expires_at": "2025-10-26T10:00:00Z"
    }
}
```

#### 步骤 6: 填答者访问问卷

```bash
export SHARE_TOKEN="abc123..."

curl -X GET "http://localhost:8080/api/v1/public/surveys/1?token=$SHARE_TOKEN"
```

响应包含问卷详情和预填数据（姓名已预填为"张三"）。

#### 步骤 7: 填答者提交答案

```bash
curl -X POST http://localhost:8080/api/v1/public/responses \
  -H "Content-Type: application/json" \
  -d '{
    "token": "abc123...",
    "answers": [
        {
            "question_id": 1,
            "value": "张三"
        },
        {
            "question_id": 2,
            "value": "满意"
        },
        {
            "question_id": 3,
            "value": ["服务态度", "响应速度"]
        }
    ]
  }'
```

#### 步骤 8: 查看填答数据

```bash
# 查看所有填答记录
curl -X GET http://localhost:8080/api/v1/surveys/1/responses \
  -H "Authorization: Bearer $TOKEN"

# 查看统计信息
curl -X GET http://localhost:8080/api/v1/surveys/1/statistics \
  -H "Authorization: Bearer $TOKEN"

# 导出为 Excel
curl -X GET "http://localhost:8080/api/v1/surveys/1/export?format=excel" \
  -H "Authorization: Bearer $TOKEN" \
  -o responses.xlsx
```

---

## 8. 常见错误和解决方法

### 8.1 认证相关错误

#### 错误：401 Unauthorized - "未授权访问"

**原因**：
- 未提供 JWT token
- Token 格式错误
- Token 已过期

**解决方法**：
1. 确保在请求头中包含 `Authorization: Bearer <token>`
2. 检查 token 格式是否正确
3. 如果 token 过期，重新登录获取新 token

```bash
# 重新登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

#### 错误：403 Forbidden - "禁止访问"

**原因**：
- 尝试访问不属于自己的资源
- 没有足够的权限

**解决方法**：
- 确保只访问自己创建的问卷
- 检查用户角色权限

### 8.2 分享链接相关错误

#### 错误：400 Bad Request - "无效的令牌"

**原因**：
- Token 格式错误
- Token 解密失败
- Token 被篡改

**解决方法**：
- 使用完整的 token 字符串
- 不要手动修改 token
- 重新生成分享链接

#### 错误：403 Forbidden - "令牌已过期"

**原因**：
- 链接已超过设置的过期时间

**解决方法**：
- 联系管理员重新生成分享链接
- 生成链接时设置更长的过期时间

#### 错误：403 Forbidden - "链接已被使用"

**原因**：
- 一次性链接已经被使用过

**解决方法**：
- 联系管理员重新生成新的分享链接
- 每个链接只能使用一次，这是系统的安全设计

### 8.3 数据验证错误

#### 错误：400 Bad Request - "数据验证失败"

**原因**：
- 必填题目未回答
- 单选/多选题的选项不在配置的选项中
- 表格题的行数不符合限制
- 数据类型不匹配

**解决方法**：

1. **必填题目未回答**：
```json
// 错误示例：缺少必填题目
{
    "answers": [
        {"question_id": 2, "value": "满意"}
    ]
}

// 正确示例：包含所有必填题目
{
    "answers": [
        {"question_id": 1, "value": "张三"},
        {"question_id": 2, "value": "满意"}
    ]
}
```

2. **选项不在范围内**：
```json
// 错误示例：选项不存在
{
    "question_id": 2,
    "value": "非常非常满意"  // 该选项不在配置中
}

// 正确示例：使用配置中的选项
{
    "question_id": 2,
    "value": "非常满意"  // 该选项在配置中
}
```

3. **表格题行数不符合**：
```json
// 错误示例：行数少于 min_rows
{
    "question_id": 4,
    "value": []  // 如果 min_rows=1，这是错误的
}

// 正确示例：符合行数限制
{
    "question_id": 4,
    "value": [
        ["张三", "30", "男"]
    ]
}
```

### 8.4 问卷状态错误

#### 错误：400 Bad Request - "问卷未发布"

**原因**：
- 尝试访问或提交草稿状态的问卷

**解决方法**：
- 管理员需要先发布问卷：
```bash
curl -X POST http://localhost:8080/api/v1/surveys/1/publish \
  -H "Authorization: Bearer $TOKEN"
```

### 8.5 限流错误

#### 错误：429 Too Many Requests

**原因**：
- 请求频率超过限制（默认 100 次/分钟）

**解决方法**：
- 等待一分钟后重试
- 优化客户端请求频率
- 联系管理员调整限流配置

### 8.6 服务器错误

#### 错误：500 Internal Server Error

**原因**：
- 服务器内部错误
- 数据库连接失败
- Redis 连接失败

**解决方法**：
- 检查服务器日志
- 确认数据库和 Redis 服务正常运行
- 联系系统管理员

---

## 9. 最佳实践

### 9.1 安全建议

1. **保护 JWT Token**：
   - 不要在 URL 中传递 token
   - 使用 HTTPS 传输
   - 定期刷新 token

2. **分享链接管理**：
   - 设置合理的过期时间
   - 不要公开分享加密 token
   - 定期清理过期链接

3. **数据验证**：
   - 前端进行基础验证
   - 后端进行完整验证
   - 防止 XSS 和 SQL 注入

### 9.2 性能优化

1. **分页查询**：
   - 始终使用分页参数
   - 避免一次性加载大量数据

2. **缓存利用**：
   - 已发布的问卷会被缓存
   - 减少不必要的重复请求

3. **批量操作**：
   - 使用题目重排序接口批量更新
   - 避免循环调用单个更新接口

### 9.3 开发建议

1. **错误处理**：
   - 始终检查响应的 `success` 字段
   - 根据错误码提供友好的用户提示

2. **请求重试**：
   - 对于网络错误，实现指数退避重试
   - 对于 429 错误，等待后重试

3. **日志记录**：
   - 记录所有 API 请求和响应
   - 便于问题排查和审计

---

## 10. 附录

### 10.1 HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录 |
| 403 | 禁止访问，权限不足 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

### 10.2 题目类型说明

| 类型 | 值 | 答案格式 | 说明 |
|------|-----|----------|------|
| 填空题 | text | string | 单行文本输入 |
| 单选题 | single | string | 从选项中选择一个 |
| 多选题 | multiple | string[] | 从选项中选择多个 |
| 表格题 | table | string[][] | 多行多列数据 |

### 10.3 配置参数

**数据库连接池**：
- 最大连接数：100
- 最大空闲连接：10
- 连接最大生命周期：1 小时

**Redis 配置**：
- 连接池大小：10
- 问卷缓存 TTL：1 小时
- 链接状态缓存 TTL：与链接过期时间相同

**限流配置**：
- 每分钟请求数：100
- 突发请求数：20

**JWT 配置**：
- 过期时间：24 小时
- 算法：HS256

**一次性链接**：
- 默认过期时间：1 小时
- 最大过期时间：7 天

---

## 联系方式

如有问题或建议，请联系开发团队。

**文档版本**: v1.0  
**最后更新**: 2025-10-25
