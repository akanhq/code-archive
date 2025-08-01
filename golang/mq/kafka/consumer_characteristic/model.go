package consumer_characteristic

type Task struct {
	ID         uint64 `json:"id"`
	Priority   int    `json:"priority"`    // 优先级，越大越优先
	Data       string `json:"data"`        // 任务内容
	ExecuteAt  int64  `json:"execute_at"`  // 执行时间戳
	Status     string `json:"status"`      // 状态: pending, processing, success, failed
	RetryCount int    `json:"retry_count"` // 重试次数
	StartedAt  int64  `json:"started_at"`  // 开始处理时间戳（用于超时检测）
}
