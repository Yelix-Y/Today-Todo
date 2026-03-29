const API_BASE = "/api/v1";
const USER_ID = 1;

const state = {
  todos: [],
  reminderConfig: {
    water_minutes: 90,
    stand_minutes: 60,
    short_video_minutes: 120,
  },
};

const el = {
  todoForm: document.getElementById("todoForm"),
  todoTitle: document.getElementById("todoTitle"),
  todoDescription: document.getElementById("todoDescription"),
  todoPriority: document.getElementById("todoPriority"),
  todoList: document.getElementById("todoList"),
  refreshBtn: document.getElementById("refreshBtn"),
  drink250Btn: document.getElementById("drink250Btn"),
  drink500Btn: document.getElementById("drink500Btn"),
  stand5Btn: document.getElementById("stand5Btn"),
  shortVideoBtn: document.getElementById("shortVideoBtn"),
  completedTodos: document.getElementById("completedTodos"),
  totalTodos: document.getElementById("totalTodos"),
  shortVideoCount: document.getElementById("shortVideoCount"),
  waterText: document.getElementById("waterText"),
  standText: document.getElementById("standText"),
  waterRing: document.getElementById("waterRing"),
  standRing: document.getElementById("standRing"),
  focusScore: document.getElementById("focusScore"),
  reminderModal: document.getElementById("reminderModal"),
  reminderType: document.getElementById("reminderType"),
  reminderTitle: document.getElementById("reminderTitle"),
  reminderMessage: document.getElementById("reminderMessage"),
  ackReminder: document.getElementById("ackReminder"),
  quickShortVideo: document.getElementById("quickShortVideo"),
};

// 统一 API 请求工具，简化错误处理。
async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(data.error || `请求失败: ${response.status}`);
  }
  return data;
}

async function loadTodos() {
  const data = await request("/todos");
  state.todos = data;
  renderTodos();
}

function renderTodos() {
  el.todoList.innerHTML = "";

  state.todos.forEach((todo) => {
    const li = document.createElement("li");
    li.className = `todo-item ${todo.completed ? "todo-done" : ""}`;

    const main = document.createElement("div");
    const title = document.createElement("p");
    title.className = "todo-title";
    title.textContent = `${todo.title} · ${mapPriority(todo.priority)}`;

    const desc = document.createElement("p");
    desc.className = "todo-desc";
    desc.textContent = todo.description || "暂无说明";

    main.append(title, desc);

    const actions = document.createElement("div");
    actions.className = "todo-actions";

    const toggleBtn = document.createElement("button");
    toggleBtn.className = "btn btn-success";
    toggleBtn.textContent = todo.completed ? "撤销完成" : "完成";
    toggleBtn.onclick = () => toggleTodo(todo);

    const editBtn = document.createElement("button");
    editBtn.className = "btn btn-ghost";
    editBtn.textContent = "编辑";
    editBtn.onclick = () => editTodo(todo);

    const deleteBtn = document.createElement("button");
    deleteBtn.className = "btn btn-danger";
    deleteBtn.textContent = "删除";
    deleteBtn.onclick = () => removeTodo(todo.id);

    actions.append(toggleBtn, editBtn, deleteBtn);
    li.append(main, actions);
    el.todoList.appendChild(li);
  });
}

async function createTodo(event) {
  event.preventDefault();
  const title = el.todoTitle.value.trim();
  if (!title) return;

  await request("/todos", {
    method: "POST",
    body: JSON.stringify({
      title,
      description: el.todoDescription.value.trim(),
      priority: el.todoPriority.value,
    }),
  });

  el.todoForm.reset();
  await Promise.all([loadTodos(), loadStats()]);
}

async function toggleTodo(todo) {
  await request(`/todos/${todo.id}`, {
    method: "PUT",
    body: JSON.stringify({ completed: !todo.completed }),
  });

  await Promise.all([loadTodos(), loadStats()]);
}

async function editTodo(todo) {
  const title = window.prompt("编辑任务标题", todo.title);
  if (title === null) return;

  const description = window.prompt("编辑任务说明", todo.description || "");
  if (description === null) return;

  await request(`/todos/${todo.id}`, {
    method: "PUT",
    body: JSON.stringify({ title: title.trim(), description: description.trim() }),
  });

  await loadTodos();
}

async function removeTodo(id) {
  await request(`/todos/${id}`, { method: "DELETE" });
  await Promise.all([loadTodos(), loadStats()]);
}

async function recordWater(amount) {
  await request("/health/water", {
    method: "POST",
    body: JSON.stringify({ user_id: USER_ID, amount }),
  });

  await loadStats();
}

async function recordStand(durationSeconds = 300) {
  await request("/health/stand", {
    method: "POST",
    body: JSON.stringify({ user_id: USER_ID, duration: durationSeconds }),
  });

  await loadStats();
}

