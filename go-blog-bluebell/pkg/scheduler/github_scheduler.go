package scheduler

import (
	"time"

	"github.com/LucienLSA/go-blog/service"
	"go.uber.org/zap"
)

// GitHubTrendingScheduler GitHub热点数据定时任务调度器
type GitHubTrendingScheduler struct {
	githubService *service.GitHubService
	ticker        *time.Ticker
	stopChan      chan bool
}

// NewGitHubTrendingScheduler 创建GitHub热点数据定时任务调度器
func NewGitHubTrendingScheduler() *GitHubTrendingScheduler {
	return &GitHubTrendingScheduler{
		githubService: service.GetGitHubService(),
		stopChan:      make(chan bool),
	}
}

// Start 启动定时任务
func (s *GitHubTrendingScheduler) Start() {
	zap.L().Info("starting github trending data scheduler")
	
	// 立即执行一次数据刷新
	s.githubService.RefreshGitHubTrendingData()
	
	// 创建定时器，每小时执行一次
	s.ticker = time.NewTicker(1 * time.Hour)
	
	go func() {
		for {
			select {
			case <-s.ticker.C:
				zap.L().Info("executing scheduled github trending data refresh")
				s.githubService.RefreshGitHubTrendingData()
			case <-s.stopChan:
				zap.L().Info("stopping github trending data scheduler")
				return
			}
		}
	}()
}

// Stop 停止定时任务
func (s *GitHubTrendingScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
}

// 全局调度器实例
var githubScheduler *GitHubTrendingScheduler

// StartGitHubTrendingScheduler 启动GitHub热点数据定时任务
func StartGitHubTrendingScheduler() {
	if githubScheduler == nil {
		githubScheduler = NewGitHubTrendingScheduler()
	}
	githubScheduler.Start()
}

// StopGitHubTrendingScheduler 停止GitHub热点数据定时任务
func StopGitHubTrendingScheduler() {
	if githubScheduler != nil {
		githubScheduler.Stop()
	}
} 