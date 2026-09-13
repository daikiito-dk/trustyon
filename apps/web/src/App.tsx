import { useEffect, useState } from 'react'

type Health = 'checking' | 'online' | 'offline'

function App() {
  const [health, setHealth] = useState<Health>('checking')

  useEffect(() => {
    fetch('/api/healthz')
      .then((response) => {
        if (!response.ok) throw new Error('API error')
        setHealth('online')
      })
      .catch(() => setHealth('offline'))
  }, [])

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
          <div className="card-title"><span>Tasks</span><strong>0</strong></div>
          <p className="muted">No tasks yet. Your next commit starts here.</p>
        </article>
        <article className="card">
          <div className="card-title"><span>Projects</span><strong>0</strong></div>
          <p className="muted">Projects will appear here.</p>
        </article>
        <article className="card">
          <div className="card-title"><span>Notes</span><strong>0</strong></div>
          <p className="muted">Capture ideas and technical notes.</p>
        </article>
      </section>
    </main>
  )
}

export default App
