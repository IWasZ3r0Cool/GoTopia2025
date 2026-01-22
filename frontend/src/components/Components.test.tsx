import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi as jest } from 'vitest'
import { Tile } from './Tile'
import { ControlPanel } from './ControlPanel'

describe('Tile Component', () => {
    it('renders without crashing', () => {
        render(<Tile x={10} y={10} onBuild={() => { }} isSelected={false} />);
        const tile = screen.getByRole('button', { name: /Tile 10,10/i });
        expect(tile).toBeInTheDocument();
    });

    it('shows building icon', () => {
        // @ts-ignore
        render(<Tile x={10} y={10} building={{ type: "HOUSE", health: 100 }} onBuild={() => { }} isSelected={false} />);
        const tile = screen.getByRole('button', { name: /Tile 10,10/i });
        expect(tile).toHaveTextContent('🏠');
    });

    it('calls onBuild when clicked', () => {
        const handleClick = jest.fn();
        render(<Tile x={10} y={10} onBuild={handleClick} isSelected={false} />);
        const tile = screen.getByRole('button', { name: /Tile 10,10/i });
        fireEvent.click(tile);
        expect(handleClick).toHaveBeenCalledWith(10, 10);
    });
});

describe('ControlPanel Component', () => {
    const mockPlayer = {
        id: "P1",
        gold: 1000,
        population: 50,
        mood: 100
    };

    it('renders player stats', () => {
        render(<ControlPanel player={mockPlayer} selectedTile={null} onBuild={() => { }} />);
        expect(screen.getByTestId('stat-gold')).toHaveTextContent('1000');
    });

    it('disables build buttons when no tile selected', () => {
        render(<ControlPanel player={mockPlayer} selectedTile={null} onBuild={() => { }} />);
        const houseBtn = screen.getByTestId('build-btn-HOUSE');
        expect(houseBtn).toBeDisabled();
    });

    it('enables build buttons when tile selected', () => {
        render(<ControlPanel player={mockPlayer} selectedTile={{ x: 10, y: 10 }} onBuild={() => { }} />);
        const houseBtn = screen.getByTestId('build-btn-HOUSE');
        expect(houseBtn).not.toBeDisabled();
    });

    it('disables expensive buttons', () => {
        const poorPlayer = { ...mockPlayer, gold: 0 };
        render(<ControlPanel player={poorPlayer} selectedTile={{ x: 10, y: 10 }} onBuild={() => { }} />);
        const houseBtn = screen.getByTestId('build-btn-HOUSE'); // Cost 150
        expect(houseBtn).toBeDisabled();
    });
});
