package com.todaytodo.app.ui

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.LinearEasing
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.scaleIn
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.weight
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Checkbox
import androidx.compose.material3.ElevatedButton
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilledTonalButton
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.todaytodo.app.data.TodoItem

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TodayTodoApp(viewModel: MainViewModel) {
    val uiState by viewModel.uiState.collectAsStateWithLifecycle()
    val snackState = remember { SnackbarHostState() }

    var titleInput by remember { mutableStateOf("") }
    var descInput by remember { mutableStateOf("") }
    var editingTodo by remember { mutableStateOf<TodoItem?>(null) }
    var editTitle by remember { mutableStateOf("") }
    var editDesc by remember { mutableStateOf("") }

    LaunchedEffect(uiState.error) {
        uiState.error?.let { msg ->
            snackState.showSnackbar(msg)
            viewModel.clearError()
        }
    }

    Scaffold(
        snackbarHost = { SnackbarHost(hostState = snackState) }
    ) { innerPadding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(Color(0xFFF4FBF9))
                .padding(innerPadding)
        ) {
            LazyColumn(
                modifier = Modifier.fillMaxSize(),
                contentPadding = PaddingValues(16.dp),
                verticalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                item {
                    HeaderCard(onRefresh = viewModel::refreshAll)
                }

                item {
                    StatsPanel(uiState = uiState)
                }

                item {
                    AddTodoPanel(
                        title = titleInput,
                        description = descInput,
                        onTitleChange = { titleInput = it },
                        onDescriptionChange = { descInput = it },
                        onSubmit = {
                            viewModel.addTodo(titleInput, descInput)
                            titleInput = ""
                            descInput = ""
                        }
                    )
                }

                item {
                    HealthActionPanel(
                        onWater250 = { viewModel.recordWater(250) },
                        onWater500 = { viewModel.recordWater(500) },
                        onStand = { viewModel.recordStand(300) },
                        onShortVideo = { viewModel.recordShortVideo(1) }
                    )
                }

                items(uiState.todos, key = { it.id }) { todo ->
                    TodoRow(
                        todo = todo,
                        onToggle = { viewModel.toggleTodo(todo) },
                        onDelete = { viewModel.deleteTodo(todo.id) },
                        onEdit = {
                            editingTodo = todo
                            editTitle = todo.title
                            editDesc = todo.description
                        }
                    )
                }
            }

            AnimatedVisibility(
                visible = uiState.activeReminder != null,
                enter = fadeIn() + scaleIn(initialScale = 0.8f),
                exit = fadeOut(),
                modifier = Modifier.align(Alignment.Center)
            ) {
                uiState.activeReminder?.let { reminder ->
                    ReminderOverlay(
                        type = reminder.type,
                        title = reminder.title,
                        message = reminder.message,
                        onDismiss = viewModel::dismissReminder
                    )
                }
            }
        }
    }

    if (editingTodo != null) {
        AlertDialog(
            onDismissRequest = { editingTodo = null },
            confirmButton = {
                TextButton(onClick = {
                    editingTodo?.let { todo ->
                        viewModel.updateTodo(todo, editTitle, editDesc)
                    }
                    editingTodo = null
                }) {
                    Text("保存")
                }
            },
            dismissButton = {
                TextButton(onClick = { editingTodo = null }) {
                    Text("取消")
                }
            },
            title = { Text("编辑任务") },
            text = {
                Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    OutlinedTextField(
                        value = editTitle,
                        onValueChange = { editTitle = it },
                        label = { Text("标题") },
                        modifier = Modifier.fillMaxWidth()
                    )
                    OutlinedTextField(
                        value = editDesc,
                        onValueChange = { editDesc = it },
                        label = { Text("说明") },
                        modifier = Modifier.fillMaxWidth()
                    )
                }
            }
        )
    }
}

@Composable
private fun HeaderCard(onRefresh: () -> Unit) {
    Card(
        shape = RoundedCornerShape(20.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column {
                Text("TODAY FOCUS SYSTEM", color = Color(0xFF0A7F6C), fontWeight = FontWeight.Bold)
                Text("Today-Todo", style = MaterialTheme.typography.headlineSmall)
                Text("任务 + 健康 + 防沉迷", color = Color(0xFF55706A))
            }
            Button(onClick = onRefresh) {
                Text("刷新")
            }
        }
    }
}

@Composable
private fun StatsPanel(uiState: MainUiState) {
    val stats = uiState.stats
    Row(horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.fillMaxWidth()) {
        MiniStatCard(
            title = "任务",
            value = "${stats?.completedTodos ?: 0}/${stats?.totalTodos ?: 0}",
            modifier = Modifier.weight(1f)
        )
        MiniStatCard(
            title = "喝水",
            value = "${stats?.waterTotal ?: 0}ml",
            modifier = Modifier.weight(1f)
        )
    }

    Spacer(modifier = Modifier.height(8.dp))

    Row(horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.fillMaxWidth()) {
        MiniStatCard(
            title = "站立",
            value = "${stats?.standTotalMinutes ?: 0}分钟",
            modifier = Modifier.weight(1f)
        )
        MiniStatCard(
            title = "短视频",
            value = "${stats?.shortVideoCount ?: 0}次",
            modifier = Modifier.weight(1f),
            accent = Color(0xFFE33E3E)
        )
    }

    Spacer(modifier = Modifier.height(8.dp))

    Card(
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color(0xFFEFFFF9))
    ) {
        Column(modifier = Modifier.padding(14.dp)) {
            Text("专注评分", color = Color(0xFF55706A))
            Text(
                text = "${stats?.focusScore ?: 0}",
                style = MaterialTheme.typography.headlineMedium,
                color = Color(0xFF0A7F6C),
                fontWeight = FontWeight.ExtraBold
            )
        }
    }
}

