import React, { useState, useEffect } from 'react';
import { ArrowRightLeft, Flame, Clock, Shield, AlertTriangle } from 'lucide-react';

interface TradeRoute {
    item_id: string;
    item_name: string;
    from_city: string;
    to_city: string;
    buy_price: number;
    sell_price: number;
    profit: number;
    roi: number;
    category: string;
    sub_category: string;
    quality: number;
    confidence_score: number; // Legacy
    data_quality: number;
    execution_score: number;
    risk_profile: string;
    is_filtered: boolean;
    high_volatility?: boolean;
    updated_at: string;
}

interface TradeRoutesTableProps {
    region: string;
}

const getConfidenceColor = (score: number) => {
    if (score >= 70) return '#10B981'; // Green
    if (score >= 35) return '#F59E0B'; // Amber
    return '#EF4444'; // Red
};

const getRiskColor = (risk: string) => {
    switch (risk) {
        case 'Bajo': return '#10B981';
        case 'Medio': return '#F59E0B';
        case 'Alto': return '#EF4444';
        default: return '#94A3B8';
    }
};

const getQualityName = (q: number) => {
    switch (q) {
        case 1: return 'Normal';
        case 2: return 'Bueno';
        case 3: return 'Notable';
        case 4: return 'Sobresaliente';
        case 5: return 'Obra Maestra';
        default: return 'Normal';
    }
};

const getQualityColor = (q: number) => {
    switch (q) {
        case 1: return '#B0C4DE'; // LightSteelBlue for Normal
        case 2: return '#32CD32'; // LimeGreen for Good
        case 3: return '#00BFFF'; // DeepSkyBlue for Outstanding
        case 4: return '#9370DB'; // MediumPurple for Excellent
        case 5: return '#FFD700'; // Gold for Masterpiece
        default: return '#B0C4DE';
    }
};

const calculateTimeSince = (dateString: string) => {
    if (!dateString || dateString.startsWith('0001')) return 'N/A';
    const now = new Date();
    const past = new Date(dateString);
    const diffMs = now.getTime() - past.getTime();

    if (diffMs < 0) return 'Ahora';

    const minutes = Math.floor(diffMs / 60000);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 0) return `${days}d ${hours % 24}h`;
    if (hours > 0) return `${hours}h ${minutes % 60}m`;
    return `${minutes}m`;
};

