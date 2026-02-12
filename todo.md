# TODO

- [x] 支持按空格分割多个关键词，使用 AND 关系搜索
- [x] 优化增强：加载数据后收集所有的 tags信息，固定到 services-container 左外侧
  - 显示每个tag 的 item 数量
  - 点击指定 tag 后页面只显示对应的 items
- [x] 增强：如果在浏览器输入了 query 参数 `refresh=true` 需要将此参数带给后端 API `/api/page/*` 刷新缓存数据
