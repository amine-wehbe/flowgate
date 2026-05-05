import { useState } from "react";

// Right panel — shows full detail of the selected request
function RequestDetail({ request }) {
  const [replayResult, setReplayResult] = useState(null);
  const handleReplay = () => {
  fetch(`http://localhost:3000/api/requests/${request.id}/replay`, { method: "POST" })
    .then((res) => res.json())
    .then((data) => setReplayResult(data));
};
  if (!request) {
    return <div className="request-detail empty">Select a request to inspect it.</div>;
  }
  
  return (
    <div className="request-detail">
      <div className="detail-section">
        <button className="replay-btn" onClick={handleReplay}>Replay</button>
        <h3>Request</h3>
        <div className="detail-row"><span>Method</span><span>{request.method}</span></div>
        <div className="detail-row"><span>URL</span><span>{request.url}</span></div>
        <div className="detail-row"><span>Protocol</span><span>{request.protocol}</span></div>
        <div className="detail-row"><span>TLS</span><span>{request.tls ? "Yes" : "No"}</span></div>
        <div className="detail-row"><span>Captured</span><span>{new Date(request.captured_at).toLocaleString()}</span></div>
      </div>

      <div className="detail-section">
        <h3>Request Headers</h3>
        <pre>{JSON.stringify(request.req_headers, null, 2)}</pre>
      </div>

      {request.req_body && (
        <div className="detail-section">
          <h3>Request Body</h3>
          <pre>{request.req_body}</pre>
        </div>
      )}

      <div className="detail-section">
        <h3>Response</h3>
        <div className="detail-row"><span>Status</span><span>{request.res_status}</span></div>
        <div className="detail-row"><span>Duration</span><span>{request.duration_ms}ms</span></div>
      </div>

      <div className="detail-section">
        <h3>Response Headers</h3>
        <pre>{JSON.stringify(request.res_headers, null, 2)}</pre>
      </div>

      {request.res_body && (
        <div className="detail-section">
          <h3>Response Body</h3>
          <pre>{request.res_body}</pre>
        </div>
      )}
      {replayResult && (
  <div className="detail-section">
    <h3>Replay Response</h3>
    <div className="detail-row"><span>Status</span><span>{replayResult.res_status}</span></div>
    <pre>{replayResult.res_body}</pre>
  </div>
)}
    </div>
  );
}

export default RequestDetail;
