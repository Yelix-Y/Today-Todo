package com.todaytodo.app.reminder

import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.os.Build
import androidx.core.app.NotificationCompat
import androidx.core.app.NotificationManagerCompat
import androidx.work.Constraints
import androidx.work.CoroutineWorker
import androidx.work.Data
import androidx.work.ExistingPeriodicWorkPolicy
import androidx.work.NetworkType
import androidx.work.PeriodicWorkRequestBuilder
import androidx.work.WorkManager
import androidx.work.WorkerParameters
import com.todaytodo.app.R
import java.util.concurrent.TimeUnit

class ReminderWorker(
    appContext: Context,
    workerParams: WorkerParameters
) : CoroutineWorker(appContext, workerParams) {

    override suspend fun doWork(): Result {
        val type = inputData.getString(KEY_TYPE) ?: "general"
        val title = inputData.getString(KEY_TITLE) ?: "提醒"
        val message = inputData.getString(KEY_MESSAGE) ?: "该行动了"

        ensureChannel()
        showNotification(type = type, title = title, message = message)
        return Result.success()
    }

    private fun ensureChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "Today Todo Reminders",
                NotificationManager.IMPORTANCE_HIGH
            ).apply {
                description = "健康与防沉迷提醒"
            }

            val manager = applicationContext.getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
            manager.createNotificationChannel(channel)
        }
    }

    private fun showNotification(type: String, title: String, message: String) {
        val icon = when (type) {
            "short-video" -> android.R.drawable.ic_dialog_alert
            "stand" -> android.R.drawable.ic_menu_compass
            else -> android.R.drawable.ic_menu_info_details
        }

        val notification = NotificationCompat.Builder(applicationContext, CHANNEL_ID)
            .setSmallIcon(icon)
            .setContentTitle(title)
            .setContentText(message)
            .setPriority(NotificationCompat.PRIORITY_HIGH)
            .setAutoCancel(true)
            .build()

        runCatching {
            NotificationManagerCompat.from(applicationContext).notify(type.hashCode(), notification)
        }
    }

    companion object {
        private const val CHANNEL_ID = "today_todo_reminders"
        private const val KEY_TYPE = "type"
        private const val KEY_TITLE = "title"
        private const val KEY_MESSAGE = "message"

        fun scheduleDefaultReminders(context: Context) {
            schedule(
                context = context,
                uniqueName = "water-reminder",
                intervalMinutes = 90,
                inputData = workData("water", "补水提醒", "建议喝 200-300ml 水，保持专注。")
            )
            schedule(
                context = context,
                uniqueName = "stand-reminder",
                intervalMinutes = 60,
                inputData = workData("stand", "站立提醒", "久坐已满 1 小时，建议起身拉伸。")
            )
            schedule(
                context = context,
                uniqueName = "short-video-reminder",
                intervalMinutes = 120,
                inputData = workData("short-video", "防沉迷提醒", "限制短视频时长，优先推进任务。")
            )
        }

        private fun schedule(context: Context, uniqueName: String, intervalMinutes: Long, inputData: Data) {
            val constraints = Constraints.Builder()
                .setRequiredNetworkType(NetworkType.CONNECTED)
                .build()

            val request = PeriodicWorkRequestBuilder<ReminderWorker>(intervalMinutes, TimeUnit.MINUTES)
                .setInputData(inputData)
                .setConstraints(constraints)
                .build()

            WorkManager.getInstance(context).enqueueUniquePeriodicWork(
                uniqueName,
                ExistingPeriodicWorkPolicy.UPDATE,
                request
            )
        }

        private fun workData(type: String, title: String, message: String): Data {
            return Data.Builder()
                .putString(KEY_TYPE, type)
                .putString(KEY_TITLE, title)
                .putString(KEY_MESSAGE, message)
                .build()
        }
    }
}
