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

const NetworkTrafficVisualizer: React.FC = () => {
  const [networkData, setNetworkData] = useState<NetworkData[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${process.env.REACT_APP_API_BASE_URL}/network-traffic`);
        setNetworkData(response.data);
      } catch (error) {
        console.error('Error fetching network data:', error);
      }
    };
    fetchData();

    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

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