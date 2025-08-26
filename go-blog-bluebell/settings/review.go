package settings

// ReviewConfig 审核配置
type ReviewConfig struct {
	// n8n配置
	N8nURL     string `mapstructure:"n8n_url"`     // n8n服务地址
	N8nTimeout int    `mapstructure:"n8n_timeout"` // n8n请求超时时间（秒）

	// 审核配置
	AutoReview    bool    `mapstructure:"auto_review"`     // 是否启用自动审核
	ReviewTimeout int     `mapstructure:"review_timeout"`  // 审核超时时间（秒）
	MaxRetryCount int     `mapstructure:"max_retry_count"` // 最大重试次数
	MinScore      float64 `mapstructure:"min_score"`       // 最低通过分数

	// 通知配置
	EnableNotification bool   `mapstructure:"enable_notification"` // 是否启用通知
	EmailTemplate      string `mapstructure:"email_template"`      // 邮件模板
	SMSTemplate        string `mapstructure:"sms_template"`        // 短信模板

	// 缓存配置
	CacheExpireTime int `mapstructure:"cache_expire_time"` // 缓存过期时间（小时）
}

// DefaultReviewConfig 默认审核配置
func DefaultReviewConfig() *ReviewConfig {
	return &ReviewConfig{
		N8nURL:             "http://8e0740914f7a.ngrok-free.app/webhook-test/post-review",
		N8nTimeout:         30,
		AutoReview:         true,
		ReviewTimeout:      300, // 5分钟
		MaxRetryCount:      3,
		MinScore:           60.0,
		EnableNotification: true,
		EmailTemplate:      "您的帖子审核结果：{status}",
		SMSTemplate:        "帖子审核{status}，详情请查看邮件",
		CacheExpireTime:    24, // 24小时
	}
}
