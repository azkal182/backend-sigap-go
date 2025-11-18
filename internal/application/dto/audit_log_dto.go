package dto

// AuditLogResponse represents audit log data in responses
type AuditLogResponse struct {
	ID            string   `json:"id"`
	ActorID       string   `json:"actor_id,omitempty"`
	ActorUsername string   `json:"username,omitempty"`
	ActorRoles    []string `json:"actor_roles,omitempty"`
	Action        string   `json:"action"`
	Resource      string   `json:"resource"`
	TargetID      string   `json:"target_id,omitempty"`
	RequestPath   string   `json:"request_path"`
	RequestMethod string   `json:"request_method"`
	StatusCode    int      `json:"status_code"`
	IPAddress     string   `json:"ip_address,omitempty"`
	UserAgent     string   `json:"user_agent,omitempty"`
	Metadata      string   `json:"metadata,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

// ListAuditLogsResponse represents paginated audit log list response
type ListAuditLogsResponse struct {
	Logs       []AuditLogResponse `json:"logs"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
