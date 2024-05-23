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
            setLoading(true); 
            setError(null); 

            const backendUrl = process.env.REACT_APP_BACKEND_URL ?? "http://localhost:3000";
            const response = await Axios.get<NetworkStats>(`${backendUrl}/network-stats`);
            
            setStats(response.data);
        } catch (error: any) {
            const errorMessage = error.response?.data?.message || error.message || "An unknown error occurred";
            setError(errorMessage);
        } finally {
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