import { useState, useEffect } from 'react';
import axios from 'axios';
import { format } from 'date-fns';
import './App.css';

function App() {
  const [logs, setLogs] = useState([]);
  const [filters, setFilters] = useState({
    service_name: '',
    level: '',
    search: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const fetchLogs = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await axios.get('/api/logs', { params: filters });
      setLogs(response.data || []); 
    } catch (err) {
      setError('Failed to fetch logs.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({ ...prev, [name]: value }));
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchLogs();
  };

  return (
    <div className="container">
      <h1>Logs Dashboard</h1>

      <form className="controls" onSubmit={handleSearch}>
        <input
          name="service_name"
          value={filters.service_name}
          onChange={handleFilterChange}
          placeholder="Service Name..."
        />
        <input
          name="level"
          value={filters.level}
          onChange={handleFilterChange}
          placeholder="Level (e.g., ERROR)..."
        />
        <input
          name="search"
          value={filters.search}
          onChange={handleFilterChange}
          placeholder="Search in message..."
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Searching...' : 'Search'}
        </button>
      </form>

      <div className="logs-container">
        {error && <div className="error-message">{error}</div>}
        {!loading && logs.length === 0 && !error && <div className="no-logs">No logs found.</div>}

        {logs.map((log, index) => (
          <div key={index} className={`log-entry ${log.level?.toLowerCase()}`}>
            <span className="timestamp">{format(new Date(log.timestamp), 'HH:mm:ss.SSS')}</span>
            <span className="level">{log.level}</span>
            <span className="service">{log.service_name}</span>
            <span className="message">{log.message}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;