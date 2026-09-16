# 资质证书管理系统

纯 Go 标准库实现的后端服务，零第三方依赖。

## 运行

```bash
cd origin
go run ./cmd/server
```

访问 http://localhost:8080 查看前端页面。

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/certificates | 创建证书 |
| GET | /api/certificates | 证书列表 |
| GET | /api/certificates/{id} | 证书详情 |
| PUT | /api/certificates/{id} | 更新证书 |
| DELETE | /api/certificates/{id} | 删除证书 |
| POST | /api/certificates/{id}/revoke | 吊销证书 |
| POST | /api/certificates/batch-import | 批量导入证书 |
| POST | /api/certificates/batch-revoke | 按持有人批量吊销 |
| POST | /api/categories | 创建分类 |
| GET | /api/categories | 分类列表 |
| GET | /api/categories/{id} | 分类详情 |
| PUT | /api/categories/{id} | 更新分类 |
| DELETE | /api/categories/{id} | 删除分类 |
| POST | /api/holders | 创建持有人 |
| GET | /api/holders | 持有人列表 |
| GET | /api/holders/{id} | 持有人详情 |
| PUT | /api/holders/{id} | 更新持有人 |
| DELETE | /api/holders/{id} | 删除持有人 |
| POST | /api/issuers | 创建机构 |
| GET | /api/issuers | 机构列表 |
| GET | /api/issuers/{id} | 机构详情 |
| PUT | /api/issuers/{id} | 更新机构 |
| DELETE | /api/issuers/{id} | 删除机构 |
| POST | /api/verifications | 创建验真记录 |
| GET | /api/verifications | 验真列表 |
| GET | /api/verifications/{id} | 验真详情 |
| PUT | /api/verifications/{id} | 更新验真 |
| DELETE | /api/verifications/{id} | 删除验真 |
| POST | /api/reminders | 创建提醒 |
| GET | /api/reminders | 提醒列表 |
| GET | /api/reminders/{id} | 提醒详情 |
| PUT | /api/reminders/{id} | 更新提醒 |
| DELETE | /api/reminders/{id} | 删除提醒 |
| POST | /api/renewals | 创建续期申请 |
| GET | /api/renewals | 续期列表 |
| GET | /api/renewals/{id} | 续期详情 |
| PUT | /api/renewals/{id} | 更新续期 |
| DELETE | /api/renewals/{id} | 删除续期 |
| POST | /api/renewals/{id}/approve | 审批通过 |
| POST | /api/renewals/{id}/reject | 驳回 |
| POST | /api/renewals/{id}/complete | 完成续期 |
| POST | /api/revocations | 创建吊销记录 |
| GET | /api/revocations | 吊销列表 |
| GET | /api/revocations/{id} | 吊销详情 |
| DELETE | /api/revocations/{id} | 删除吊销记录 |
| POST | /api/attachments | 创建附件 |
| GET | /api/attachments | 附件列表 |
| GET | /api/attachments/{id} | 附件详情 |
| DELETE | /api/attachments/{id} | 删除附件 |
| POST | /api/notifications | 创建通知 |
| GET | /api/notifications | 通知列表 |
| GET | /api/notifications/{id} | 通知详情 |
| POST | /api/notifications/{id}/send | 发送通知 |
| DELETE | /api/notifications/{id} | 删除通知 |
| POST | /api/audit-logs | 创建审计日志 |
| GET | /api/audit-logs | 审计日志列表 |
| GET | /api/audit-logs/{id} | 审计日志详情 |
| DELETE | /api/audit-logs/{id} | 删除审计日志 |
| GET | /api/stats/overview | 总览统计 |
| GET | /api/stats/by-category | 按分类统计 |
| GET | /api/stats/by-status | 按状态统计 |
| GET | /api/stats/monthly-trend | 按月签发趋势 |
| GET | /api/stats/top-issuers | 机构签发排行 |
| GET | /api/stats/expiring-soon | 即将到期证书 |
| GET | /api/stats/export | 全量导出快照 |
