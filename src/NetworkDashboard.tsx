import React, { useState, useEffect } from 'react';
import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer
} from 'recharts';
import axios from 'axios';

interface NetworkData {
  timestamp: string;
  throughput: number;
  latency: number;
}

interface FetchError {
  message: string;
}

const NetworkTrafficVisualizer: React.FC = () => {
  const [networkData, setNetworkData] = useState<NetworkData[]>([]);
  const [fetchError, setFetchError] = useState<FetchError | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${process.env.REACT_APP_API_BASE_URL}/network-traffic`);
        setNetworkData(response.data);
        // Reset error state in case of successful fetch
        setFetchError(null);
      } catch (error) {
        console.error('Error fetching network data:', error);
        // Capture axios error message or use a fallback
        const errorMessage = axios.isAxiosError(error) && error.message ? error.message : 'Failed to fetch network data';
        setFetchError({ message: errorMessage });
      }
    };
    fetchData();

    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  if (fetchError) {
    return <div>Error: {fetchError.message}. Please try again later.</div>;
  }

  return (
    <div style={{ width: '100%', height: 300 }}>
      <ResponsiveContainer>
        <LineChart data={networkData}>
          <Line type="monotone" dataKey="throughput" stroke="#8884d8" />
          <CartesianGrid stroke="#ccc" />
          <XAxis dataKey="timestamp" />
          <YAxis />
          <Tooltip />
        </LineChart>
      </ResponsiveContainer>

      <ResponsiveContainer>
        <LineChart data={networkData}>
          <Line type="monotone" dataKey="latency" stroke="#82ca9d" />
          <CartesianGrid stroke="#ccc" />
          <XAxis dataKey="timestamp" />
          <YAxis />
          <Tooltip />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default NetworkTrafficVisualizer;