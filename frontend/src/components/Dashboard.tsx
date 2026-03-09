import React, { useState, useEffect } from 'react';
import { LayoutDashboard, Activity, Wifi, DollarSign, TrendingUp, Cpu, Globe, ArrowRightLeft, Settings as SettingsIcon } from 'lucide-react';
import OpportunityTable from './OpportunityTable';
import TradeRoutesTable from './TradeRoutesTable';
import Settings from './settings/Settings';

export type Region = 'Americas' | 'Europe' | 'Asia';
export type View = 'dashboard' | 'opportunities' | 'routes' | 'settings';

const Dashboard: React.FC = () => {
    const [view, setView] = useState<View>('dashboard');
    const [region, setRegion] = useState<Region>((localStorage.getItem('albion_region') as Region) || 'Americas');
    const [health, setHealth] = useState({ ingestion: false, calculation: false, metadata: false });
    const [stats, setStats] = useState({ total_profit: 0, active_flips: 0 });

    // Persist Region
    const handleRegionChange = (r: Region) => {
        setRegion(r);
        localStorage.setItem('albion_region', r);
    };

    // Poll Health & Stats
    useEffect(() => {
        const fetchData = async () => {
            try {
                const i = await fetch('http://localhost:8080/health').then(r => r.ok).catch(() => false);
                const c = await fetch('http://localhost:8081/health').then(r => r.ok).catch(() => false);
                const m = await fetch('http://localhost:8082/health').then(r => r.ok).catch(() => false);
                setHealth({ ingestion: i, calculation: c, metadata: m });

                if (c) {
                    const s = await fetch('http://localhost:8081/stats').then(r => r.json()).catch(() => ({ total_profit: 0, active_flips: 0 }));
                    setStats(s);
                }
            } catch (e) {
                // ignore
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 5000);
        return () => clearInterval(interval);
    }, []);

    const formatPrice = (p: number) => {
        if (p === undefined || p === null || isNaN(p)) return '0';
        if (p >= 1000000) return (p / 1000000).toFixed(1) + 'M';
        if (p >= 1000) return (p / 1000).toFixed(1) + 'k';
        return p.toFixed(0);
    };

    const allSystemsGo = health.ingestion && health.calculation && health.metadata;

    const NavItem = ({ icon: Icon, label, id }: any) => (
        <button
            onClick={() => setView(id)}
            style={{
                background: 'none',
                border: 'none',
                color: view === id ? 'var(--color-text-primary)' : 'var(--color-text-secondary)',
                fontWeight: view === id ? 600 : 400,
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
                padding: '0.5rem',
                borderBottom: view === id ? '2px solid var(--color-accent)' : '2px solid transparent'
            }}
        >
            <Icon size={18} />
            {label}
        </button>
    );

    return (
        <div style={{ minHeight: '100vh', background: 'var(--color-bg-primary)', color: 'var(--color-text-primary)' }}>
            <header style={{
                padding: '1rem 2rem',
                borderBottom: '1px solid var(--glass-border)',
                background: 'var(--glass-bg)',
                backdropFilter: 'blur(10px)',
                display: 'flex',
                alignItems: 'center',
                gap: '2rem'
            }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <div style={{
                        width: '32px',
                        height: '32px',
                        background: 'linear-gradient(135deg, var(--color-primary), var(--color-accent))',
                        borderRadius: '8px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}>
                        <Activity size={20} color="white" />
                    </div>
                    <h1 style={{ fontSize: '1.25rem', fontWeight: 700, letterSpacing: '-0.025em', marginRight: '2rem' }}>Albion Analytics <span style={{ opacity: 0.5, fontWeight: 400 }}>Pro</span></h1>
                </div>

                <nav style={{ display: 'flex', gap: '1rem' }}>
                    <NavItem icon={LayoutDashboard} label="Escritorio" id="dashboard" />
                    <NavItem icon={TrendingUp} label="Oportunidades" id="opportunities" />
                    <NavItem icon={ArrowRightLeft} label="Rutas" id="routes" />
                    <NavItem icon={SettingsIcon} label="Ajustes" id="settings" />
                </nav>

                <div style={{ marginLeft: 'auto', display: 'flex', gap: '1rem', alignItems: 'center' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', background: 'rgba(255,255,255,0.03)', padding: '0.25rem 0.5rem', borderRadius: '0.5rem', border: '1px solid var(--glass-border)' }}>
                        <Globe size={14} className="text-secondary" />
                        <select
                            value={region}
                            onChange={(e) => handleRegionChange(e.target.value as Region)}
                            style={{ background: 'transparent', border: 'none', color: 'white', fontSize: '0.8rem', outline: 'none', cursor: 'pointer' }}
                        >
                            <option value="Americas" style={{ background: '#1a1d21' }}>🇺🇸 Américas</option>
                            <option value="Europe" style={{ background: '#1a1d21' }}>🇪🇺 Europa</option>
                            <option value="Asia" style={{ background: '#1a1d21' }}>🇸🇬 Asia</option>
                        </select>
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem', color: 'var(--color-success)', background: 'rgba(34, 197, 94, 0.1)', padding: '0.25rem 0.75rem', borderRadius: '99px' }}>
                        <Wifi size={14} />
                        <span>Conectado</span>
                    </div>
                </div>
            </header>

            {view === 'dashboard' ? (
                <main className="container">
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '1rem', marginBottom: '2rem' }}>
                        <div className="card">
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                <div>
                                    <div style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>Ganancia Total (Sesión)</div>
                                    <div style={{ fontSize: '1.5rem', fontWeight: 700, marginTop: '0.25rem' }}>{formatPrice(stats.total_profit)}</div>
                                </div>
                                <div style={{ padding: '0.5rem', background: 'rgba(59, 130, 246, 0.1)', borderRadius: '8px', color: 'var(--color-primary)' }}>
                                    <DollarSign size={20} />
                                </div>
                            </div>
                        </div>
                        <div className="card">
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                <div>
                                    <div style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>Oportunidades Activas</div>
                                    <div style={{ fontSize: '1.5rem', fontWeight: 700, marginTop: '0.25rem' }}>{stats.active_flips}</div>
                                </div>
                                <div style={{ padding: '0.5rem', background: 'rgba(234, 179, 8, 0.1)', borderRadius: '8px', color: 'var(--color-warning)' }}>
                                    <TrendingUp size={20} />
                                </div>
                            </div>
                        </div>
                        <div className="card">
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                <div>
                                    <div style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>Sniffer Privado</div>
                                    <div style={{ fontSize: '1.5rem', fontWeight: 700, marginTop: '0.25rem', color: 'var(--color-success)' }}>Activo</div>
                                </div>
                                <div style={{ padding: '0.5rem', background: 'rgba(34, 197, 94, 0.1)', borderRadius: '8px', color: 'var(--color-success)' }}>
                                    <Cpu size={20} />
                                </div>
                            </div>
                        </div>
                        <div className="card">
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                <div>
                                    <div style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>Estado del Sistema</div>
                                    <div style={{ fontSize: '0.875rem', fontWeight: 600, marginTop: '0.25rem', display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
                                        <span style={{ color: health.ingestion ? 'var(--color-success)' : 'var(--color-danger)' }}>Ingesta: {health.ingestion ? 'OK' : 'OFF'}</span>
                                        <span style={{ color: health.calculation ? 'var(--color-success)' : 'var(--color-danger)' }}>Cálculo: {health.calculation ? 'OK' : 'OFF'}</span>
                                        <span style={{ color: health.metadata ? 'var(--color-success)' : 'var(--color-danger)' }}>Metadatos: {health.metadata ? 'OK' : 'OFF'}</span>
                                    </div>
                                </div>
                                <div style={{ padding: '0.5rem', background: allSystemsGo ? 'rgba(34, 197, 94, 0.1)' : 'rgba(239, 68, 68, 0.1)', borderRadius: '8px', color: allSystemsGo ? 'var(--color-success)' : 'var(--color-danger)' }}>
                                    <Activity size={20} />
                                </div>
                            </div>
                        </div>
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                        <div className="card-header" style={{ marginBottom: '-1rem', paddingLeft: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--color-primary)' }}>
                            <TrendingUp size={18} />
                            <span style={{ fontWeight: 600 }}>Mejores Flips (Black Market)</span>
                        </div>
                        <OpportunityTable region={region} />

                        <div className="card-header" style={{ marginBottom: '-1rem', paddingLeft: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--color-accent)' }}>
                            <ArrowRightLeft size={18} />
                            <span style={{ fontWeight: 600 }}>Rutas de Arbitraje</span>
                        </div>
                        <TradeRoutesTable region={region} />
                    </div>
                </main>
            ) : view === 'opportunities' ? (
                <main className="container" style={{ height: 'calc(100vh - 120px)' }}>
                    <OpportunityTable region={region} />
                </main>
            ) : view === 'routes' ? (
                <main className="container" style={{ height: 'calc(100vh - 120px)' }}>
                    <TradeRoutesTable region={region} />
                </main>
            ) : view === 'settings' ? (
                <main className="container">
                    <Settings />
                </main>
            ) : null}
        </div>
    );
};

export default Dashboard;
