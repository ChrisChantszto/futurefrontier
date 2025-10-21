import { useState, useEffect } from 'react'
import { Search, RefreshCw, Download } from 'lucide-react'
import axios from 'axios'
import { format } from 'date-fns'

interface Log {
  timestamp: string
  method: string
  path: string
  status_code: number
  latency_ms: number
  client_ip: string
  error_message?: string
}

export default function Discover() {
  const [logs, setLogs] = useState<Log[]>([])
  const [loading, setLoading] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')
  const [timeRange, setTimeRange] = useState('1h')

  const fetchLogs = async () => {
    setLoading(true)
    try {
      const response = await axios.get('http://localhost:8080/api/logs', {
        params: {
          size: 100,
          search: searchQuery,
          status: statusFilter
        },
        withCredentials: true
      })

      if (response.data.success) {
        setLogs(response.data.data.logs || [])
      }
    } catch (error) {
      console.error('Failed to fetch logs:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchLogs()
  }, [searchQuery, statusFilter])

  const filteredLogs = logs // Backend handles filtering now

  const getStatusColor = (status: number) => {
    if (status < 300) return '#48bb78'
    if (status < 400) return '#ed8936'
    return '#f56565'
  }

  const exportLogs = () => {
    const csv = [
      ['Timestamp', 'Method', 'Path', 'Status', 'Latency (ms)', 'IP', 'Error'].join(','),
      ...filteredLogs.map(log => [
        log.timestamp,
        log.method,
        log.path,
        log.status_code,
        log.latency_ms,
        log.client_ip,
        log.error_message || ''
      ].join(','))
    ].join('\n')

    const blob = new Blob([csv], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `logs-${Date.now()}.csv`
    a.click()
  }

  return (
    <div style={{ padding: '2rem', maxWidth: '1400px', margin: '0 auto' }}>
      <div style={{ marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: 'bold', color: '#1a202c', marginBottom: '0.5rem' }}>
          Discover Logs
        </h1>
        <p style={{ color: '#718096' }}>
          Search and filter API logs from Elasticsearch
        </p>
      </div>

      {/* Filters */}
      <div style={{
        background: 'white',
        padding: '1.5rem',
        borderRadius: '12px',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
        marginBottom: '1.5rem'
      }}>
        <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', alignItems: 'center' }}>
          {/* Search */}
          <div style={{ flex: '1', minWidth: '300px', position: 'relative' }}>
            <Search size={20} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#a0aec0' }} />
            <input
              type="text"
              placeholder="Search logs (path, error message)..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              style={{
                width: '100%',
                padding: '0.75rem 0.75rem 0.75rem 2.5rem',
                border: '1px solid #e2e8f0',
                borderRadius: '8px',
                fontSize: '0.875rem',
                boxSizing: 'border-box'
              }}
            />
          </div>

          {/* Status Filter */}
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            style={{
              padding: '0.75rem',
              border: '1px solid #e2e8f0',
              borderRadius: '8px',
              fontSize: '0.875rem',
              background: 'white',
              cursor: 'pointer'
            }}
          >
            <option value="all">All Status</option>
            <option value="success">Success (2xx-3xx)</option>
            <option value="error">Errors (4xx-5xx)</option>
          </select>

          {/* Time Range */}
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
            style={{
              padding: '0.75rem',
              border: '1px solid #e2e8f0',
              borderRadius: '8px',
              fontSize: '0.875rem',
              background: 'white',
              cursor: 'pointer'
            }}
          >
            <option value="15m">Last 15 minutes</option>
            <option value="1h">Last hour</option>
            <option value="24h">Last 24 hours</option>
            <option value="7d">Last 7 days</option>
          </select>

          {/* Actions */}
          <button
            onClick={fetchLogs}
            disabled={loading}
            style={{
              padding: '0.75rem 1rem',
              background: '#667eea',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              cursor: loading ? 'not-allowed' : 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              fontSize: '0.875rem',
              fontWeight: '500'
            }}
          >
            <RefreshCw size={16} />
            Refresh
          </button>

          <button
            onClick={exportLogs}
            style={{
              padding: '0.75rem 1rem',
              background: '#48bb78',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              fontSize: '0.875rem',
              fontWeight: '500'
            }}
          >
            <Download size={16} />
            Export CSV
          </button>
        </div>

        <div style={{ marginTop: '1rem', fontSize: '0.875rem', color: '#718096' }}>
          Showing {filteredLogs.length} of {logs.length} logs
        </div>
      </div>

      {/* Logs Table */}
      <div style={{
        background: 'white',
        borderRadius: '12px',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
        overflow: 'hidden'
      }}>
        <div style={{ overflowX: 'auto' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ background: '#f7fafc', borderBottom: '2px solid #e2e8f0' }}>
                <th style={{ padding: '1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: '600', color: '#4a5568', textTransform: 'uppercase' }}>Timestamp</th>
                <th style={{ padding: '1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: '600', color: '#4a5568', textTransform: 'uppercase' }}>Method</th>
                <th style={{ padding: '1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: '600', color: '#4a5568', textTransform: 'uppercase' }}>Path</th>
                <th style={{ padding: '1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: '600', color: '#4a5568', textTransform: 'uppercase' }}>Status</th>
                <th style={{ padding: '1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: '600', color: '#4a5568', textTransform: 'uppercase' }}>Latency</th>
                <th style={{ padding: '1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: '600', color: '#4a5568', textTransform: 'uppercase' }}>IP</th>
                <th style={{ padding: '1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: '600', color: '#4a5568', textTransform: 'uppercase' }}>Error</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={7} style={{ padding: '3rem', textAlign: 'center', color: '#a0aec0' }}>
                    Loading logs...
                  </td>
                </tr>
              ) : filteredLogs.length === 0 ? (
                <tr>
                  <td colSpan={7} style={{ padding: '3rem', textAlign: 'center', color: '#a0aec0' }}>
                    No logs found
                  </td>
                </tr>
              ) : (
                filteredLogs.map((log, index) => (
                  <tr key={index} style={{ borderBottom: '1px solid #e2e8f0' }}>
                    <td style={{ padding: '1rem', fontSize: '0.875rem', color: '#4a5568', fontFamily: 'monospace' }}>
                      {log.timestamp ? format(new Date(log.timestamp), 'MMM dd, HH:mm:ss') : 'N/A'}
                    </td>
                    <td style={{ padding: '1rem' }}>
                      <span style={{
                        padding: '0.25rem 0.5rem',
                        background: '#edf2f7',
                        borderRadius: '4px',
                        fontSize: '0.75rem',
                        fontWeight: '600',
                        color: '#2d3748'
                      }}>
                        {log.method}
                      </span>
                    </td>
                    <td style={{ padding: '1rem', fontSize: '0.875rem', color: '#2d3748', fontFamily: 'monospace' }}>
                      {log.path}
                    </td>
                    <td style={{ padding: '1rem' }}>
                      <span style={{
                        padding: '0.25rem 0.5rem',
                        background: `${getStatusColor(log.status_code)}20`,
                        color: getStatusColor(log.status_code),
                        borderRadius: '4px',
                        fontSize: '0.75rem',
                        fontWeight: '600'
                      }}>
                        {log.status_code}
                      </span>
                    </td>
                    <td style={{ padding: '1rem', fontSize: '0.875rem', color: '#4a5568' }}>
                      {log.latency_ms}ms
                    </td>
                    <td style={{ padding: '1rem', fontSize: '0.875rem', color: '#4a5568', fontFamily: 'monospace' }}>
                      {log.client_ip}
                    </td>
                    <td style={{ padding: '1rem', fontSize: '0.875rem', color: '#e53e3e', maxWidth: '200px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                      {log.error_message || '-'}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
