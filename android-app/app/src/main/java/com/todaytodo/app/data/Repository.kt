package com.todaytodo.app.data

import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class TodayTodoRepository(
    private val apiService: ApiService
) {
    suspend fun getTodos(): List<TodoItem> = apiService.getTodos()

    suspend fun createTodo(title: String, description: String, priority: String) {
        apiService.createTodo(CreateTodoRequest(title = title, description = description, priority = priority))
    }

    suspend fun toggleTodo(todo: TodoItem) {
        apiService.updateTodo(
            id = todo.id,
            request = UpdateTodoRequest(completed = !todo.completed)
        )
    }

    suspend fun updateTodoText(todo: TodoItem, title: String, description: String) {
        apiService.updateTodo(
            id = todo.id,
            request = UpdateTodoRequest(title = title, description = description)
        )
    }

    suspend fun deleteTodo(id: Long) {
        apiService.deleteTodo(id)
    }

    suspend fun recordWater(amount: Int) {
        apiService.recordWater(WaterRequest(amount = amount))
    }

    suspend fun recordStand(durationSeconds: Int = 300) {
        apiService.recordStand(StandRequest(duration = durationSeconds))
    }

    suspend fun recordShortVideo(count: Int = 1) {
        apiService.recordShortVideo(ShortVideoRequest(count = count))
    }

    suspend fun getDailyProgress(): DailyProgressData {
        return apiService.getDailyProgress().data
    }

    suspend fun getReminderConfig(): Map<String, Int> {
        return apiService.getReminderConfig().data
    }

    companion object {
        // Android 模拟器访问主机服务请使用 10.0.2.2。
        private const val BASE_URL = "http://10.0.2.2:8080/"

        fun createDefault(): TodayTodoRepository {
            val logger = HttpLoggingInterceptor().apply {
                level = HttpLoggingInterceptor.Level.BASIC
            }

            val okHttp = OkHttpClient.Builder()
                .addInterceptor(logger)
                .build()

            val retrofit = Retrofit.Builder()
                .baseUrl(BASE_URL)
                .addConverterFactory(GsonConverterFactory.create())
                .client(okHttp)
                .build()

            val api = retrofit.create(ApiService::class.java)
            return TodayTodoRepository(api)
        }
    }
}
