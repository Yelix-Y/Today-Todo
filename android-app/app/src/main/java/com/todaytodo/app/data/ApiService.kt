package com.todaytodo.app.data

import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path
import retrofit2.http.Query

interface ApiService {
    @GET("api/v1/todos")
    suspend fun getTodos(): List<TodoItem>

    @POST("api/v1/todos")
    suspend fun createTodo(@Body request: CreateTodoRequest): TodoItem

    @PUT("api/v1/todos/{id}")
    suspend fun updateTodo(
        @Path("id") id: Long,
        @Body request: UpdateTodoRequest
    ): TodoItem

    @DELETE("api/v1/todos/{id}")
    suspend fun deleteTodo(@Path("id") id: Long)

    @POST("api/v1/health/water")
    suspend fun recordWater(@Body request: WaterRequest)

    @POST("api/v1/health/stand")
    suspend fun recordStand(@Body request: StandRequest)

    @POST("api/v1/health/short-video")
    suspend fun recordShortVideo(@Body request: ShortVideoRequest)

    @GET("api/v1/health/daily-progress")
    suspend fun getDailyProgress(
        @Query("user_id") userId: Int = 1
    ): ApiDataWrapper<DailyProgressData>

    @GET("api/v1/reminders/config")
    suspend fun getReminderConfig(): ApiDataWrapper<Map<String, Int>>
}
