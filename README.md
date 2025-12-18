# LINUX DO Credit

🚀 Linux Do 社区 Credit (Content Distribution Kit) 快速分享平台

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-16-black.svg)](https://nextjs.org/)
[![React](https://img.shields.io/badge/React-19-blue.svg)](https://reactjs.org/)

[![GitHub release](https://img.shields.io/github/v/release/linux-do/pay?include_prereleases)](https://github.com/linux-do/pay/releases)
[![GitHub stars](https://img.shields.io/github/stars/linux-do/pay)](https://github.com/linux-do/pay/stargazers) 
[![GitHub forks](https://img.shields.io/github/forks/linux-do/pay)](https://github.com/linux-do/pay/network)
[![GitHub issues](https://img.shields.io/github/issues/linux-do/pay)](https://github.com/linux-do/pay/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/linux-do/pay)](https://github.com/linux-do/pay/pulls)
[![GitHub contributors](https://img.shields.io/github/contributors/linux-do/pay)](https://github.com/linux-do/pay/graphs/contributors)

[![Backend Build](https://github.com/linux-do/pay/actions/workflows/build_backend.yml/badge.svg)](https://github.com/linux-do/pay/actions/workflows/build_backend.yml)
[![Frontend Build](https://github.com/linux-do/pay/actions/workflows/build_frontend.yml/badge.svg)](https://github.com/linux-do/pay/actions/workflows/build_frontend.yml)
[![Docker Build](https://github.com/linux-do/pay/actions/workflows/build_image.yml/badge.svg)](https://github.com/linux-do/pay/actions/workflows/build_image.yml)
[![CodeQL](https://github.com/linux-do/pay/actions/workflows/codeql.yml/badge.svg)](https://github.com/linux-do/pay/actions/workflows/codeql.yml)
[![ESLint](https://github.com/linux-do/pay/actions/workflows/eslint.yml/badge.svg)](https://github.com/linux-do/pay/actions/workflows/eslint.yml)

## 📖 项目简介

LINUX DO Credit 是一个为 Linux Do 社区打造的积分服务平台，旨在提供一系列积分相关服务，为社区开发者提供积分流转基础框架。
### ✨ 主要特性

- 🔐 **OAuth2 认证** - 集成 Linux Do 社区账号系统
- 🛡️ **风险控制** - 完善的信任等级和风险评估系统
- 📊 **实时监控** - 详细的分发统计和用户行为分析
- 🎨 **现代化界面** - 基于 Next.js 16 和 React 19 的响应式设计
- ⚡ **高性能** - Go 后端 + Redis 缓存 + PostgreSQL 数据库

## 🏗️ 架构概览

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Database      │
│   (Next.js)     │◄──►│     (Go)        │◄──►│  (PostgreSQL)   │
│                 │    │                 │    │                 │
│ • React 19      │    │ • Gin Framework │    │ • PostgreSQL    │
│ • TypeScript    │    │ • OAuth2        │    │ • Redis Cache   │
│ • Tailwind CSS  │    │ • OpenTelemetry │    │ • Session Store │
│ • Shadcn UI     │    │ • Swagger API   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ 技术栈

### 后端
- **Go 1.25.4** - 主要开发语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **Redis** - 缓存和会话存储
- **PostgreSQL** - 主数据库
- **OpenTelemetry** - 可观测性
- **Swagger** - API 文档

### 前端
- **Next.js 16** - React 框架
- **React 19** - UI 库
- **TypeScript** - 类型安全
- **Tailwind CSS 4** - 样式框架
- **Shadcn UI** - 组件库
- **Lucide Icons** - 图标库

## 📋 环境要求

- **Go** >= 1.25.4
- **Node.js** >= 18.0
- **PostgreSQL** >= 18
- **Redis** >= 6.0
- **pnpm** >= 8.0 (推荐)