async function recordShortVideo(count = 1) {
  await request("/health/short-video", {
    method: "POST",
    body: JSON.stringify({ user_id: USER_ID, count }),
  });

  await loadStats();
}

async function loadStats() {
  const result = await request(`/health/daily-progress?user_id=${USER_ID}`);
  const stats = result.data;

  el.completedTodos.textContent = stats.completed_todos;
  el.totalTodos.textContent = stats.total_todos;
  el.shortVideoCount.textContent = stats.short_video_count;
  el.waterText.textContent = `${stats.water_total}ml`;
  el.standText.textContent = `${stats.stand_total_minutes}分钟`;
  el.focusScore.textContent = stats.focus_score;

  el.waterRing.style.setProperty("--p", Number(stats.water_progress || 0).toFixed(1));
  el.standRing.style.setProperty("--p", Number(stats.stand_progress || 0).toFixed(1));
}

function showReminder(reminder) {
  el.reminderType.textContent = reminder.type === "short-video" ? "防沉迷提醒" : "健康提醒";
  el.reminderTitle.textContent = reminder.title || "行动提醒";
  el.reminderMessage.textContent = reminder.message || "请及时完成当前提醒动作";
  el.reminderModal.classList.remove("hidden");
}

function hideReminder() {
  el.reminderModal.classList.add("hidden");
}

function setupSSE() {
  const source = new EventSource(`${API_BASE}/reminders/stream`);

  source.addEventListener("reminder", (event) => {
    try {
      const data = JSON.parse(event.data);
      showReminder(data);
    } catch (error) {
      console.error("提醒解析失败", error);
    }
  });

  source.onerror = () => {
    source.close();
    // SSE 中断后延迟重连，避免频繁请求。
    setTimeout(setupSSE, 5000);
  };
}

async function loadReminderConfig() {
  try {
    const result = await request("/reminders/config");
    state.reminderConfig = {
      ...state.reminderConfig,
      ...result.data,
    };
  } catch (error) {
    console.warn("使用默认提醒配置", error.message);
  }
}

function setupLocalFallbackReminders() {
  // 本地兜底提醒：SSE 不可用时仍然可触发视觉提醒。
  setInterval(() => {
    showReminder({
      type: "water",
      title: "补水提醒",
      message: "喝一杯水，维持注意力与代谢状态。",
    });
  }, state.reminderConfig.water_minutes * 60 * 1000);

  setInterval(() => {
    showReminder({
      type: "stand",
      title: "站立拉伸提醒",
      message: "请离开座位 5 分钟，缓解久坐压力。",
    });
  }, state.reminderConfig.stand_minutes * 60 * 1000);

  setInterval(() => {
    showReminder({
      type: "short-video",
      title: "短视频防沉迷提醒",
      message: "设置短视频时长上限，优先完成今日核心任务。",
    });
  }, state.reminderConfig.short_video_minutes * 60 * 1000);
}

function mapPriority(priority) {
  switch (priority) {
    case "high":
      return "高";
    case "low":
      return "低";
    default:
      return "中";
  }
}

function bindEvents() {
  el.todoForm.addEventListener("submit", (event) => {
    createTodo(event).catch((error) => window.alert(error.message));
  });

  el.refreshBtn.addEventListener("click", () => {
    Promise.all([loadTodos(), loadStats()]).catch((error) => window.alert(error.message));
  });

  el.drink250Btn.addEventListener("click", () => recordWater(250).catch((error) => window.alert(error.message)));
  el.drink500Btn.addEventListener("click", () => recordWater(500).catch((error) => window.alert(error.message)));
  el.stand5Btn.addEventListener("click", () => recordStand(300).catch((error) => window.alert(error.message)));
  el.shortVideoBtn.addEventListener("click", () => {
    recordShortVideo(1).then(() => {
      showReminder({
        type: "short-video",
        title: "已记录短视频行为",
        message: "已计入 1 次，建议回到任务列表继续推进。",
      });
    }).catch((error) => window.alert(error.message));
  });

  el.ackReminder.addEventListener("click", hideReminder);
  el.quickShortVideo.addEventListener("click", () => {
    recordShortVideo(1)
      .then(hideReminder)
      .catch((error) => window.alert(error.message));
  });

  el.reminderModal.addEventListener("click", (event) => {
    if (event.target === el.reminderModal) {
      hideReminder();
    }
  });
}

async function init() {
  bindEvents();
  await loadReminderConfig();
  await Promise.all([loadTodos(), loadStats()]);
  setupSSE();
  setupLocalFallbackReminders();
}

init().catch((error) => {
  console.error(error);
  window.alert(`初始化失败: ${error.message}`);
});
