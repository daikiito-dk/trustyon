import { useEffect, useState } from 'react'

type Health = 'checking' | 'online' | 'offline'
type Project = { id: number; name: string; status: string }
type Task = { id: number; title: string; status: string; priority: string }
type Note = { id: number; title: string; body: string }

async function getJSON<T>(path: string): Promise<T> {
  const response = await fetch(path)
  if (!response.ok) throw new Error(`API ${response.status}`)
  return response.json()
}

function App() {
  const [health, setHealth] = useState<Health>('checking')
  const [projects, setProjects] = useState<Project[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)

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

  const openTasks = tasks.filter((task) => task.status !== 'done').slice(0, 5)

  return (
    <main className="shell">
      <header className="header">
        <div>
          <span className="eyebrow">PERSONAL DEVELOPER OS</span>
          <h1>Trustyon</h1>
        </div>
        <div className={`status ${health}`}>
          <span /> API {health}
        </div>
      </header>

      <section className="hero">
        <p className="eyebrow">TODAY</p>
        <h2>Build. Learn. Ship.</h2>
        <p className="muted">Your development history, projects, tasks, and notes — in one place.</p>
      </section>

      <section className="grid">
        <article className="card wide">
          <div className="card-title"><span>Tasks</span><strong>{loading ? '—' : tasks.length}</strong></div>
          {openTasks.length > 0 ? (
            <ul className="task-list">
              {openTasks.map((task) => <li key={task.id}><span>□</span>{task.title}<small>{task.priority}</small></li>)}
            </ul>
          ) : <p className="muted">{loading ? 'Loading tasks…' : 'No open tasks. Your next commit starts here.'}</p>}
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
