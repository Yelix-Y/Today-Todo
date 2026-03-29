package com.todaytodo.app.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.todaytodo.app.data.DailyProgressData
import com.todaytodo.app.data.TodayTodoRepository
import com.todaytodo.app.data.TodoItem
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

data class ReminderBanner(
    val type: String,
    val title: String,
    val message: String
)

data class MainUiState(
    val todos: List<TodoItem> = emptyList(),
    val stats: DailyProgressData? = null,
    val loading: Boolean = false,
    val error: String? = null,
    val activeReminder: ReminderBanner? = null,
    val reminderConfig: Map<String, Int> = emptyMap()
)

class MainViewModel(
    private val repository: TodayTodoRepository = TodayTodoRepository.createDefault()
) : ViewModel() {

    private val _uiState = MutableStateFlow(MainUiState())
    val uiState: StateFlow<MainUiState> = _uiState.asStateFlow()

    init {
        refreshAll()
    }

    fun refreshAll() {
        viewModelScope.launch {
            _uiState.update { it.copy(loading = true, error = null) }
            runCatching {
                val todos = repository.getTodos()
                val stats = repository.getDailyProgress()
                val config = repository.getReminderConfig()
                Triple(todos, stats, config)
            }.onSuccess { result ->
                _uiState.update {
                    it.copy(
                        todos = result.first,
                        stats = result.second,
                        reminderConfig = result.third,
                        loading = false,
                        error = null
                    )
                }
            }.onFailure { error ->
                _uiState.update {
                    it.copy(loading = false, error = error.message ?: "请求失败")
                }
            }
        }
    }

    fun addTodo(title: String, description: String, priority: String = "medium") {
        if (title.isBlank()) return

        viewModelScope.launch {
            runCatching {
                repository.createTodo(title, description, priority)
            }.onSuccess {
                refreshAll()
            }.onFailure { e ->
                _uiState.update { it.copy(error = e.message) }
            }
        }
    }

    fun toggleTodo(todo: TodoItem) {
        viewModelScope.launch {
            runCatching {
                repository.toggleTodo(todo)
            }.onSuccess {
                refreshAll()
            }.onFailure { e ->
                _uiState.update { it.copy(error = e.message) }
            }
        }
    }

    fun updateTodo(todo: TodoItem, title: String, description: String) {
        viewModelScope.launch {
            runCatching {
                repository.updateTodoText(todo, title, description)
            }.onSuccess {
                refreshAll()
            }.onFailure { e ->
                _uiState.update { it.copy(error = e.message) }
            }
        }
    }

    fun deleteTodo(id: Long) {
        viewModelScope.launch {
            runCatching {
                repository.deleteTodo(id)
            }.onSuccess {
                refreshAll()
            }.onFailure { e ->
                _uiState.update { it.copy(error = e.message) }
            }
        }
    }

    fun recordWater(amount: Int) {
        viewModelScope.launch {
            runCatching {
                repository.recordWater(amount)
            }.onSuccess {
                refreshAll()
                showReminder(
                    type = "water",
                    title = "补水已记录",
                    message = "继续保持，目标 2000ml。"
                )
            }.onFailure { e ->
                _uiState.update { it.copy(error = e.message) }
            }
        }
    }

    fun recordStand(durationSeconds: Int = 300) {
        viewModelScope.launch {
            runCatching {
                repository.recordStand(durationSeconds)
            }.onSuccess {
                refreshAll()
                showReminder(
                    type = "stand",
                    title = "站立已记录",
                    message = "你已经打断了久坐，节律很好。"
                )
            }.onFailure { e ->
                _uiState.update { it.copy(error = e.message) }
            }
        }
    }

    fun recordShortVideo(count: Int = 1) {
        viewModelScope.launch {
            runCatching {
                repository.recordShortVideo(count)
            }.onSuccess {
                refreshAll()
                showReminder(
                    type = "short-video",
                    title = "防沉迷提醒",
                    message = "已记录短视频次数，建议立刻回到任务清单。"
                )
            }.onFailure { e ->
                _uiState.update { it.copy(error = e.message) }
            }
        }
    }

    fun showReminder(type: String, title: String, message: String) {
        _uiState.update {
            it.copy(activeReminder = ReminderBanner(type = type, title = title, message = message))
        }
    }

    fun dismissReminder() {
        _uiState.update { it.copy(activeReminder = null) }
    }

    fun clearError() {
        _uiState.update { it.copy(error = null) }
    }
}
