# Changelog

## v0.0.3
_发布日期：2026-10-24_

- build: require published go-events v0.0.5

## v0.0.2
_发布日期：2026-10-12_

- build: 依赖改为已发布版本，移除本地 replace
- 更新 go-sqlserver-storage
- feat: 实现 Capabilities() 方法
- 升级go-storage-test-helper版本
- 把SqlServer主键重复的错误转换为StorageLock框架内部的Version Miss的错误
- SqlServerConnectionManager支持从sql.DB创建
- 修改为基于sql based storage实现
- 升级storage版本
- 增加github action ci
- 完善文档
- 完善API
- 增加本地测试docker脚本
- 增加examples

## v0.0.1
_发布日期：2023-08-06_

- rename func

## Unreleased

- build: 升级依赖到 go-storage/go-storage-lock v0.1.0