@Composable
private fun MiniStatCard(
    title: String,
    value: String,
    modifier: Modifier = Modifier,
    accent: Color = Color(0xFF16302B)
) {
    Card(
        modifier = modifier,
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White)
    ) {
        Column(modifier = Modifier.padding(12.dp)) {
            Text(title, color = Color(0xFF55706A))
            Text(value, fontWeight = FontWeight.ExtraBold, color = accent)
        }
    }
}

@Composable
private fun AddTodoPanel(
    title: String,
    description: String,
    onTitleChange: (String) -> Unit,
    onDescriptionChange: (String) -> Unit,
    onSubmit: () -> Unit
) {
    Card(shape = RoundedCornerShape(16.dp)) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(12.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Text("新增任务", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
            OutlinedTextField(
                value = title,
                onValueChange = onTitleChange,
                label = { Text("任务标题") },
                modifier = Modifier.fillMaxWidth()
            )
            OutlinedTextField(
                value = description,
                onValueChange = onDescriptionChange,
                label = { Text("任务说明") },
                modifier = Modifier.fillMaxWidth()
            )
            ElevatedButton(onClick = onSubmit, modifier = Modifier.fillMaxWidth()) {
                Text("添加任务")
            }
        }
    }
}

@Composable
private fun HealthActionPanel(
    onWater250: () -> Unit,
    onWater500: () -> Unit,
    onStand: () -> Unit,
    onShortVideo: () -> Unit
) {
    Card(shape = RoundedCornerShape(16.dp)) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(12.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Text("健康与防沉迷", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.fillMaxWidth()) {
                FilledTonalButton(onClick = onWater250, modifier = Modifier.weight(1f)) {
                    Text("+250ml")
                }
                FilledTonalButton(onClick = onWater500, modifier = Modifier.weight(1f)) {
                    Text("+500ml")
                }
            }
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.fillMaxWidth()) {
                FilledTonalButton(onClick = onStand, modifier = Modifier.weight(1f)) {
                    Text("站立 5 分钟")
                }
                FilledTonalButton(onClick = onShortVideo, modifier = Modifier.weight(1f)) {
                    Text("刷短视频 +1")
                }
            }
        }
    }
}

@Composable
private fun TodoRow(
    todo: TodoItem,
    onToggle: () -> Unit,
    onDelete: () -> Unit,
    onEdit: () -> Unit
) {
    val doneAlpha = if (todo.completed) 0.6f else 1f
    Card(shape = RoundedCornerShape(14.dp), colors = CardDefaults.cardColors(containerColor = Color.White)) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(10.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Checkbox(checked = todo.completed, onCheckedChange = { onToggle() })
            Column(modifier = Modifier.weight(1f).alpha(doneAlpha)) {
                Text(todo.title, fontWeight = FontWeight.Bold)
                if (todo.description.isNotBlank()) {
                    Text(todo.description, color = Color(0xFF55706A))
                }
            }
            IconButton(onClick = onEdit) {
                Icon(Icons.Filled.Edit, contentDescription = "编辑")
            }
            IconButton(onClick = onDelete) {
                Icon(Icons.Filled.Delete, contentDescription = "删除")
            }
        }
    }
}

@Composable
private fun ReminderOverlay(
    type: String,
    title: String,
    message: String,
    onDismiss: () -> Unit
) {
    val infinite = rememberInfiniteTransition(label = "reminder")
    val alpha by infinite.animateFloat(
        initialValue = 0.65f,
        targetValue = 1f,
        animationSpec = infiniteRepeatable(
            animation = tween(durationMillis = 800, easing = LinearEasing),
            repeatMode = RepeatMode.Reverse
        ),
        label = "overlay-alpha"
    )

    val accent = if (type == "short-video") Color(0xFFE33E3E) else Color(0xFF0EA58C)

    Card(
        modifier = Modifier
            .padding(20.dp)
            .fillMaxWidth()
            .alpha(alpha)
            .border(2.dp, accent, RoundedCornerShape(18.dp)),
        shape = RoundedCornerShape(18.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White)
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Box(
                    modifier = Modifier
                        .size(10.dp)
                        .background(accent, CircleShape)
                )
                Spacer(modifier = Modifier.size(8.dp))
                Text("实时提醒", color = accent, fontWeight = FontWeight.Bold)
            }
            Text(title, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.ExtraBold)
            Text(message, color = Color(0xFF55706A))
            Button(onClick = onDismiss, modifier = Modifier.fillMaxWidth()) {
                Text("我已处理")
            }
        }
    }
}
