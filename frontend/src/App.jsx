import { useState, useEffect, useCallback } from "react";
import RequestList from "./components/RequestList";
import RequestDetail from "./components/RequestDetail";
import FilterBar from "./components/FilterBar";
import "./App.css";

const API = "http://localhost:3000";

function matchesFilters(req, filters) {
  if (filters.method && req.method !== filters.method) return false;
  if (filters.status && !String(req.res_status).startsWith(filters.status)) return false;
  if (filters.search && !req.url.toLowerCase().includes(filters.search.toLowerCase())) return false;
  return true;
}

function App() {
  const [requests, setRequests] = useState([]);
  const [selected, setSelected] = useState(null);
  const [filters, setFilters] = useState({ method: "", status: "", search: "" });

  // Build query string from active filters
  const buildQuery = useCallback((f) => {
    const params = new URLSearchParams();
    if (f.method) params.set("method", f.method);
    if (f.status) params.set("status", f.status);
    if (f.search) params.set("search", f.search);
    const qs = params.toString();
    return qs ? `?${qs}` : "";
  }, []);

  // Re-fetch when filters change
  useEffect(() => {
    fetch(`${API}/api/requests${buildQuery(filters)}`)
      .then((res) => res.json())
      .then((json) => setRequests(json || []));
  }, [filters, buildQuery]);

  // WebSocket live feed — drop arrivals that don't match active filters
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:3000/ws`);
    ws.onmessage = (event) => {
      const newItem = JSON.parse(event.data);
      setRequests((prev) => {
        if (!matchesFilters(newItem, filters)) return prev;
        return [newItem, ...prev];
      });
    };
    return () => ws.close();
  }, [filters]);

  const handleClear = () => {
    fetch(`${API}/api/requests`, { method: "DELETE" }).then(() => setRequests([]));
  };

  return (
    <div className="app-wrapper">
      <header className="topbar">
        <span className="topbar-logo">Flowgate</span>
        <span className="topbar-sub">HTTP/HTTPS Intercepting Proxy</span>
      </header>
      <div className="app">
        <div className="left-panel">
          <FilterBar filters={filters} onChange={setFilters} />
          <RequestList requests={requests} selected={selected} onSelect={setSelected} onClear={handleClear} />
        </div>
        <RequestDetail request={selected} />
      </div>
    </div>
  );
}

export default App;
