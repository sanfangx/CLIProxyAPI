# CLIProxyAPI 巨魔 (TrollStore) iOS 版使用说明

本项目是专门针对 iOS 巨魔 (TrollStore) 环境适配优化的 [CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) 版本。通过单文件纯静态编译、内置 Mozilla 根证书与离线 Web 管理面板，将一台闲置的 iPhone 或 iPad 变身为随身、24/7 低功耗运行的个人 AI 模型中转网关。

---

## 一、 特性亮点

1. **单文件静态二进制**：纯 Go 语言静态编译，**不需要在 iOS 上安装 Python、Node.js 或任何额外环境**。
2. **极低资源开销**：日常静默运行物理内存仅 **25MB ~ 45MB**，毫秒级启动。
3. **内置 Mozilla 根证书支持**：彻底解决纯 Go 在 iOS 系统下缺失 `/etc/ssl/cert.pem` 导致的 `x509: certificate signed by unknown authority` 报错。
4. **离线 WebUI 控制中心**：预置 `management.html`，手机 Safari 直接打开 `http://127.0.0.1:8317/management.html` 即可视化管理与完成 OAuth 登录。
5. **协议全兼容**：提供标准的 OpenAI (`/v1/chat/completions`)、Claude (`/v1/messages`) 与 Gemini (`/v1beta`) 格式接口，支持多账号自动轮询负载均衡。

---

## 二、 安装与运行指南

### 方法一：通过 TrollService 守护进程管理器运行（推荐）

1. 在手机上安装并打开 **TrollService** 应用。
2. 点击右上方「添加服务」，选择下载好的 `cli-proxy-api` 二进制文件：
   * **服务 ID**：`cliproxyapi`
   * **名称**：`CLI Proxy API`
   * **监听端口**：`8317`
   * **运行参数**：`--config config/config.yaml`
3. 部署完成后点击「启动服务」，TrollService 会自动为程序注入 CoreTrust 绕过签名并以 root/persona 99 身份在后台常驻运行。
4. 点击「打开管理面板」或在 Safari 中访问 `http://127.0.0.1:8317/management.html`。

### 方法二：通过 Filza / 终端命令行运行

1. 将解压出来的整个文件夹放置于 `/var/mobile/Documents/service/cliproxyapi`。
2. 为可执行程序赋予权限与 TrollStore 签名：
   ```bash
   chmod 755 bin/cli-proxy-api
   ldid -SRootHelper.entitlements bin/cli-proxy-api
   fastPathSign bin/cli-proxy-api
   ```
3. 启动后台服务：
   ```bash
   cd /var/mobile/Documents/service/cliproxyapi
   nohup ./bin/cli-proxy-api --config config/config.yaml > runtime/service.log 2>&1 &
   ```

---

## 三、 客户端接入方法

### 1. iPhone 本地客户端（如 OpenCat / Chatbox iOS）
* **API 地址 / Base URL**：`http://127.0.0.1:8317/v1`
* **API Key**：`sk-trollservice-local-key`

### 2. 局域网电脑客户端（VSCode / Cursor / Claude Code）
* **API 地址 / Base URL**：`http://<手机局域网IP>:8317/v1` (例如 `http://192.168.1.50:8317/v1`)
* **API Key**：`sk-trollservice-local-key`

---

## 四、 目录结构说明

```text
cliproxyapi/
├── bin/
│   └── cli-proxy-api        # iOS arm64 可执行程序
├── config/
│   ├── config.yaml          # 服务主配置文件
│   ├── cacert.pem           # Mozilla CA 证书
│   └── .env                 # 环境变量配置
├── static/
│   └── management.html      # 离线 Web 管理控制台
├── auth/                    # OAuth 账号凭据存储目录 (*.json)
└── runtime/
    ├── service.pid          # 进程 PID
    └── service.log          # 运行日志
```
