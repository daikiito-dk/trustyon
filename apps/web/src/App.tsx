import { FormEvent, useEffect, useState } from 'react'

type Health = 'checking' | 'online' | 'offline'
type Project = { id: number; name: string; status: string }
type Task = { id: number; title: string; status: string; priority: string }
type Note = { id: number; title: string; body: string }

async function getJSON<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(path, options)
  if (!response.ok) throw new Error(`API ${response.status}`)
  if (response.status === 204) return undefined as T
  return response.json()
}

function App() {
  const [health, setHealth] = useState<Health>('checking')
  const [projects, setProjects] = useState<Project[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)
  const [newTask, setNewTask] = useState('')
  const [priority, setPriority] = useState('medium')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    Promise.all([
      getJSON<Project[]>('/api/projects'),
      getJSON<Task[]>('/api/tasks'),
      getJSON<Note[]>('/api/notes'),
    ])
      .then(([projectData, taskData, noteData]) => {
        setProjects(projectData)
        setTasks(taskData)
        setNotes(noteData)
        setHealth('online')
      })
      .catch(() => setHealth('offline'))
      .finally(() => setLoading(false))
  }, [])

  const addTask = async (event: FormEvent) => {
    event.preventDefault()
    const title = newTask.trim()
    if (!title || saving) return
    setSaving(true)
    try {
      const task = await getJSON<Task>('/api/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, status: 'todo', priority }),
      })
      setTasks((current) => [task, ...current])
      setNewTask('')
      setHealth('online')
    } catch {
      setHealth('offline')
    } finally {
      setSaving(false)
    }
  }

  const toggleTask = async (task: Task) => {
    const status = task.status === 'done' ? 'todo' : 'done'
    try {
      const updated = await getJSON<Task>(`/api/tasks/${task.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...task, status }),
      })
      setTasks((current) => current.map((item) => item.id === updated.id ? updated : item))
    } catch {
      setHealth('offline')
    }
  }

  const deleteTask = async (id: number) => {
    try {
      await getJSON<void>(`/api/tasks/${id}`, { method: 'DELETE' })
      setTasks((current) => current.filter((task) => task.id !== id))
    } catch {
      setHealth('offline')
    }
  }

  const openTasks = tasks.filter((task) => task.status !== 'done').slice(0, 5)
  const doneTasks = tasks.filter((task) => task.status === 'done').slice(0, 3)

  return (
    <main className="shell">
      <header className="header">
        <div>
          <span className="eyebrow">PERSONAL DEVELOPER OS</span>
          <h1>Trustyon</h1>
        </div>
        <div className={`status ${health}`}><span /> API {health}</div>
      </header>

      <section className="hero">
        <p className="eyebrow">TODAY</p>
        <h2>Build. Learn. Ship.</h2>
        <p className="muted">Your development history, projects, tasks, and notes — in one place.</p>
      </section>

      <section className="grid">
        <article className="card wide">
          <div className="card-title"><span>Tasks</span><strong>{loading ? '—' : tasks.length}</strong></div>
          <form className="task-form" onSubmit={addTask}>
            <input value={newTask} onChange={(event) => setNewTask(event.target.value)} placeholder="What needs to ship?" aria-label="New task" />
            <select value={priority} onChange={(event) => setPriority(event.target.value)} aria-label="Priority">
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
            <button type="submit" disabled={saving || !newTask.trim()}>{saving ? 'Adding…' : 'Add task'}</button>
          </form>
          {openTasks.length > 0 ? (
            <ul className="task-list">
              {openTasks.map((task) => (
                <li key={task.id}>
                  <button className="check" onClick={() => toggleTask(task)} aria-label={`Complete ${task.title}`}>□</button>
                  <span className="task-name">{task.title}</span>
                  <small>{task.priority}</small>
                  <button className="delete" onClick={() => deleteTask(task.id)} aria-label={`Delete ${task.title}`}>×</button>
                </li>
              ))}
            </ul>
          ) : <p className="muted">{loading ? 'Loading tasks…' : 'No open tasks. Add your next move above.'}</p>}
          {doneTasks.length > 0 && (
            <div className="completed">
              <p className="eyebrow">COMPLETED</p>
              {doneTasks.map((task) => (
                <div className="completed-task" key={task.id}>
                  <button className="check" onClick={() => toggleTask(task)} aria-label={`Reopen ${task.title}`}>✓</button>
                  <span>{task.title}</span>
                  <button className="delete" onClick={() => deleteTask(task.id)} aria-label={`Delete ${task.title}`}>×</button>
                </div>
              ))}
            </div>
          )}
        </article>

        <article className="card">
          <div className="card-title"><span>Projects</span><strong>{loading ? '—' : projects.length}</strong></div>
          {projects.length > 0 ? (
            <ul className="simple-list">{projects.slice(0, 5).map((project) => <li key={project.id}><span>●</span>{project.name}</li>)}</ul>
          ) : <p className="muted">{loading ? 'Loading projects…' : 'No projects yet.'}</p>}
        </article>

        <article className="card">
          <div className="card-title"><span>Notes</span><strong>{loading ? '—' : notes.length}</strong></div>
          {notes.length > 0 ? (
            <ul className="simple-list">{notes.slice(0, 5).map((note) => <li key={note.id}><span>↳</span>{note.title}</li>)}</ul>
          ) : <p className="muted">{loading ? 'Loading notes…' : 'Capture ideas and technical notes.'}</p>}
        </article>
      </section>
    </main>
  )
}

export default App
