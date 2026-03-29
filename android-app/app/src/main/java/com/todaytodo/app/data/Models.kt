package com.todaytodo.app.data

import com.google.gson.annotations.SerializedName

// TodoItem 对应后端 Todo 模型。
data class TodoItem(
    val id: Long,
    val title: String,
    val description: String,
    val completed: Boolean,
    val priority: String,
    @SerializedName("created_at")
    val createdAt: String? = null
)

// 创建任务请求体。
data class CreateTodoRequest(
    val title: String,
    val description: String,
    val priority: String = "medium"
)

// 更新任务请求体。
data class UpdateTodoRequest(
    val title: String? = null,
    val description: String? = null,
    val priority: String? = null,
    val completed: Boolean? = null
)

data class WaterRequest(
    @SerializedName("user_id")
    val userId: Int = 1,
    val amount: Int
)

data class StandRequest(
    @SerializedName("user_id")
    val userId: Int = 1,
    val duration: Int
)

data class ShortVideoRequest(
    @SerializedName("user_id")
    val userId: Int = 1,
    val count: Int = 1
)

// DailyProgressData 展示当天四类统计。
data class DailyProgressData(
    val date: String,
    @SerializedName("completed_todos")
    val completedTodos: Int,
    @SerializedName("total_todos")
    val totalTodos: Int,
    @SerializedName("water_total")
    val waterTotal: Int,
    @SerializedName("water_target")
    val waterTarget: Int,
    @SerializedName("water_progress")
    val waterProgress: Double,
    @SerializedName("stand_total_minutes")
    val standTotalMinutes: Int,
    @SerializedName("stand_target")
    val standTarget: Int,
    @SerializedName("stand_progress")
    val standProgress: Double,
    @SerializedName("short_video_count")
    val shortVideoCount: Int,
    @SerializedName("focus_score")
    val focusScore: Int
)

data class ApiDataWrapper<T>(
    val data: T
)