const TradeRoutesTable: React.FC<TradeRoutesTableProps> = ({ region }) => {
    const [routes, setRoutes] = useState<TradeRoute[]>([]);
    const [selectedCategory, setSelectedCategory] = useState<string>('Todos');
    const [loading, setLoading] = useState(true);

    // Advanced Filters
    const [selectedTiers, setSelectedTiers] = useState<number[]>([4, 5, 6, 7, 8]);
    const [searchTerm, setSearchTerm] = useState<string>('');
    const [minProfit, setMinProfit] = useState<number>(0);
    const [selectedWeaponSubCat, setSelectedWeaponSubCat] = useState<string>('Todas');

    const weaponSubCats = [
        'Todas', 'Espadas', 'Hachas', 'Mazas', 'Martillos', 'Guantes de Guerra',
        'Ballestas', 'Arcos', 'Lanzas', 'Dagas', 'Quarterstaffs',
        'Bastón de Fuego', 'Bastón de Hielo', 'Bastón Maldito', 'Bastón Sagrado', 'Bastón Natural',
        'Bastón Arcano', 'Bastón Cambiaformas'
    ];

    const toggleTier = (t: number) => {
        setSelectedTiers(prev => prev.includes(t) ? prev.filter(x => x !== t) : [...prev, t]);
    };

    const categories = ['Todos', 'Recursos', 'Consumibles', 'Monturas', 'Equipamiento', 'Cosméticos'];

    const categoryMap: Record<string, string[]> = {
        'Todos': ['Resource', 'Consumable', 'Mount', 'Weapon', 'Armor', 'Accessory', 'Offhand', 'Tome', 'Skin', 'Luxury'],
        'Recursos': ['Resource'],
        'Consumibles': ['Consumable'],
        'Monturas': ['Mount'],
        'Equipamiento': ['Weapon', 'Armor', 'Accessory', 'Offhand', 'Tome'],
        'Cosméticos': ['Skin', 'Luxury']
    };

    useEffect(() => {
        const fetchRoutes = async () => {
            try {
                const res = await fetch(`http://localhost:8081/calculate/routes?region=${region}`);
                if (res.ok) {
                    const data = await res.json();
                    setRoutes(data);
                }
            } catch (err) {
            } finally {
                setLoading(false);
            }
        };
        fetchRoutes();
        const interval = setInterval(fetchRoutes, 15000);
        return () => clearInterval(interval);
    }, [region]);

    const filteredRoutes = (routes || [])
        .filter(r => {
            if (!r) return false;
            const matchCat = selectedCategory === 'Todos' || (r.category && categoryMap[selectedCategory].includes(r.category));
            const matchSearch = (r.item_name || '').toLowerCase().includes(searchTerm.toLowerCase());
            const matchProfit = (r.profit || 0) >= minProfit;

            // Detect Tier from ID
            const id = r.item_id || '';
            const tierMatch = id.match(/T(\d)/);
            const tier = tierMatch ? parseInt(tierMatch[1]) : 0;
            const matchTier = selectedTiers.includes(tier);

            // Sub-category check for weapons
            let matchSubCat = true;
            if (selectedCategory === 'Equipamiento' && r.category === 'Weapon' && selectedWeaponSubCat !== 'Todas') {
                matchSubCat = r.sub_category === selectedWeaponSubCat;
            }

            return matchCat && matchSearch && matchProfit && matchTier && matchSubCat;
        })
        .sort((a, b) => b.profit - a.profit);

    const formatPrice = (p: number) => {
        if (p >= 1000000) return (p / 1000000).toFixed(1) + 'M';
        if (p >= 1000) return (p / 1000).toFixed(1) + 'k';
        return p.toLocaleString();
    };

    return (
        <div className="card h-full overflow-hidden flex flex-col gap-4">
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.6rem', padding: '0.6rem', background: 'rgba(255,255,255,0.02)', borderRadius: '0.6rem', border: '1px solid var(--glass-border)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h3 style={{ fontSize: '0.9rem', fontWeight: 700, display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
                        <ArrowRightLeft className="text-primary" size={16} />
                        Arbitraje Global Pro
                    </h3>
                    <div style={{ display: 'flex', gap: '0.3rem' }}>
                        {categories.map(cat => (
                            <button
                                key={cat}
                                onClick={() => setSelectedCategory(cat)}
                                style={{
                                    padding: '0.15rem 0.4rem',
                                    borderRadius: '0.3rem',
                                    fontSize: '0.6rem',
                                    border: '1px solid var(--glass-border)',
                                    background: selectedCategory === cat ? 'var(--color-accent)' : 'rgba(255,255,255,0.03)',
                                    color: 'white',
                                    cursor: 'pointer'
                                }}
                            >
                                {cat}
                            </button>
                        ))}
                    </div>
                </div>

                {selectedCategory === 'Equipamiento' && (
                    <div style={{ display: 'flex', gap: '0.3rem', flexWrap: 'wrap', paddingBottom: '0.3rem', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                        <span style={{ fontSize: '0.55rem', color: 'var(--color-text-secondary)', alignSelf: 'center', marginRight: '0.4rem' }}>ARMAS:</span>
                        {weaponSubCats.map(sub => (
                            <button
                                key={sub}
                                onClick={() => setSelectedWeaponSubCat(sub)}
                                style={{
                                    padding: '0.1rem 0.35rem',
                                    borderRadius: '0.25rem',
                                    fontSize: '0.55rem',
                                    border: '1px solid var(--glass-border)',
                                    background: selectedWeaponSubCat === sub ? 'var(--color-accent)' : 'transparent',
                                    color: 'white',
                                    cursor: 'pointer'
                                }}
                            >
                                {sub}
                            </button>
                        ))}
                    </div>
                )}

                <div style={{ display: 'flex', gap: '0.75rem', alignItems: 'center', flexWrap: 'wrap' }}>
                    <div style={{ display: 'flex', gap: '0.2rem', alignItems: 'center' }}>
                        <span style={{ fontSize: '0.6rem', color: 'var(--color-text-secondary)' }}>TIER:</span>
                        {[4, 5, 6, 7, 8].map(t => (
                            <button key={t} onClick={() => toggleTier(t)} style={{ width: '20px', height: '20px', fontSize: '0.55rem', borderRadius: '3px', border: '1px solid var(--glass-border)', background: selectedTiers.includes(t) ? 'var(--color-accent)' : 'transparent', color: 'white', cursor: 'pointer' }}>T{t}</button>
                        ))}
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.3rem' }}>
                        <span style={{ fontSize: '0.6rem', color: 'var(--color-text-secondary)' }}>GAN.:</span>
                        <input type="number" value={minProfit} onChange={(e) => setMinProfit(Number(e.target.value))} style={{ width: '60px', background: 'rgba(0,0,0,0.2)', border: '1px solid var(--glass-border)', borderRadius: '3px', color: 'white', fontSize: '0.65rem', padding: '1px 3px' }} />
                    </div>
                    <input
                        type="text"
                        placeholder="Filtrar objeto..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        style={{ flex: 1, minWidth: '100px', background: 'rgba(0,0,0,0.2)', border: '1px solid var(--glass-border)', borderRadius: '3px', color: 'white', fontSize: '0.65rem', padding: '2px 6px' }}
                    />
                </div>
            </div>

            <div style={{ flex: 1, overflowY: 'auto' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.8rem' }}>
                    <thead style={{ position: 'sticky', top: 0, background: 'var(--color-bg)', zIndex: 10 }}>
                        <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--glass-border)', color: 'var(--color-text-secondary)' }}>
                            <th style={{ padding: '0.5rem' }}>OBJETO</th>
                            <th style={{ padding: '0.5rem' }}>RUTA</th>
                            <th style={{ padding: '0.5rem' }}>TIEMPO</th>
                            <th style={{ padding: '0.5rem' }}>EJECUCIÓN</th>
                            <th style={{ padding: '0.5rem' }}>COMPRA</th>
                            <th style={{ padding: '0.5rem' }}>VENTA</th>
                            <th style={{ padding: '0.5rem' }}>GANANCIA</th>
                            <th style={{ padding: '0.5rem' }}>ROI</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-glass-border">
                        {(filteredRoutes || []).map((route, i) => {
                            // Parse Tier and Enchantment
                            const tierMatch = route.item_id.match(/T(\d)/);
                            const tier = tierMatch ? tierMatch[1] : '?';
                            const enchantMatch = route.item_id.match(/@(\d)/);
                            const enchant = enchantMatch ? `.${enchantMatch[1]}` : '.0';

                            // Determine badge color for enchantment
                            let enchantColor = 'rgba(255,255,255,0.2)';
                            if (enchant === '.1') enchantColor = '#50C878'; // Green
                            if (enchant === '.2') enchantColor = '#40E0D0'; // Turquoise/Blue
                            if (enchant === '.3') enchantColor = '#9370DB'; // Purple
                            if (enchant === '.4') enchantColor = '#FFD700'; // Gold

                            const qualityName = getQualityName(route.quality);
                            const qualityColor = getQualityColor(route.quality);
                            const timeSince = calculateTimeSince(route.updated_at);

                            // Hot logic: < 30 mins and profitable
                            const isHot = timeSince.includes('m') && !timeSince.includes('h') && !timeSince.includes('d') && parseInt(timeSince) < 30;

                            return (
                                <tr key={i} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                                    <td style={{ padding: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
                                        <div style={{ position: 'relative' }}>
                                            <img src={`https://render.albiononline.com/v1/item/${route.item_id}.png?size=48`} alt="" style={{ width: '32px', height: '32px' }} />
                                            <div style={{
                                                position: 'absolute',
                                                bottom: -2,
                                                right: -4,
                                                background: 'rgba(0,0,0,0.8)',
                                                color: 'white',
                                                fontSize: '0.55rem',
                                                padding: '0 2px',
                                                borderRadius: '2px',
                                                border: '1px solid rgba(255,255,255,0.2)'
                                            }}>T{tier}</div>
                                        </div>
                                        <div>
                                            <div style={{ fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.4rem' }}>
                                                {route.item_name}
                                                <span style={{
                                                    fontSize: '0.6rem',
                                                    padding: '1px 4px',
                                                    borderRadius: '3px',
                                                    background: enchantColor,
                                                    color: enchantColor === '#FFD700' ? 'black' : 'white',
                                                    fontWeight: 700,
                                                    boxShadow: '0 0 5px ' + enchantColor
                                                }}>{enchant}</span>
                                                {isHot && <Flame size={12} className="text-orange-500 animate-pulse" />}
                                            </div>
                                            <div style={{ display: 'flex', gap: '0.3rem', alignItems: 'center', marginTop: '0.1rem' }}>
                                                <span style={{ fontSize: '0.55rem', opacity: 0.5 }}>{route.sub_category}</span>
                                                <span style={{
                                                    fontSize: '0.55rem',
                                                    color: qualityColor,
                                                    border: `1px solid ${qualityColor}`,
                                                    padding: '0 3px',
                                                    borderRadius: '2px',
                                                    fontWeight: 500
                                                }}>{qualityName}</span>
                                            </div>
                                        </div>
                                    </td>
                                    <td style={{ padding: '0.5rem' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.3rem', fontSize: '0.7rem' }}>
                                            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                                                <span style={{ color: 'var(--color-primary)' }}>{route.from_city}</span>
                                            </div>
                                            <ArrowRightLeft size={10} />
                                            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                                                <span style={{ color: 'var(--color-accent)' }}>{route.to_city}</span>
                                            </div>
                                        </div>
                                    </td>
                                    <td style={{ padding: '0.5rem', fontSize: '0.7rem', color: isHot ? 'var(--color-accent)' : 'var(--color-text-secondary)' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.2rem' }}>
                                            <Clock size={10} />
                                            {timeSince}
                                        </div>
                                    </td>
                                    <td style={{ padding: '0.5rem' }}>
                                        <div
                                            title={`Calidad de Datos: ${route.data_quality}%\nPerfiles: ${route.risk_profile}`}
                                            style={{
                                                display: 'flex',
                                                flexDirection: 'column',
                                                gap: '0.2rem',
                                                fontSize: '0.65rem',
                                                fontWeight: 600,
                                            }}
                                        >
                                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.3rem', color: getConfidenceColor(route.execution_score) }}>
                                                <Shield size={12} />
                                                {route.execution_score}%
                                                {route.is_filtered && <AlertTriangle size={10} className="text-amber-500" />}
                                            </div>
                                            <div style={{
                                                fontSize: '0.55rem',
                                                padding: '1px 4px',
                                                borderRadius: '3px',
                                                background: 'rgba(255,255,255,0.05)',
                                                border: `1px solid ${getRiskColor(route.risk_profile)}`,
                                                color: getRiskColor(route.risk_profile),
                                                width: 'fit-content'
                                            }}>
                                                Riesgo: {route.risk_profile}
                                            </div>
                                        </div>
                                    </td>
                                    <td style={{ padding: '0.5rem', color: 'var(--color-text-secondary)' }}>{formatPrice(route.buy_price)}</td>
                                    <td style={{ padding: '0.5rem', color: 'var(--color-text-primary)' }}>{formatPrice(route.sell_price)}</td>
                                    <td style={{ padding: '0.5rem', color: 'var(--color-success)', fontWeight: 700 }}>+{formatPrice(route.profit)}</td>
                                    <td style={{ padding: '0.5rem', fontWeight: 600 }}>{route.roi.toFixed(1)}%</td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
                {filteredRoutes.length === 0 && !loading && (
                    <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--color-text-secondary)' }}>
                        No se detectaron rutas de arbitraje o ajusta los filtros.
                    </div>
                )}
            </div>
        </div>
    );
};

export default TradeRoutesTable;
