import React, { useState, useEffect } from 'react';

interface MarketOrder {
    OrderID: number;
    ItemID: string;
    LocationID: number;
    Quality: number;
    UnitPrice: number;
    Amount: number;
    AuctionType: string;
    Expires: string;
    UpdatedAt: string;
    Source: string;
    Confidence: string;
    Region: string;
}

interface Opportunity {
    OriginalOrder: MarketOrder;
    Profit: number;
    ROI: number;
    Confidence: string; // Legacy
    DataQuality: number;
    ExecutionScore: number;
    RiskProfile: string;
    Category: string;
    SubCategory: string;
    ItemName: string;
    BuyLocation: string;
    SellLocation: string;
    SellPrice?: number;
}

interface OpportunityTableProps {
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

const OpportunityTable: React.FC<OpportunityTableProps> = ({ region }) => {
    const [opportunities, setOpportunities] = useState<Opportunity[]>([]);
    const [selectedCategory, setSelectedCategory] = useState<string>('Todos');
    const [isPremium, setIsPremium] = useState<boolean>(true);
    const [sortKey, setSortKey] = useState<'Profit' | 'ROI'>('ROI');
    const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

    // Advanced Filters
    const [minProfit, setMinProfit] = useState<number>(0);
    const [selectedTiers, setSelectedTiers] = useState<number[]>([4, 5, 6, 7, 8]);
    const [selectedEnchants, setSelectedEnchants] = useState<number[]>([0, 1, 2, 3, 4]);
    const [searchTerm, setSearchTerm] = useState<string>('');
    const [selectedWeaponSubCat, setSelectedWeaponSubCat] = useState<string>('Todas');

    const weaponSubCats = ['Todas', 'Espadas', 'Hachas', 'Mazas', 'Martillos', 'Ballestas', 'Arcos', 'Lanzas', 'Dagas', 'Bastón de Fuego', 'Bastón de Hielo', 'Bastón Maldito', 'Bastón Sagrado', 'Bastón Natural'];

    const toggleTier = (t: number) => {
        setSelectedTiers(prev => prev.includes(t) ? prev.filter(x => x !== t) : [...prev, t]);
    };
    const toggleEnchant = (e: number) => {
        setSelectedEnchants(prev => prev.includes(e) ? prev.filter(x => x !== e) : [...prev, e]);
    };

    // Only BM-relevant categories
    const categories = ['Todos', 'Armas', 'Armaduras', 'Accesorios'];

    const categoryMap: Record<string, string[]> = {
        'Todos': ['Weapon', 'Armor', 'Accessory', 'Offhand', 'Tome'],
        'Armas': ['Weapon'],
        'Armaduras': ['Armor'],
        'Accesorios': ['Accessory', 'Offhand']
    };

    const filteredOpportunities = (opportunities || [])
        .filter(o => {
            if (!o || !o.OriginalOrder) return false;

            const matchRegion = o.OriginalOrder.Region === region;
            const matchCat = selectedCategory === 'Todos' || (o.Category && categoryMap[selectedCategory].includes(o.Category));
            const matchProfit = o.Profit >= minProfit;
            const matchSearch = (o.ItemName || '').toLowerCase().includes(searchTerm.toLowerCase());

            // Tier/Enchant detection from ID (e.g. T4_...LEVEL1@1)
            const id = o.OriginalOrder.ItemID || '';
            const tierMatch = id.match(/T(\d)/);
            const tier = tierMatch ? parseInt(tierMatch[1]) : 0;
            const enchant = id.includes('@') ? parseInt(id.split('@')[1]) : 0;

            const matchTier = selectedTiers.includes(tier);
            const matchEnchant = selectedEnchants.includes(enchant);

            // Sub-category check for weapons
            let matchSubCat = true;
            if (selectedCategory === 'Armas' && selectedWeaponSubCat !== 'Todas') {
                matchSubCat = o.SubCategory === selectedWeaponSubCat;
            }

            return matchRegion && matchCat && matchProfit && matchSearch && matchTier && matchEnchant && matchSubCat;
        })
        .sort((a, b) => {
            const factor = sortOrder === 'asc' ? 1 : -1;
            return (a[sortKey] - b[sortKey]) * factor;
        });

    const toggleSort = (key: 'Profit' | 'ROI') => {
        if (sortKey === key) {
            setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
        } else {
            setSortKey(key);
            setSortOrder('desc');
        }
    };

    const togglePremium = async () => {
        const next = !isPremium;
        setIsPremium(next);
        try {
            await fetch('http://localhost:8081/calculate/premium', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ premium: next })
            });
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetch('http://localhost:8081/calculate/premium').then(res => res.json()).then(data => setIsPremium(data.premium)).catch(() => { });
        const ws = new WebSocket('ws://localhost:8081/ws');
        ws.onmessage = (event) => {
            try {
                const opp: Opportunity = JSON.parse(event.data);
                setOpportunities((prev) => {
                    const existing = (prev || []).find(o => o.OriginalOrder.OrderID === opp.OriginalOrder.OrderID);
                    if (existing) return prev;
                    return [opp, ...(prev || [])].slice(0, 50);
                });
            } catch (e) { }
        };
        return () => ws.close();
    }, []);

    return (
        <div className="card" style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem', marginBottom: '1rem', padding: '0.75rem', background: 'rgba(255,255,255,0.02)', borderRadius: '0.75rem', border: '1px solid var(--glass-border)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <div style={{ display: 'flex', gap: '0.5rem' }}>
                        {categories.map(cat => (
                            <button
                                key={cat}
                                onClick={() => setSelectedCategory(cat)}
                                style={{
                                    padding: '0.3rem 0.6rem',
                                    borderRadius: '0.5rem',
                                    fontSize: '0.7rem',
                                    fontWeight: 600,
                                    border: '1px solid var(--glass-border)',
                                    background: selectedCategory === cat ? 'var(--color-primary)' : 'rgba(255,255,255,0.03)',
                                    color: 'white',
                                    cursor: 'pointer'
                                }}
                            >
                                {cat}
                            </button>
                        ))}
                    </div>

                    <div
                        onClick={togglePremium}
                        style={{ fontSize: '0.7rem', color: isPremium ? 'var(--color-success)' : 'var(--color-text-secondary)', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '0.4rem' }}
                    >
                        <div style={{ width: '8px', height: '8px', borderRadius: '50%', background: isPremium ? 'var(--color-success)' : '#444' }} />
                        {isPremium ? 'PREMIUM' : 'ESTÁNDAR'}
                    </div>
                </div>

                {selectedCategory === 'Armas' && (
                    <div style={{ display: 'flex', gap: '0.3rem', flexWrap: 'wrap', paddingBottom: '0.5rem', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                        <span style={{ fontSize: '0.6rem', color: 'var(--color-text-secondary)', alignSelf: 'center', marginRight: '0.4rem' }}>TIPO ARMA:</span>
                        {weaponSubCats.map(sub => (
                            <button
                                key={sub}
                                onClick={() => setSelectedWeaponSubCat(sub)}
                                style={{
                                    padding: '0.15rem 0.4rem',
                                    borderRadius: '0.3rem',
                                    fontSize: '0.6rem',
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

                <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', alignItems: 'center' }}>
                    <div style={{ display: 'flex', gap: '0.25rem', alignItems: 'center' }}>
                        <span style={{ fontSize: '0.65rem', color: 'var(--color-text-secondary)', marginRight: '0.25rem' }}>TIER:</span>
                        {[4, 5, 6, 7, 8].map(t => (
                            <button key={t} onClick={() => toggleTier(t)} style={{ width: '22px', height: '22px', fontSize: '0.6rem', borderRadius: '4px', border: '1px solid var(--glass-border)', background: selectedTiers.includes(t) ? 'var(--color-accent)' : 'transparent', color: 'white', cursor: 'pointer' }}>T{t}</button>
                        ))}
                    </div>
                    <div style={{ display: 'flex', gap: '0.25rem', alignItems: 'center' }}>
                        <span style={{ fontSize: '0.65rem', color: 'var(--color-text-secondary)', marginRight: '0.25rem' }}>ENC:</span>
                        {[0, 1, 2, 3, 4].map(e => (
                            <button key={e} onClick={() => toggleEnchant(e)} style={{ width: '22px', height: '22px', fontSize: '0.6rem', borderRadius: '4px', border: '1px solid var(--glass-border)', background: selectedEnchants.includes(e) ? 'var(--color-success)' : 'transparent', color: 'white', cursor: 'pointer' }}>.{e}</button>
                        ))}
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <span style={{ fontSize: '0.65rem', color: 'var(--color-text-secondary)' }}>GAN. MÍN:</span>
                        <input type="number" value={minProfit} onChange={(e) => setMinProfit(Number(e.target.value))} style={{ width: '70px', background: 'rgba(0,0,0,0.2)', border: '1px solid var(--glass-border)', borderRadius: '4px', color: 'white', fontSize: '0.7rem', padding: '2px 4px' }} />
                    </div>
                    <input
                        type="text"
                        placeholder="Buscar por nombre..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        style={{ flex: 1, minWidth: '120px', background: 'rgba(0,0,0,0.2)', border: '1px solid var(--glass-border)', borderRadius: '4px', color: 'white', fontSize: '0.7rem', padding: '4px 8px' }}
                    />
                </div>
            </div>

            <div style={{ flex: 1, overflowY: 'auto' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.85rem' }}>
                    <thead style={{ position: 'sticky', top: 0, background: 'var(--color-bg)', zIndex: 10 }}>
                        <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--glass-border)', color: 'var(--color-text-secondary)' }}>
                            <th style={{ padding: '0.5rem' }}>OBJETO</th>
                            <th style={{ padding: '0.5rem' }}>CALIDAD</th>
                            <th style={{ padding: '0.5rem' }}>ORIGEN</th>
                            <th style={{ padding: '0.5rem' }}>COMPRA</th>
                            <th style={{ padding: '0.5rem' }}>VENTA EST.</th>
                            <th style={{ padding: '0.5rem' }}>EJECUCIÓN</th>
                            <th style={{ padding: '0.5rem', cursor: 'pointer' }} onClick={() => toggleSort('Profit')}>GANANCIA</th>
                            <th style={{ padding: '0.5rem', cursor: 'pointer' }} onClick={() => toggleSort('ROI')}>ROI</th>
                        </tr>
                    </thead>
                    <tbody>
                        {(filteredOpportunities || []).map((opp) => {
                            const qualityName = getQualityName(opp.OriginalOrder.Quality);
                            const qualityColor = getQualityColor(opp.OriginalOrder.Quality);

                            return (
                                <tr key={opp.OriginalOrder.OrderID} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                                    <td style={{ padding: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                        <img src={`https://render.albiononline.com/v1/item/${opp.OriginalOrder.ItemID}.png?size=32`} alt="" style={{ width: '24px', height: '24px' }} />
                                        <div>
                                            <div style={{ fontWeight: 600 }}>{opp.ItemName}</div>
                                            <div style={{ fontSize: '0.65rem', opacity: 0.5 }}>{opp.Category} - {opp.SubCategory}</div>
                                        </div>
                                    </td>
                                    <td style={{ padding: '0.5rem' }}>
                                        <span style={{
                                            fontSize: '0.6rem',
                                            color: qualityColor,
                                            border: `1px solid ${qualityColor}`,
                                            padding: '1px 4px',
                                            borderRadius: '3px',
                                            fontWeight: 600
                                        }}>{qualityName}</span>
                                    </td>
                                    <td style={{ padding: '0.5rem', fontSize: '0.75rem' }}>{opp.BuyLocation}</td>
                                    <td style={{ padding: '0.5rem', color: 'var(--color-text-secondary)' }}>{opp.OriginalOrder.UnitPrice.toLocaleString()}</td>
                                    <td style={{ padding: '0.5rem', color: 'var(--color-text-primary)' }}>{(opp.SellPrice || 0).toLocaleString()}</td>
                                    <td style={{ padding: '0.5rem' }}>
                                        <div
                                            title={`Calidad: ${opp.DataQuality}%`}
                                            style={{ display: 'flex', flexDirection: 'column', gap: '0.1rem' }}
                                        >
                                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.3rem', fontSize: '0.75rem', fontWeight: 700, color: getConfidenceColor(opp.ExecutionScore) }}>
                                                {opp.ExecutionScore}%
                                            </div>
                                            <div style={{ fontSize: '0.55rem', color: getRiskColor(opp.RiskProfile), fontWeight: 600 }}>
                                                {opp.RiskProfile}
                                            </div>
                                        </div>
                                    </td>
                                    <td style={{ padding: '0.5rem', color: 'var(--color-success)', fontWeight: 700 }}>+{Math.floor(opp.Profit).toLocaleString()}</td>
                                    <td style={{ padding: '0.5rem', fontWeight: 600 }}>{opp.ROI.toFixed(1)}%</td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
                {filteredOpportunities.length === 0 && (
                    <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--color-text-secondary)', fontSize: '0.8rem' }}>
                        Buscando oportunidades en Mercado Negro...
                    </div>
                )}
            </div>
        </div>
    );
};

export default OpportunityTable;
