import { useState, useEffect } from 'react'
import { Search, RefreshCw, Download, ChevronLeft, ChevronRight } from 'lucide-react'
import axios from 'axios'
import { format } from 'date-fns'
import { API_URL } from '../config'

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
  const [useCustomRange, setUseCustomRange] = useState(false)
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [totalLogs, setTotalLogs] = useState(0)
  const logsPerPage = 100

  const fetchLogs = async () => {
    setLoading(true)
    try {
      let params: any = {
        size: 10000, // Fetch all logs
        search: searchQuery,
        status: statusFilter,
        time_range: useCustomRange ? 'custom' : timeRange
      }

      // Add custom date range if selected
      if (useCustomRange && startDate && endDate) {
        params.start_date = startDate
        params.end_date = endDate
      }

      const response = await axios.get(`${API_URL}/api/logs`, {
        params,
        withCredentials: true
      })

      if (response.data.success) {
        const allLogs = response.data.data.logs || []
        setTotalLogs(allLogs.length)
        setLogs(allLogs)
        setCurrentPage(1) // Reset to first page on new fetch
      }
    } catch (error) {
      console.error('Failed to fetch logs:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchLogs()
  }, [searchQuery, statusFilter, timeRange, useCustomRange, startDate, endDate])

  // Pagination logic
  const totalPages = Math.ceil(totalLogs / logsPerPage)
  const startIndex = (currentPage - 1) * logsPerPage
  const endIndex = startIndex + logsPerPage
  const currentLogs = logs.slice(startIndex, endIndex)

  const goToNextPage = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }

  const goToPrevPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }

  const goToPage = (page: number) => {
    setCurrentPage(page)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const getStatusColor = (status: number) => {
    if (status < 300) return '#48bb78'
    if (status < 400) return '#ed8936'
    return '#f56565'
  }

  const exportLogs = () => {
    const csv = [
      ['Timestamp', 'Method', 'Path', 'Status', 'Latency (ms)', 'IP', 'Error'].join(','),
      ...logs.map(log => [
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
              color: '#2d3748',
              cursor: 'pointer'
            }}
          >
            <option value="all">All Status</option>
            <option value="success">Success (2xx-3xx)</option>
            <option value="error">Errors (4xx-5xx)</option>
          </select>

          {/* Time Range */}
          <select
            value={useCustomRange ? 'custom' : timeRange}
            onChange={(e) => {
              if (e.target.value === 'custom') {
                setUseCustomRange(true)
              } else {
                setUseCustomRange(false)
                setTimeRange(e.target.value)
              }
            }}
            style={{
              padding: '0.75rem',
              border: '1px solid #e2e8f0',
              borderRadius: '8px',
              fontSize: '0.875rem',
              background: 'white',
              color: '#2d3748',
              cursor: 'pointer'
            }}
          >
            <option value="15m">Last 15 minutes</option>
            <option value="1h">Last hour</option>
            <option value="6h">Last 6 hours</option>
            <option value="24h">Last 24 hours</option>
            <option value="7d">Last 7 days</option>
            <option value="30d">Last 30 days</option>
            <option value="all">All data</option>
            <option value="custom">Custom range</option>
          </select>

          {/* Custom Date Range */}
          {useCustomRange && (
            <>
              <input
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                style={{
                  padding: '0.75rem',
                  border: '1px solid #e2e8f0',
                  borderRadius: '8px',
                  fontSize: '0.875rem',
                  background: 'white',
                  color: '#2d3748',
                  cursor: 'pointer'
                }}
              />
              <span style={{ color: '#718096', fontWeight: '500' }}>to</span>
              <input
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                style={{
                  padding: '0.75rem',
                  border: '1px solid #e2e8f0',
                  borderRadius: '8px',
                  fontSize: '0.875rem',
                  background: 'white',
                  color: '#2d3748',
                  cursor: 'pointer'
                }}
              />
            </>
          )}

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
          Showing {startIndex + 1}-{Math.min(endIndex, totalLogs)} of {totalLogs} logs (Page {currentPage} of {totalPages})
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
              ) : currentLogs.length === 0 ? (
                <tr>
                  <td colSpan={7} style={{ padding: '3rem', textAlign: 'center', color: '#a0aec0' }}>
                    No logs found
                  </td>
                </tr>
              ) : (
                currentLogs.map((log, index) => (
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

      {/* Pagination */}
      {totalPages > 1 && (
        <div style={{
          marginTop: '1.5rem',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          gap: '0.5rem'
        }}>
          <button
            onClick={goToPrevPage}
            disabled={currentPage === 1}
            style={{
              padding: '0.5rem 1rem',
              background: currentPage === 1 ? '#e2e8f0' : '#667eea',
              color: currentPage === 1 ? '#a0aec0' : 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: currentPage === 1 ? 'not-allowed' : 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              fontSize: '0.875rem',
              fontWeight: '500'
            }}
          >
            <ChevronLeft size={16} />
            Previous
          </button>

          <div style={{ display: 'flex', gap: '0.25rem' }}>
            {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
              let pageNum: number
              if (totalPages <= 5) {
                pageNum = i + 1
              } else if (currentPage <= 3) {
                pageNum = i + 1
              } else if (currentPage >= totalPages - 2) {
                pageNum = totalPages - 4 + i
              } else {
                pageNum = currentPage - 2 + i
              }

              return (
                <button
                  key={pageNum}
                  onClick={() => goToPage(pageNum)}
                  style={{
                    padding: '0.5rem 0.75rem',
                    background: currentPage === pageNum ? '#667eea' : 'white',
                    color: currentPage === pageNum ? 'white' : '#4a5568',
                    border: '1px solid #e2e8f0',
                    borderRadius: '6px',
                    cursor: 'pointer',
                    fontSize: '0.875rem',
                    fontWeight: currentPage === pageNum ? '600' : '500',
                    minWidth: '40px'
                  }}
                >
                  {pageNum}
                </button>
              )
            })}
          </div>

          <button
            onClick={goToNextPage}
            disabled={currentPage === totalPages}
            style={{
              padding: '0.5rem 1rem',
              background: currentPage === totalPages ? '#e2e8f0' : '#667eea',
              color: currentPage === totalPages ? '#a0aec0' : 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: currentPage === totalPages ? 'not-allowed' : 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              fontSize: '0.875rem',
              fontWeight: '500'
            }}
          >
            Next
            <ChevronRight size={16} />
          </button>
        </div>
      )}
    </div>
  )
}
