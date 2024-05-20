import { useState, useEffect } from 'react';
import Axios from 'axios';

interface NetworkStats {
    packetsSent: number;
    packetsReceived: number;
    errorRate: number;
}

export const useNetworkStats = () => {
    const [stats, setStats] = useState<NetworkStats | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    const fetchStats = async () => {
        try {
            const response = await Axios.get<NetworkStats>(`${process.env.REACT_APP_BACKEND_URL}/network-stats`);
            setStats(response.data);
            setLoading(false);
            setError(null);
        } catch (error: any) {
            setError(error.message || "An unknown error occurred");
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchStats();
        
        const interval = setInterval(fetchStats, 5000);
        
        return () => clearInterval(interval);
    }, []);

    return { stats, loading, error };
};