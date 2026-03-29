package services

import (
	"Today-Todo/models"
	"context"
	"log"
	"sync"
	"time"
)

// Scheduler 任务调度器
type Scheduler struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	waterInterval time.Duration // 喝水提醒间隔
	standInterval time.Duration // 站立提醒间隔
	videoInterval time.Duration // 短视频提醒间隔
	stateMachines map[uint]*models.StateMachine
	mu            sync.RWMutex
	hub           *ReminderHub
}

// NewScheduler 创建新的调度器
func NewScheduler(hub *ReminderHub) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		ctx:           ctx,
		cancel:        cancel,
		waterInterval: 90 * time.Minute, // 每90分钟提醒喝水
		standInterval: 60 * time.Minute, // 每60分钟提醒站立
		videoInterval: 120 * time.Minute,
		stateMachines: make(map[uint]*models.StateMachine),
		hub:           hub,
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	log.Println("任务调度器启动")

	// 并行启动多个goroutine处理不同类型的任务
	s.wg.Add(4)

	// 1. 喝水提醒goroutine
	go s.waterReminderWorker()

	// 2. 站立提醒goroutine
	go s.standReminderWorker()

	// 3. 状态监控goroutine
	go s.stateMonitorWorker()

	// 4. 短视频提醒 goroutine
	go s.shortVideoReminderWorker()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	log.Println("任务调度器停止")
	s.cancel()
	s.wg.Wait()
}

// waterReminderWorker 喝水提醒工作goroutine
func (s *Scheduler) waterReminderWorker() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.waterInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.sendWaterReminder()
		}
	}
}

// standReminderWorker 站立提醒工作goroutine
func (s *Scheduler) standReminderWorker() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.standInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.sendStandReminder()
		}
	}
}

// stateMonitorWorker 状态监控工作goroutine
func (s *Scheduler) stateMonitorWorker() {
	defer s.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkAndSwitchStates()
		}
	}
}

// shortVideoReminderWorker 防沉迷提醒 goroutine。
func (s *Scheduler) shortVideoReminderWorker() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.videoInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.sendShortVideoReminder()
		}
	}
}

// sendWaterReminder 发送喝水提醒
func (s *Scheduler) sendWaterReminder() {
	log.Println("💧 喝水提醒：该喝水了！建议饮水200-300ml")
	if s.hub != nil {
		s.hub.Broadcast(ReminderEvent{
			Type:      "water",
			Title:     "补水时间到",
			Message:   "建议补充 200-300ml 水，保持专注状态。",
			Timestamp: time.Now(),
		})
	}
}

// sendStandReminder 发送站立提醒
func (s *Scheduler) sendStandReminder() {
	log.Println("🚶 站立提醒：久坐提醒，建议站立活动5-10分钟")
	if s.hub != nil {
		s.hub.Broadcast(ReminderEvent{
			Type:      "stand",
			Title:     "起身活动一下",
			Message:   "已经久坐 1 小时，建议站立拉伸 5 分钟。",
			Timestamp: time.Now(),
		})
	}
}

// sendShortVideoReminder 发送防沉迷提醒。
func (s *Scheduler) sendShortVideoReminder() {
	log.Println("📵 防沉迷提醒：减少短视频连续刷屏，回到你的目标任务")
	if s.hub != nil {
		s.hub.Broadcast(ReminderEvent{
			Type:      "short-video",
			Title:     "防沉迷提醒",
			Message:   "短视频可设 10 分钟上限，继续专注会更高效。",
			Timestamp: time.Now(),
		})
	}
}

// checkAndSwitchStates 检查并切换状态
func (s *Scheduler) checkAndSwitchStates() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for userID, sm := range s.stateMachines {
		state := sm.GetCurrentState()
		duration := time.Since(sm.LastChange)

		// 工作超过50分钟，建议休息
		if state == models.StateWorking && duration > 50*time.Minute {
			log.Printf("用户 %d 已工作 %.0f 分钟，建议切换到休息状态",
				userID, duration.Minutes())
		}

		// 休息超过15分钟，建议回到工作
		if state == models.StateResting && duration > 15*time.Minute {
			log.Printf("用户 %d 已休息 %.0f 分钟，可以恢复工作了",
				userID, duration.Minutes())
		}
	}
}

// RegisterUser 注册用户到调度器
func (s *Scheduler) RegisterUser(userID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stateMachines[userID] = models.NewStateMachine(userID)
}

// GetStateMachine 获取用户状态机
func (s *Scheduler) GetStateMachine(userID uint) *models.StateMachine {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stateMachines[userID]
}

// ReminderConfig 返回提醒间隔配置（分钟），便于前端兜底定时器使用。
func (s *Scheduler) ReminderConfig() map[string]int {
	return map[string]int{
		"water_minutes":       int(s.waterInterval.Minutes()),
		"stand_minutes":       int(s.standInterval.Minutes()),
		"short_video_minutes": int(s.videoInterval.Minutes()),
	}
}
