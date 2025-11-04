# 服务端部署指南（面板）

本文介绍两种部署方式：

- 二进制一键脚本部署（systemd，Linux）
- Docker Compose 部署

在开始前，请准备：
- 一台 Linux 服务器（建议 Ubuntu 20.04+/Debian 11+/CentOS 8+）
- 已开放面板端口（默认 6365 提供 API；前端静态资源可由反代提供 HTTPS）
- MySQL 数据库（或在 Docker Compose 中随容器启动）

---
## 方式一：二进制一键脚本部署（Linux）

脚本位置：`scripts/install_server.sh`

步骤：
1）下载并执行安装脚本（root 权限）：

```bash
curl -fsSL https://raw.githubusercontent.com/NiuStar/flux-panel/refs/heads/main/scripts/install_server.sh -o install_server.sh \
  && sudo bash install_server.sh
```

2）按提示选择：
- 是否使用下载代理前缀（可为空）
- CPU 架构（默认自动识别）
- 选择从 GitHub Releases 下载预编译，或本地源码编译（需要已安装 Go）

3）服务与配置：
- systemd 服务名：`flux-panel`
- 可执行文件：`/usr/local/bin/flux-panel-server`
- 工作目录：`/opt/flux-panel`
- 环境配置：`/etc/default/flux-panel`

环境变量说明：
```
PORT=6365               # 面板后端监听端口
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=flux_panel
DB_USER=flux
DB_PASSWORD=123456
```

4）常用命令：
```bash
sudo systemctl status flux-panel
sudo systemctl restart flux-panel
sudo journalctl -u flux-panel -f
```

> 首次启动会自动创建数据库（如权限允许）与管理员账号（admin_user/admin_user），请尽快登录修改密码。

---
## 方式二：Docker Compose 部署

仓库内提供 `docker-compose-v4.yml`，可与 `panel_install.sh` 搭配使用或手动部署。

1）准备环境与变量
- 确保 Docker 与 Docker Compose 可用
- 准备 `.env` 文件（样例变量参考 `panel_install.sh` 交互生成），至少包括：
  - 面板访问域名/端口
  - 数据库相关变量（DB_HOST/DB_NAME/DB_USER/DB_PASSWORD）

2）启动服务
```bash
docker compose -f docker-compose-v4.yml up -d
```

3）反向代理（可选）
- 使用仓库内 `proxy.sh` 可快速配置 Caddy/Nginx 反代，或自行配置 HTTPS 证书与反代至后端端口（默认 6365）

4）升级/重启
```bash
docker compose -f docker-compose-v4.yml pull
docker compose -f docker-compose-v4.yml up -d
```

---
## 端口与安全
- 后端默认监听 6365（可通过 `PORT` 修改）
- 建议将前端静态资源置于反代服务器并启用 HTTPS
- 不要在公开渠道泄露 `.env`、数据库密码、JWT 等敏感信息

