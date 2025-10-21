import { useState, useEffect } from 'react'
import { Send, Sparkles, TrendingUp, AlertTriangle, Activity, Trash2 } from 'lucide-react'
import axios from 'axios'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { API_URL } from '../config'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

function formatAIResponse(data: any): string {
  if (!data) return 'No data received'

  let formatted = ''

  // Priority 1: Handle direct response field (from RAG chat)
  if (data.response) {
    formatted += data.response
    // Only show logs count if context was actually used
    if (data.context_used && data.logs_count !== undefined) {
      formatted += `\n\n---\n\n📈 *Analyzed ${data.logs_count} logs*`
    }
    return formatted
  }

  // Priority 2: Handle analysis.text field (from pattern detection, etc.)
  if (data.analysis?.text) {
    formatted += data.analysis.text
    if (data.logs_count !== undefined) {
      formatted += `\n\n---\n\n📈 *Analyzed ${data.logs_count} logs*`
    }
    return formatted
  }

  // Priority 3: Handle legacy analysis object
  if (data.analysis && typeof data.analysis === 'object' && !data.analysis.text) {
    formatted += `📊 **AI Analysis**\n\n${JSON.stringify(data.analysis, null, 2)}\n\n`
  }

  // Handle patterns
  if (data.patterns && data.patterns.length > 0) {
    formatted += '## 🔍 Detected Patterns\n\n'
    data.patterns.forEach((pattern: any, i: number) => {
      formatted += `**Pattern ${i + 1}:** ${pattern.description || pattern.error_pattern}\n`
      formatted += `- Occurrences: ${pattern.count || pattern.occurrences}\n`
      if (pattern.endpoints) {
        formatted += `- Affected Endpoints: ${pattern.endpoints.join(', ')}\n`
      }
      formatted += '\n'
    })
  }

  // Handle failing APIs
  if (data.failing_apis && data.failing_apis.length > 0) {
    formatted += '## ⚠️ Failing APIs\n\n'
    data.failing_apis.forEach((api: any) => {
      formatted += `- **${api.endpoint}**: ${api.error_count} errors (${(api.failure_rate * 100).toFixed(1)}% failure rate)\n`
    })
    formatted += '\n'
  }

  // Show log count only if we have formatted content
  if (formatted && data.logs_count !== undefined) {
    formatted += `\n📈 *Analyzed ${data.logs_count} logs*\n`
  }

  return formatted || JSON.stringify(data, null, 2)
}

