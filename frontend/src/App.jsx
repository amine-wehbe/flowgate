import { useState, useEffect } from "react";
import RequestList from "./components/RequestList";
import RequestDetail from "./components/RequestDetail";
import "./App.css";

function App() {
  const [requests, setRequests] = useState([]);
  const [selected, setSelected] = useState(null);
  const handleClear = () => {
  fetch("http://localhost:3000/api/requests", { method: "DELETE" })
    .then(() => setRequests([]));
};

  // Load existing requests from DB on mount
  useEffect(() => {
    fetch("http://localhost:3000/api/requests")
      .then((res) => res.json())
      .then((json) => setRequests(json || []));
  }, []);

  // Open WebSocket once — prepend each new request to the top of the list
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:3000/ws");
    ws.onmessage = (event) => {
      const newItem = JSON.parse(event.data);
      setRequests((prev) => [newItem, ...prev]);
    };
    return () => ws.close();
  }, []);

  return (
    <div className="app-wrapper">
      <header className="topbar">
        <span className="topbar-logo">Flowgate</span>
        <span className="topbar-sub">HTTP/HTTPS Intercepting Proxy</span>
      </header>
      <div className="app">
        <RequestList requests={requests} selected={selected} onSelect={setSelected} onClear={handleClear}/>
        <RequestDetail request={selected} />
      </div>
    </div>
  );
}

export default App;
