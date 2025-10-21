import { useState, useEffect } from 'react'
import { Send, Sparkles, TrendingUp, AlertTriangle, Activity, Trash2 } from 'lucide-react'
import axios from 'axios'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

function formatAIResponse(data: any): string {
  if (!data) return 'No data received'

  let formatted = ''

  // Handle traffic analysis response
  if (data.analysis?.text) {
    formatted += '📊 **AI Analysis**\n\n'
    formatted += data.analysis.text.replace(/\*\*/g, '').substring(0, 1500)
    if (data.analysis.text.length > 1500) formatted += '...\n\n'
    formatted += '\n\n'
  }

  // Handle failing APIs
  if (data.failing_apis && data.failing_apis.length > 0) {
    formatted += '⚠️ **Failing APIs**\n\n'
    data.failing_apis.slice(0, 5).forEach((api: any, i: number) => {
      formatted += `${i + 1}. ${api.endpoint}\n`
      formatted += `   • Failure Rate: ${(api.failure_rate * 100).toFixed(1)}%\n`
      formatted += `   • Error Count: ${api.error_count}/${api.total_requests}\n`
      formatted += `   • Avg Latency: ${Math.round(api.avg_latency?.value || 0)}ms\n\n`
    })
  }

  // Handle patterns
  if (data.patterns && data.patterns.length > 0) {
    formatted += '🔍 **Detected Patterns**\n\n'
    data.patterns.slice(0, 5).forEach((pattern: any, i: number) => {
      formatted += `${i + 1}. ${pattern.endpoint || pattern.error_pattern}\n`
      formatted += `   • Occurrences: ${pattern.count || pattern.occurrences}\n`
      if (pattern.failure_rate) formatted += `   • Failure Rate: ${(pattern.failure_rate * 100).toFixed(1)}%\n`
      if (pattern.error_message) formatted += `   • Error: ${pattern.error_message}\n`
      formatted += '\n'
    })
  }

  // Handle recommendations
  if (data.recommendations && data.recommendations.length > 0) {
    formatted += '💡 **Recommendations**\n\n'
    data.recommendations.slice(0, 5).forEach((rec: any, i: number) => {
      formatted += `${i + 1}. ${rec}\n`
    })
    formatted += '\n'
  }

  // Handle root cause
  if (data.root_cause) {
    formatted += '🎯 **Root Cause**\n\n'
    formatted += data.root_cause + '\n\n'
  }

  // Show log count
  if (data.logs_count !== undefined) {
    formatted += `\n📈 Analyzed ${data.logs_count} logs\n`
  }

  return formatted || JSON.stringify(data, null, 2)
}

export default function AIChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)

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
    { label: 'Detect Patterns', icon: TrendingUp, endpoint: '/api/ai/detect-patterns', body: { time_range: '1h', limit: 500 } },
    { label: 'Find Failures', icon: AlertTriangle, endpoint: '/api/ai/analyze-failures', body: { time_range: '1h' } },
    { label: 'Traffic Analysis', icon: Activity, endpoint: '/api/ai/analyze-traffic-failures', body: { time_range: '1h', min_failure_rate: 0.1 } },
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
      const body = customBody || {
        message,
        time_range: '1h',
        limit: 100
      }

      const response = await axios.post(`http://localhost:8080${endpoint}`, body, {
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
    sendMessage(action.label, action.endpoint, action.body)
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

      {/* Quick Actions */}
      <div style={{ padding: '1.5rem 2rem', background: 'white', borderBottom: '1px solid #e2e8f0' }}>
        <div style={{ fontSize: '0.875rem', fontWeight: '600', color: '#4a5568', marginBottom: '0.75rem' }}>
          Quick Actions:
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
