// Header integrated directly into Dashboard for better state sharing if needed in future
// Skipping deletion to avoid errors if imported elsewhere, though it's not.
import React from 'react';

const Header: React.FC = () => {
    return (
        <header style={{
            padding: '1rem 2rem',
            borderBottom: '1px solid var(--glass-border)',
            background: 'var(--glass-bg)',
            backdropFilter: 'blur(var(--glass-blur))',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center'
        }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <h1 style={{ margin: 0, fontSize: '1.25rem', fontWeight: 700, letterSpacing: '-0.025em' }}>
                    Albion<span style={{ color: 'var(--color-accent)' }}>Analytics</span>
                </h1>
                <nav style={{ display: 'flex', gap: '1.5rem', marginLeft: '2rem' }}>
                    <a href="#" style={{ textDecoration: 'none', color: 'var(--color-text-primary)', fontSize: '0.875rem', fontWeight: 500 }}>Market</a>
                    <a href="#" style={{ textDecoration: 'none', color: 'var(--color-text-secondary)', fontSize: '0.875rem', fontWeight: 500, transition: 'color 0.2s' }}>Crafting</a>
                    <a href="#" style={{ textDecoration: 'none', color: 'var(--color-text-secondary)', fontSize: '0.875rem', fontWeight: 500, transition: 'color 0.2s' }}>Builds</a>
                </nav>
            </div>
            <div>
                <button className="btn" style={{ fontSize: '0.875rem' }}>
                    Connect Private Sniffer
                </button>
            </div>
        </header>
    );
};

export default Header;
