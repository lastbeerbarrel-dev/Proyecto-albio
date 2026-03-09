import React, { useState, useEffect } from 'react';
import { Save, Search, Settings as SettingsIcon } from 'lucide-react';

interface ItemMetadata {
    id: string;
    name: string;
    tier: number;
    base_ip: number;
    reference_price: number;
    prices?: Record<string, number>;
}

const Settings: React.FC = () => {
    const [items, setItems] = useState<Record<string, ItemMetadata>>({});
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState('');
    const [searchTerm, setSearchTerm] = useState('');

    const fetchItems = async () => {
        try {
            const res = await fetch('http://localhost:8082/items');
            if (res.ok) {
                const data = await res.json();
                setItems(data);
            }
        } catch (err) {
            setMessage('Error al cargar objetos');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchItems();
    }, []);

    const handleSave = async () => {
        setSaving(true);
        setMessage('');
        try {
            const res = await fetch('http://localhost:8082/items', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(items)
            });
            if (res.ok) {
                setMessage('¡Configuración guardada con éxito!');
                setTimeout(() => setMessage(''), 3000);
            } else {
                setMessage('Error al guardar la configuración');
            }
        } catch (err) {
            setMessage('Error de red al guardar');
        } finally {
            setSaving(false);
        }
    };

    const updateRegionalPrice = (id: string, region: string, value: number) => {
        setItems(prev => {
            const item = prev[id];
            const updatedPrices = { ...(item.prices || {}), [region]: value };
            return {
                ...prev,
                [id]: { ...item, prices: updatedPrices, reference_price: value } // Keep reference_price in sync with last edit for safety
            };
        });
    };

    const filteredItems = Object.entries(items).filter(([id, item]) =>
        id.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading) return <div style={{ padding: '2rem', textAlign: 'center' }}>Cargando Base de Datos...</div>;

    const RegionInput = ({ id, region, value, flag }: any) => (
        <div style={{ flex: 1 }}>
            <label style={{ display: 'block', fontSize: '0.65rem', color: 'var(--color-text-secondary)', marginBottom: '0.25rem' }}>{flag} {region === 'Americas' ? 'Américas' : region}</label>
            <input
                type="number"
                value={value || 0}
                onChange={(e) => updateRegionalPrice(id, region, Number(e.target.value))}
                style={{
                    width: '100%',
                    padding: '0.4rem',
                    background: 'rgba(0,0,0,0.3)',
                    border: '1px solid var(--glass-border)',
                    color: 'white',
                    borderRadius: '4px',
                    fontSize: '0.8rem',
                    outline: 'none'
                }}
            />
        </div>
    );

    return (
        <div style={{ padding: '2rem', maxWidth: '1000px', margin: '0 auto' }}>
            <div className="card" style={{ marginBottom: '1.5rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                    <SettingsIcon size={24} className="text-secondary" />
                    <div>
                        <h2 style={{ margin: 0, fontSize: '1.25rem' }}>Precios de Referencia Regionales</h2>
                        <p style={{ margin: '0.25rem 0 0 0', color: 'var(--color-text-secondary)', fontSize: '0.875rem' }}>
                            Actualiza los precios de venta objetivo para Américas, Europa y Asia.
                        </p>
                    </div>
                </div>
                <button
                    className="btn"
                    onClick={handleSave}
                    disabled={saving}
                    style={{ background: 'var(--color-primary)', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
                >
                    <Save size={18} />
                    {saving ? 'Guardando...' : 'Guardar Todos los Cambios'}
                </button>
            </div>

            {message && (
                <div style={{
                    padding: '1rem',
                    background: message.includes('éxito') ? 'rgba(34, 197, 94, 0.1)' : 'rgba(239, 68, 68, 0.1)',
                    color: message.includes('éxito') ? 'var(--color-success)' : 'var(--color-danger)',
                    borderRadius: '0.5rem',
                    marginBottom: '1rem',
                    textAlign: 'center',
                    fontWeight: 600,
                    border: '1px solid currentColor'
                }}>
                    {message}
                </div>
            )}

            <div className="card" style={{ padding: '0.75rem', marginBottom: '1.5rem' }}>
                <div style={{ position: 'relative', display: 'flex', alignItems: 'center' }}>
                    <Search size={18} style={{ position: 'absolute', left: '1rem', color: 'var(--color-text-secondary)' }} />
                    <input
                        type="text"
                        placeholder="Buscar por ID o Nombre (ej. 'T4 Bag' o 'Metal Bar')..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        style={{
                            width: '100%',
                            padding: '0.75rem 0.75rem 0.75rem 3rem',
                            background: 'rgba(0,0,0,0.2)',
                            border: '1px solid var(--glass-border)',
                            color: 'white',
                            borderRadius: '0.5rem',
                            outline: 'none'
                        }}
                    />
                </div>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                {filteredItems.map(([id, item]) => (
                    <div key={id} className="card" style={{ background: 'rgba(255,255,255,0.02)', border: '1px solid var(--glass-border)', padding: '1rem' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem', alignItems: 'center' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                                <img
                                    src={`https://render.albiononline.com/v1/item/${id}.png?size=32`}
                                    alt={id}
                                    style={{ width: '32px', height: '32px', borderRadius: '4px', background: 'rgba(255,255,255,0.05)' }}
                                />
                                <div>
                                    <div style={{ fontSize: '0.7rem', color: 'var(--color-accent)', fontWeight: 600 }}>{id}</div>
                                    <div style={{ fontWeight: 600, fontSize: '0.9rem' }}>{item.name}</div>
                                </div>
                            </div>
                            <div style={{ fontSize: '0.75rem', padding: '0.2rem 0.5rem', background: 'rgba(255,255,255,0.05)', borderRadius: '4px' }}>
                                T{item.tier}
                            </div>
                        </div>

                        <div style={{ display: 'flex', gap: '0.75rem' }}>
                            <RegionInput id={id} region="Americas" value={item.prices?.Americas} flag="🇺🇸" />
                            <RegionInput id={id} region="Europe" value={item.prices?.Europe} flag="🇪🇺" />
                            <RegionInput id={id} region="Asia" value={item.prices?.Asia} flag="🇸🇬" />
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default Settings;
