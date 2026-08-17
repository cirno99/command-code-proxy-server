# CommandCode Proxy Server 配置指南

## 概述

CommandCode Proxy Server 是一个 OpenAI 兼容的代理服务器，可以将 CommandCode API 包装为 OpenAI 兼容的接口，让 pi 可以使用 CommandCode 的模型。

## 快速开始

### 1. 启动代理服务器

```bash
# 使用默认配置启动
./bin/command-code-proxy

# 使用自定义端口启动
./bin/command-code-proxy -port 8080

# 使用默认 API 密钥启动
./bin/command-code-proxy -api-key your-commandcode-api-key
```

代理服务器默认运行在 `http://127.0.0.1:55990`

### 2. 配置 pi 的 models.json

将 `models.json.example` 的内容复制到 `~/.pi/agent/models.json`：

```bash
cp models.json.example ~/.pi/agent/models.json
```

或者手动编辑 `~/.pi/agent/models.json`，添加 `commandcode` 提供者配置。

## 配置详解

### 基本配置

```json
{
  "providers": {
    "commandcode": {
      "baseUrl": "http://127.0.0.1:55990",
      "api": "openai-completions",
      "models": [
        {
          "id": "deepseek-v4-flash",
          "name": "DeepSeek V4 Flash",
          "reasoning": true,
          "input": ["text"],
          "contextWindow": 200000,
          "maxTokens": 200000
        }
      ]
    }
  }
}
```

### 配置字段说明

#### 提供者级别

| 字段 | 说明 |
|------|------|
| `baseUrl` | 代理服务器地址，默认 `http://127.0.0.1:55990` |
| `api` | API 类型，使用 `openai-completions` |
| `apiKey` | 可选，CommandCode API 密钥（如果不在启动时指定） |

#### 模型级别

| 字段 | 说明 |
|------|------|
| `id` | 模型 ID（使用别名或完整 ID） |
| `name` | 显示名称 |
| `reasoning` | 是否支持推理（思考链） |
| `input` | 输入类型：`["text"]` 或 `["text", "image"]` |
| `contextWindow` | 上下文窗口大小（token 数） |
| `maxTokens` | 最大输出 token 数 |

### 支持的模型

#### DeepSeek 系列
- `deepseek-v4-pro` / `deepseek-v4` / `deepseek-pro` → `deepseek/deepseek-v4-pro`
- `deepseek-v4-flash` / `deepseek-flash` → `deepseek/deepseek-v4-flash`

#### MiniMax 系列
- `minimax-m2.7` / `minimax2.7` → `MiniMaxAI/MiniMax-M2.7`
- `minimax-m2.5` / `minimax2.5` / `minimax` → `MiniMaxAI/MiniMax-M2.5`
- `minimax-m3` / `minimax3` → `MiniMaxAI/MiniMax-M3`

#### GLM 系列
- `glm-5.1` → `zai-org/GLM-5.1`
- `glm-5` → `zai-org/GLM-5`

#### Kimi 系列
- `kimi-k2.6` / `kimi2.6` → `moonshotai/Kimi-K2.6`
- `kimi-k2.5` / `kimi2.5` → `moonshotai/Kimi-K2.5`

#### Qwen 系列
- `qwen-3.6-max-preview` / `qwen3.6-max` → `Qwen/Qwen3.6-Max-Preview`
- `qwen-3.6-plus` / `qwen3.6-plus` / `qwen3.6` → `Qwen/Qwen3.6-Plus`
- `qwen-3.7-max-free` / `qwen3.7-max-free` → `Qwen/Qwen3.7-Max-Free`
- `qwen-3.7-max` / `qwen3.7-max` → `Qwen/Qwen3.7-Max`

#### Step 系列
- `step-3.5-flash` / `step3.5` → `stepfun/Step-3.5-Flash`
- `step-3.7-flash` / `step3.7` → `stepfun/Step-3.7-Flash`

#### Gemini 系列
- `gemini-3.1-flash-lite` / `gemini-flash-lite` → `google/gemini-3.1-flash-lite`

#### Mimo 系列
- `mimo-v2.5-pro` / `mimo-pro` → `xiaomi/mimo-v2.5-pro`
- `mimo-v2.5` / `mimo` → `xiaomi/mimo-v2.5`

### API 密钥管理

代理服务器按以下顺序查找 API 密钥：

1. 客户端请求的 `Authorization` 头
2. 启动时的 `-api-key` 参数
3. 如果都没有，返回 `401 Unauthorized`

#### 启动时指定密钥

```bash
./bin/command-code-proxy -api-key your-commandcode-api-key
```

