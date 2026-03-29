package models

import (
	"sync"
	"time"
)

// UserState 用户状态枚举
type UserState string

const (
	StateWorking  UserState = "working"  // 工作中
	StateResting  UserState = "resting"  // 休息中
	StateIdle     UserState = "idle"     // 空闲
	StateExercise UserState = "exercise" // 运动中
)

// StateMachine 状态机
type StateMachine struct {
	mu           sync.RWMutex
	UserID       uint
	CurrentState UserState
	LastChange   time.Time
	WorkDuration time.Duration // 工作时长
	RestDuration time.Duration // 休息时长
}

// NewStateMachine 创建新的状态机
func NewStateMachine(userID uint) *StateMachine {
	return &StateMachine{
		UserID:       userID,
		CurrentState: StateIdle,
		LastChange:   time.Now(),
	}
}

// TransitionTo 状态转换
func (sm *StateMachine) TransitionTo(newState UserState) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 验证状态转换的合法性
	if !sm.isValidTransition(newState) {
		return false
	}

	// 记录上一个状态的持续时间
	duration := time.Since(sm.LastChange)
	if sm.CurrentState == StateWorking {
		sm.WorkDuration += duration
	} else if sm.CurrentState == StateResting {
		sm.RestDuration += duration
	}

	sm.CurrentState = newState
	sm.LastChange = time.Now()
	return true
}

// GetCurrentState 获取当前状态
func (sm *StateMachine) GetCurrentState() UserState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.CurrentState
}

// isValidTransition 验证状态转换的合法性
func (sm *StateMachine) isValidTransition(newState UserState) bool {
	// 定义状态转换规则
	validTransitions := map[UserState][]UserState{
		StateIdle:     {StateWorking, StateExercise},
		StateWorking:  {StateResting, StateIdle, StateExercise},
		StateResting:  {StateWorking, StateIdle},
		StateExercise: {StateResting, StateIdle},
	}

	allowedStates, exists := validTransitions[sm.CurrentState]
	if !exists {
		return false
	}

	for _, allowed := range allowedStates {
		if allowed == newState {
			return true
		}
	}
	return false
}
