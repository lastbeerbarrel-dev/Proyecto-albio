import React from 'react';
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer
} from 'recharts';

interface PriceHistoryProps {
    itemId: string;
}

const PriceHistory: React.FC<PriceHistoryProps> = ({ itemId }) => {
    // Mock Data for MVP - In real app, fetch from Metadata Service
    const data = [
        { time: '10:00', price: 4200 },
        { time: '10:05', price: 4150 },
        { time: '10:10', price: 4300 },
        { time: '10:15', price: 4250 },
        { time: '10:20', price: 4100 },
        { time: '10:25', price: 4400 },
    ];

    return (
        <div style={{ padding: '1rem', background: 'rgba(0,0,0,0.2)', borderRadius: '0.5rem', marginTop: '0.5rem' }}>
            <h4 style={{ margin: '0 0 1rem 0', fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                Price Trend (Last Hour) - {itemId}
            </h4>
            <div style={{ width: '100%', height: 200 }}>
                <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={data}>
                        <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                        <XAxis dataKey="time" stroke="rgba(255,255,255,0.5)" fontSize={12} />
                        <YAxis stroke="rgba(255,255,255,0.5)" fontSize={12} domain={['auto', 'auto']} />
                        <Tooltip
                            contentStyle={{ background: 'var(--glass-bg)', border: '1px solid var(--glass-border)', color: 'white' }}
                        />
                        <Line type="monotone" dataKey="price" stroke="var(--color-accent)" strokeWidth={2} dot={{ r: 4 }} activeDot={{ r: 6 }} />
                    </LineChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
};

export default PriceHistory;
