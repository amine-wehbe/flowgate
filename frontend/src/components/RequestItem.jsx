// Single row in the request list — highlights when selected
function RequestItem({ request, isSelected, onSelect }) {
  const statusClass =
    request.res_status >= 500 ? "status-5xx" :
    request.res_status >= 400 ? "status-4xx" :
    request.res_status >= 300 ? "status-3xx" : "status-2xx";

  return (
    <div
      className={`request-item ${isSelected ? "selected" : ""}`}
      onClick={() => onSelect(request)}
    >
      <span className={`method method-${request.method?.toLowerCase()}`}>{request.method}</span>
      <span className="host">{request.host}</span>
      <span className="path">{request.path}</span>
      <span className={`status ${statusClass}`}>{request.res_status}</span>
      <span className="duration">{request.duration_ms}ms</span>
    </div>
  );
}

export default RequestItem;