export default function AIChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [timeRange, setTimeRange] = useState('1h')
  const [useCustomRange, setUseCustomRange] = useState(false)
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')

  // Load chat history from localStorage on mount
  useEffect(() => {
    const savedMessages = localStorage.getItem('aiChatHistory')
    if (savedMessages) {
      try {
        const parsed = JSON.parse(savedMessages)
        setMessages(parsed.map((m: any) => ({
          ...m,
          timestamp: new Date(m.timestamp)
        })))
      } catch (error) {
        console.error('Failed to load chat history:', error)
      }
    }
  }, [])

  // Save chat history to localStorage whenever messages change
  useEffect(() => {
    if (messages.length > 0) {
      localStorage.setItem('aiChatHistory', JSON.stringify(messages))
    }
  }, [messages])

  const clearHistory = () => {
    setMessages([])
    localStorage.removeItem('aiChatHistory')
  }

  const quickActions = [
    { label: 'Detect Patterns', icon: TrendingUp, endpoint: '/api/ai/detect-patterns' },
    { label: 'Find Failures', icon: AlertTriangle, endpoint: '/api/ai/analyze-failures' },
    { label: 'Traffic Analysis', icon: Activity, endpoint: '/api/ai/analyze-traffic-failures' },
  ]

  const sendMessage = async (message: string, endpoint = '/api/ai/chat', customBody?: any) => {
    if (!message.trim() && !customBody) return

    const userMessage: Message = {
      role: 'user',
      content: message || `Running: ${endpoint}`,
      timestamp: new Date()
    }

    setMessages(prev => [...prev, userMessage])
    setInput('')
    setLoading(true)

    try {
      let body = customBody || {
        message,
        time_range: useCustomRange ? 'custom' : timeRange,
        limit: 2000 // Increased from 200 to analyze more logs
      }

      // Add custom date range if selected
      if (useCustomRange && startDate && endDate) {
        body = {
          ...body,
          start_date: startDate,
          end_date: endDate
        }
      }

      const response = await axios.post(`${API_URL}${endpoint}`, body, {
        withCredentials: true
      })

      const assistantMessage: Message = {
        role: 'assistant',
        content: formatAIResponse(response.data.data),
        timestamp: new Date()
      }

      setMessages(prev => [...prev, assistantMessage])
    } catch (error: any) {
      const errorMessage: Message = {
        role: 'assistant',
        content: `Error: ${error.response?.data?.message || error.message}`,
        timestamp: new Date()
      }
      setMessages(prev => [...prev, errorMessage])
    } finally {
      setLoading(false)
    }
  }

  const handleQuickAction = (action: typeof quickActions[0]) => {
    let body: any = {
      time_range: useCustomRange ? 'custom' : timeRange,
      limit: 2000, // Increased from 1000 to analyze more logs
      min_failure_rate: 0.1
    }

    // Add custom date range if selected
    if (useCustomRange && startDate && endDate) {
      body.start_date = startDate
      body.end_date = endDate
    }

    sendMessage(action.label, action.endpoint, body)
  }

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column', background: '#f7fafc' }}>
      {/* Header */}
      <div style={{
        background: 'white',
        padding: '1.5rem 2rem',
        borderBottom: '1px solid #e2e8f0',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <div>
          <h1 style={{ fontSize: '2rem', fontWeight: 'bold', color: '#1a202c', marginBottom: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <Sparkles size={32} color="#667eea" />
            AI Chat Assistant
          </h1>
          <p style={{ color: '#718096' }}>
            Ask questions about your API logs or use quick actions below
          </p>
        </div>
        {messages.length > 0 && (
          <button
            onClick={clearHistory}
            style={{
              padding: '0.75rem 1rem',
              background: '#f56565',
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
            <Trash2 size={16} />
            Clear History
          </button>
        )}
      </div>

      {/* Quick Actions & Time Range */}
      <div style={{ padding: '1.5rem 2rem', background: 'white', borderBottom: '1px solid #e2e8f0' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.75rem' }}>
          <div style={{ fontSize: '0.875rem', fontWeight: '600', color: '#4a5568' }}>
            Quick Actions:
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexWrap: 'wrap' }}>
            <label style={{ fontSize: '0.875rem', color: '#4a5568', fontWeight: '500' }}>
              Time Range:
            </label>
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
                padding: '0.5rem 0.75rem',
                border: '1px solid #e2e8f0',
                borderRadius: '6px',
                fontSize: '0.875rem',
                cursor: 'pointer',
                background: 'white',
                color: '#2d3748'
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
            {useCustomRange && (
              <>
                <input
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  style={{
                    padding: '0.5rem 0.75rem',
                    border: '1px solid #e2e8f0',
                    borderRadius: '6px',
                    fontSize: '0.875rem',
                    background: 'white',
                    color: '#2d3748',
                    cursor: 'pointer'
                  }}
                />
                <span style={{ color: '#718096' }}>to</span>
                <input
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                  style={{
                    padding: '0.5rem 0.75rem',
                    border: '1px solid #e2e8f0',
                    borderRadius: '6px',
                    fontSize: '0.875rem',
                    background: 'white',
                    color: '#2d3748',
                    cursor: 'pointer'
                  }}
                />
              </>
            )}
          </div>
        </div>
        <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
          {quickActions.map((action) => {
            const Icon = action.icon
            return (
              <button
                key={action.label}
                onClick={() => handleQuickAction(action)}
                disabled={loading}
                style={{
                  padding: '0.75rem 1rem',
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  color: 'white',
                  border: 'none',
                  borderRadius: '8px',
                  cursor: loading ? 'not-allowed' : 'pointer',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '0.5rem',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  opacity: loading ? 0.6 : 1
                }}
              >
                <Icon size={16} />
                {action.label}
              </button>
            )
          })}
        </div>
      </div>

      {/* Messages */}
      <div style={{ flex: 1, overflowY: 'auto', padding: '2rem' }}>
        {messages.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '4rem 2rem', color: '#a0aec0' }}>
            <Sparkles size={64} style={{ margin: '0 auto 1rem', opacity: 0.3 }} />
            <h2 style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>Start a conversation</h2>
            <p>Ask me anything about your API logs or use the quick actions above</p>
            <div style={{ marginTop: '2rem', textAlign: 'left', maxWidth: '600px', margin: '2rem auto 0' }}>
              <p style={{ fontWeight: '600', marginBottom: '0.5rem', color: '#4a5568' }}>Example questions:</p>
              <ul style={{ listStyle: 'none', padding: 0 }}>
                <li style={{ padding: '0.5rem 0', color: '#718096' }}>• "Show me all failing APIs"</li>
                <li style={{ padding: '0.5rem 0', color: '#718096' }}>• "Why is /api/users/list failing?"</li>
                <li style={{ padding: '0.5rem 0', color: '#718096' }}>• "Which endpoints have high latency?"</li>
                <li style={{ padding: '0.5rem 0', color: '#718096' }}>• "Find patterns in error logs"</li>
              </ul>
            </div>
          </div>
        ) : (
          <div style={{ maxWidth: '900px', margin: '0 auto' }}>
            {messages.map((message, index) => (
              <div
                key={index}
                style={{
                  marginBottom: '1.5rem',
                  display: 'flex',
                  justifyContent: message.role === 'user' ? 'flex-end' : 'flex-start'
                }}
              >
                <div
                  style={{
                    maxWidth: '80%',
                    padding: '1rem 1.25rem',
                    borderRadius: '12px',
                    background: message.role === 'user'
                      ? 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
                      : 'white',
                    color: message.role === 'user' ? 'white' : '#2d3748',
                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                  }}
                >
                  <div style={{ fontSize: '0.75rem', opacity: 0.7, marginBottom: '0.5rem' }}>
                    {message.role === 'user' ? 'You' : 'AI Assistant'} • {message.timestamp.toLocaleTimeString()}
                  </div>
                  <div style={{
                    fontSize: '0.9375rem',
                    lineHeight: '1.7'
                  }} className="markdown-content">
                    {message.role === 'assistant' ? (
                      <ReactMarkdown remarkPlugins={[remarkGfm]}>
                        {message.content}
                      </ReactMarkdown>
                    ) : (
                      message.content
                    )}
                  </div>
                </div>
              </div>
            ))}
            {loading && (
              <div style={{ textAlign: 'center', padding: '1rem', color: '#a0aec0' }}>
                <div style={{ display: 'inline-block', animation: 'pulse 1.5s ease-in-out infinite' }}>
                  AI is thinking...
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Input */}
      <div style={{
        background: 'white',
        padding: '1.5rem 2rem',
        borderTop: '1px solid #e2e8f0',
        boxShadow: '0 -1px 3px rgba(0,0,0,0.1)'
      }}>
        <form
          onSubmit={(e) => {
            e.preventDefault()
            sendMessage(input)
          }}
          style={{ maxWidth: '900px', margin: '0 auto', display: 'flex', gap: '1rem' }}
        >
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ask a question about your logs..."
            disabled={loading}
            style={{
              flex: 1,
              padding: '1rem',
              border: '1px solid #e2e8f0',
              borderRadius: '8px',
              fontSize: '1rem',
              boxSizing: 'border-box'
            }}
          />
          <button
            type="submit"
            disabled={loading || !input.trim()}
            style={{
              padding: '1rem 2rem',
              background: loading || !input.trim() ? '#a0aec0' : 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              cursor: loading || !input.trim() ? 'not-allowed' : 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem',
              fontSize: '1rem',
              fontWeight: '600'
            }}
          >
            <Send size={20} />
            Send
          </button>
        </form>
      </div>
    </div>
  )
}
