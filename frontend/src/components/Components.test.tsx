import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { ControlPanel } from './ControlPanel'
import { IslandGrid } from './IslandGrid'
import { JoinScreen } from './JoinScreen'
import { Tile } from './Tile'

const player = { id: 'player-1', name: 'Ada', gold: 1000, population: 100, populationCapacity: 100, mood: 100 }

describe('Tile', () => {
  it('is keyboard-accessible and reports its contents', () => {
    render(<Tile x={10} y={10} building={{ type: 'HOUSE', health: 100 }} onSelect={() => undefined} isSelected />)
    const tile = screen.getByRole('button', { name: 'House at 10, 10' })
    expect(tile).toHaveAttribute('aria-pressed', 'true')
    expect(tile).toHaveTextContent('⌂')
  })

  it('selects its coordinates', () => {
    const onSelect = vi.fn()
    render(<Tile x={12} y={14} onSelect={onSelect} isSelected={false} />)
    fireEvent.click(screen.getByRole('button', { name: 'Empty tile at 12, 14' }))
    expect(onSelect).toHaveBeenCalledWith(12, 14)
  })
})

describe('IslandGrid', () => {
  it('keeps interactive tiles usable inside a scrollable map', () => {
    render(
      <IslandGrid
        island={{ ownerId: 'player-1', x: 10, y: 10, width: 25, height: 25, buildings: {} }}
        ownerName="Ada"
        onTileClick={() => undefined}
      />,
    )

    expect(screen.getByRole('region', { name: 'Scrollable island map' })).toHaveAttribute('tabindex', '0')
    expect(screen.getByTestId('island-grid-player-1')).toHaveStyle({ minWidth: '900px' })
  })
})

describe('ControlPanel', () => {
  it('renders formatted stats and requires a selection', () => {
    render(<ControlPanel player={player} selectedTile={null} onBuild={() => undefined} />)
    expect(screen.getByTestId('stat-gold')).toHaveTextContent('1,000')
    expect(screen.getByTestId('build-btn-HOUSE')).toBeDisabled()
  })

  it('allows an affordable building on an empty selected tile', () => {
    render(<ControlPanel player={player} selectedTile={{ x: 10, y: 10 }} onBuild={() => undefined} />)
    expect(screen.getByTestId('build-btn-HOUSE')).toBeEnabled()
  })

  it('prevents building on an occupied tile or without enough gold', () => {
    const { rerender } = render(
      <ControlPanel player={player} selectedTile={{ x: 10, y: 10 }} selectedBuilding={{ type: 'FARM', health: 100 }} onBuild={() => undefined} />,
    )
    expect(screen.getByTestId('build-btn-HOUSE')).toBeDisabled()

    rerender(<ControlPanel player={{ ...player, gold: 0 }} selectedTile={{ x: 11, y: 10 }} onBuild={() => undefined} />)
    expect(screen.getByTestId('build-btn-HOUSE')).toBeDisabled()
  })
})

describe('JoinScreen', () => {
  it('trims and submits a player name', () => {
    const onJoin = vi.fn()
    render(<JoinScreen status="idle" error={null} onJoin={onJoin} />)
    fireEvent.change(screen.getByLabelText('Ruler name'), { target: { value: '  Grace  ' } })
    fireEvent.click(screen.getByRole('button', { name: 'Set sail' }))
    expect(onJoin).toHaveBeenCalledWith('Grace')
  })

  it('shows connection errors', () => {
    render(<JoinScreen status="idle" error="Server unavailable" onJoin={() => undefined} />)
    expect(screen.getByRole('alert')).toHaveTextContent('Server unavailable')
  })
})
