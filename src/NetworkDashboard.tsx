import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

interface NetworkData {
  timestamp: string;
  throughput: number;
  latency: number;
}

interface FetchError {
  message: string;
}

const NetworkTrafficVisualizer: React.FC = () => {
  const [networkData, setNetworkData] = useState<NetworkWealth[]>([]);
  const [fetchError, setFetchError] = useState<FetchError | null>(null);

  // Data fetching function
  const fetchData = async () => {
    try {
      const response = await axios.get(`${process.env.REACT_APP_API_BASE_URL}/network-traffic`);
      setNetworkData(response.data);
      setFetchError(null);
    } catch (error) {
      console.error('Error fetching network data:', error);
      const errorMessage = axios.isAxiosError(error) && error.message ? error.message : 'Failed to fetch network data';
      setFetchError({ message: errorMessage });
    }
  };

  useEffect(() => {
    let interval: NodeJS.Timer;
    const startFetching = () => {
      fetchData();
      interval = setInterval(fetchData, 5000);
    };
    
    const stopFetching = () => {
      clearInterval(interval);
    };

    // Start fetching when the component mounts or becomes visible
    if (document.visibilityState === "visible") {
      startFetching();
    } else {
      stopFetching();
    }

    document.addEventListener("visibilitychange", () => {
      if (document.visibilityState === "visible") {
        startFetching();
      } else {
        stopFetching();
      }
    });

    return () => {
      stopFetching();
      document.removeEventListener("visibilitychange", startFetching);
      document.removeEventListener("visibilitychange", stopFetching);
    };
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

export default NetworkTraffic Visualizer;