#### 在 models.json 中指定密钥

```json
{
  "providers": {
    "commandcode": {
      "baseUrl": "http://127.0.0.1:55990",
      "api": "openai-completions",
      "apiKey": "your-commandcode-api-key",
      "models": [...]
    }
  }
}
```

#### 使用环境变量

```json
{
  "providers": {
    "commandcode": {
      "baseUrl": "http://127.0.0.1:55990",
      "api": "openai-completions",
      "apiKey": "$COMMANDCODE_API_KEY",
      "models": [...]
    }
  }
}
```

## 测试配置

### 1. 检查代理服务器状态

```bash
curl http://127.0.0.1:55990/health
```

### 2. 列出可用模型

```bash
curl http://127.0.0.1:55990/v1/models
```

### 3. 测试聊天完成

```bash
curl http://127.0.0.1:55990/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-commandcode-api-key" \
  -d '{
    "model": "deepseek-v4-flash",
    "messages": [
      {"role": "user", "content": "Hello"}
    ],
    "stream": false
  }'
```

## 使用示例

### 在 pi 中使用

1. 启动代理服务器：
   ```bash
   ./bin/command-code-proxy -api-key your-key
   ```

2. 在 pi 中选择模型：
   ```bash
   /model
   ```
   然后选择 `commandcode` 提供者下的模型

3. 或者使用命令行参数：
   ```bash
   pi --model commandcode/deepseek-v4-flash
   ```

### 完整配置示例

```json
{
  "providers": {
    "sudocode": {
      "baseUrl": "https://api.sudorelay.com",
      "api": "openai-completions",
      "models": [
        {
          "id": "gpt-5.6-luna",
          "name": "GPT-5.6 Luna",
          "reasoning": true,
          "input": ["text", "image"],
          "contextWindow": 1050000,
          "maxTokens": 128000
        }
      ]
    },
    "commandcode": {
      "baseUrl": "http://127.0.0.1:55990",
      "api": "openai-completions",
      "models": [
        {
          "id": "deepseek-v4-pro",
          "name": "DeepSeek V4 Pro",
          "reasoning": true,
          "input": ["text"],
          "contextWindow": 200000,
          "maxTokens": 200000
        },
        {
          "id": "deepseek-v4-flash",
          "name": "DeepSeek V4 Flash",
          "reasoning": true,
          "input": ["text"],
          "contextWindow": 200000,
          "maxTokens": 200000
        },
        {
          "id": "minimax-m3",
          "name": "MiniMax M3",
          "reasoning": true,
          "input": ["text"],
          "contextWindow": 128000,
          "maxTokens": 128000
        }
      ]
    }
  }
}
```

## 故障排除

### 代理服务器无法启动

1. 检查端口是否被占用：
   ```bash
   lsof -i :55990
   ```

2. 检查防火墙设置

### 模型不可用

1. 确保代理服务器正在运行
2. 检查 API 密钥是否正确
3. 查看代理服务器日志

### 连接超时

1. 检查网络连接
2. 确认代理服务器地址和端口正确
3. 检查防火墙设置

## 高级配置

### 自定义模型参数

如果需要自定义模型参数，可以在模型配置中添加 `samplingParams`：

```json
{
  "id": "deepseek-v4-flash",
  "name": "DeepSeek V4 Flash",
  "reasoning": true,
  "input": ["text"],
  "contextWindow": 200000,
  "maxTokens": 200000,
  "samplingParams": {
    "temperature": 0.7,
    "top_p": 0.9
  }
}
```

### 多模型配置

可以配置多个模型，pi 会在 `/model` 命令中显示所有可用模型：

```json
{
  "providers": {
    "commandcode": {
      "baseUrl": "http://127.0.0.1:55990",
      "api": "openai-completions",
      "models": [
        {
          "id": "deepseek-v4-pro",
          "name": "DeepSeek V4 Pro",
          "reasoning": true,
          "input": ["text"],
          "contextWindow": 200000,
          "maxTokens": 200000
        },
        {
          "id": "deepseek-v4-flash",
          "name": "DeepSeek V4 Flash",
          "reasoning": true,
          "input": ["text"],
          "contextWindow": 200000,
          "maxTokens": 200000
        },
        {
          "id": "minimax-m3",
          "name": "MiniMax M3",
          "reasoning": true,
          "input": ["text"],
          "contextWindow": 128000,
          "maxTokens": 128000
        }
      ]
    }
  }
}
```

## 相关链接

- 项目仓库: https://github.com/dev2k6/command-code-proxy-server
- CommandCode API: https://api.commandcode.ai