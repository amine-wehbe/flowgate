import RequestItem from "./RequestItem";

// Left panel — renders the full list of captured requests
function RequestList({ requests, selected, onSelect, onClear }) {
  return (
    <div className="request-list">
      <div className="request-list-toolbar">
        <button onClick={onClear}>Clear</button>
      </div>
      <div className="request-list-header">
        <span>Method</span>
        <span>Host</span>
        <span>Path</span>
        <span>Status</span>
        <span>Duration</span>
      </div>
      {requests && requests.length > 0 ? requests.map((req) => (
        <RequestItem
          key={req.id || req.url + req.captured_at}
          request={req}
          isSelected={selected === req}
          onSelect={onSelect}
        />
      )) : (
        <div className="empty-list">No requests to display</div>
      )}
    </div>
  );
}

export default RequestList;
