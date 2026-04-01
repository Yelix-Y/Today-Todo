import { useEffect, useMemo, useState } from "react";

const API_BASE = "/api/v1";
const USER_ID = 1;

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  const data = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(data.error || `请求失败: ${response.status}`);
  return data;
}

function mapPriority(priority) {
  if (priority === "high") return "高";
  if (priority === "low") return "低";
  return "中";
}

function mapRisk(level) {
  if (level === "high") return "高";
  if (level === "medium") return "中";
  return "低";
}

export default function App() {
  const [todos, setTodos] = useState([]);
  const [stats, setStats] = useState(null);
  const [insight, setInsight] = useState(null);
  const [form, setForm] = useState({ title: "", description: "", priority: "medium" });
  const [toast, setToast] = useState("");
  const [reminder, setReminder] = useState(null);

  const completionText = useMemo(() => {
    if (!stats) return "0/0";
    return `${stats.completed_todos}/${stats.total_todos}`;
  }, [stats]);

  const showToast = (msg) => {
    setToast(msg);
    window.clearTimeout(window.__todoToastTimer);
    window.__todoToastTimer = window.setTimeout(() => setToast(""), 1800);
  };

  const loadAll = async () => {
    const [todoData, statsData, insightData] = await Promise.all([
      request("/todos"),
      request(`/health/daily-progress?user_id=${USER_ID}`),
      request(`/insights/today?user_id=${USER_ID}`),
    ]);

    setTodos(todoData);
    setStats(statsData.data);
    setInsight(insightData.data);
  };

  useEffect(() => {
    loadAll().catch((err) => showToast(err.message));
  }, []);

  useEffect(() => {
    const source = new EventSource(`${API_BASE}/reminders/stream`);
    source.addEventListener("reminder", (event) => {
      try {
        const payload = JSON.parse(event.data);
        setReminder(payload);
      } catch {
        // ignore malformed event
      }
    });
    source.onerror = () => {
      source.close();
    };

    return () => source.close();
  }, []);

  const submitTodo = async (event) => {
    event.preventDefault();
    if (!form.title.trim()) return;
    await request("/todos", {
      method: "POST",
      body: JSON.stringify({
        title: form.title.trim(),
        description: form.description.trim(),
        priority: form.priority,
      }),
    });
    setForm({ title: "", description: "", priority: "medium" });
    await loadAll();
    showToast("任务已创建");
  };

  const toggleTodo = async (todo) => {
    await request(`/todos/${todo.id}`, {
      method: "PUT",
      body: JSON.stringify({ completed: !todo.completed }),
    });
    await loadAll();
  };

  const deleteTodo = async (id) => {
    await request(`/todos/${id}`, { method: "DELETE" });
    await loadAll();
  };

  const record = async (path, body, msg) => {
    await request(path, {
      method: "POST",
      body: JSON.stringify(body),
    });
    await loadAll();
    showToast(msg);
  };

  return (
    <div className="page">
      <main className="shell app-grid">
        <aside className="rail card">
          <p className="eyebrow">TODAY FOCUS SYSTEM</p>
          <h2>Today-Todo</h2>
          <p className="sub">Chat 风格控制台</p>
          <button className="btn ghost block" onClick={() => loadAll().catch((err) => showToast(err.message))}>刷新数据</button>
          <div className="rail-meta">
            <p>任务完成：{completionText}</p>
            <p>专注评分：{stats ? stats.focus_score : 0}</p>
            <p>风险等级：{insight ? mapRisk(insight.risk_level) : "-"}</p>
          </div>
        </aside>

        <section className="workspace">
          <header className="card topbar">
            <div>
              <h1>Today-Todo React Console</h1>
              <p className="sub">REST + SSE 实时工作流</p>
            </div>
            <button className="btn ghost" onClick={() => loadAll().catch((err) => showToast(err.message))}>刷新</button>
          </header>

          <section className="stats">
            <article className="card stat"><h3>任务完成</h3><p className="big">{completionText}</p></article>
            <article className="card stat"><h3>喝水</h3><p className="big">{stats ? `${stats.water_total}ml` : "0ml"}</p></article>
            <article className="card stat"><h3>站立</h3><p className="big">{stats ? `${stats.stand_total_minutes}分钟` : "0分钟"}</p></article>
            <article className="card stat"><h3>专注评分</h3><p className="big">{stats ? stats.focus_score : 0}</p></article>
          </section>

          <section className="grid">
            <article className="card chat-panel">
              <div className="panel-head"><h2>任务清单</h2></div>
              <form className="form" onSubmit={submitTodo}>
                <input placeholder="任务标题" value={form.title} onChange={(e) => setForm((s) => ({ ...s, title: e.target.value }))} required />
                <input placeholder="任务说明（可选）" value={form.description} onChange={(e) => setForm((s) => ({ ...s, description: e.target.value }))} />
                <select value={form.priority} onChange={(e) => setForm((s) => ({ ...s, priority: e.target.value }))}>
                  <option value="high">高优先级</option>
                  <option value="medium">中优先级</option>
                  <option value="low">低优先级</option>
                </select>
                <button className="btn primary" type="submit">新增任务</button>
              </form>

              <ul className="todo-list">
                {todos.map((todo) => (
                  <li key={todo.id} className={`todo-item ${todo.completed ? "done" : ""}`}>
                    <div>
                      <p className="todo-title">{todo.title} · {mapPriority(todo.priority)}</p>
                      <p className="todo-desc">{todo.description || "暂无说明"}</p>
                    </div>
                    <div className="actions">
                      <button className="btn success" onClick={() => toggleTodo(todo)}>{todo.completed ? "撤销" : "完成"}</button>
                      <button className="btn danger" onClick={() => deleteTodo(todo.id)}>删除</button>
                    </div>
                  </li>
                ))}
              </ul>
            </article>

            <article className="card">
              <div className="panel-head"><h2>健康与洞察</h2></div>
              <div className="actions quick">
                <button className="btn" onClick={() => record("/health/water", { user_id: USER_ID, amount: 250 }, "喝水 +250ml")}>喝水 +250</button>
                <button className="btn" onClick={() => record("/health/water", { user_id: USER_ID, amount: 500 }, "喝水 +500ml")}>喝水 +500</button>
                <button className="btn" onClick={() => record("/health/stand", { user_id: USER_ID, duration: 300 }, "站立 +5分钟")}>站立 +5分钟</button>
                <button className="btn" onClick={() => record("/health/short-video", { user_id: USER_ID, count: 1 }, "记录短视频 1 次")}>短视频 +1</button>
              </div>

              <div className="insight">
                <p>风险等级：<strong>{insight ? mapRisk(insight.risk_level) : "-"}</strong></p>
                <p>节奏判断：<strong>{insight?.momentum || "-"}</strong></p>
                <p>{insight?.suggested_action || "系统正在生成建议..."}</p>
              </div>
            </article>
          </section>
        </section>
      </main>

      {toast ? <div className="toast">{toast}</div> : null}

      {reminder ? (
        <section className="modal-wrap" onClick={() => setReminder(null)}>
          <div className="modal card" onClick={(e) => e.stopPropagation()}>
            <p className="eyebrow">{reminder.type === "short-video" ? "防沉迷提醒" : "健康提醒"}</p>
            <h3>{reminder.title || "提醒"}</h3>
            <p>{reminder.message || "请及时处理"}</p>
            <button className="btn primary" onClick={() => setReminder(null)}>我知道了</button>
          </div>
        </section>
      ) : null}
    </div>
  );
}
