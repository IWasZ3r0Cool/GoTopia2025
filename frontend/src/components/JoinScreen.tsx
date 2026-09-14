import { useState, type FormEvent } from 'react'
import type { ConnectionStatus } from '../types/game'

interface JoinScreenProps {
  status: ConnectionStatus
  error: string | null
  onJoin: (name: string) => void
}

export function JoinScreen({ status, error, onJoin }: JoinScreenProps) {
  const [name, setName] = useState('')
  const trimmedName = name.trim()

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (trimmedName) onJoin(trimmedName)
  }

  return (
    <main className="join-screen">
      <section className="join-card" aria-labelledby="game-title">
        <div className="brand-mark" aria-hidden="true">G</div>
        <p className="eyebrow">A four-island strategy game</p>
        <h1 id="game-title">GoTopia</h1>
        <p className="join-copy">
          Grow a small island into a thriving society. Build wisely, earn gold each turn,
          and keep an eye on the rulers across the water.
        </p>
        <form onSubmit={submit} className="join-form">
          <label htmlFor="player-name">Ruler name</label>
          <div className="join-controls">
            <input
              id="player-name"
              value={name}
              onChange={(event) => setName(event.target.value)}
              maxLength={24}
              autoComplete="nickname"
              autoFocus
              placeholder="Enter your name"
              disabled={status === 'connecting'}
            />
            <button type="submit" className="primary-button" disabled={!trimmedName || status === 'connecting'}>
              {status === 'connecting' ? 'Joining…' : 'Set sail'}
            </button>
          </div>
        </form>
        {error && <p className="error-banner" role="alert">{error}</p>}
        <div className="join-rules" aria-label="Game basics">
          <span><strong>4</strong> rulers</span>
          <span><strong>60s</strong> turns</span>
          <span><strong>1,000</strong> starting gold</span>
        </div>
      </section>
    </main>
  )
}
