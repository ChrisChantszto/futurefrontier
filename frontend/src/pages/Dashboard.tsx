import { useState, useEffect } from 'react'
import { Activity, AlertTriangle, TrendingUp, Clock, Database, Zap } from 'lucide-react'
import { Link } from 'react-router-dom'
import axios from 'axios'

export default function Dashboard() {
  const [stats, setStats] = useState({
    totalLogs: 0,
    errorRate: 0,
    avgLatency: 0,
    activeEndpoints: 0
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/logs/stats', {
          withCredentials: true
        })

        if (response.data.success) {
          const data = response.data.data
          setStats({
            totalLogs: data.total_logs || 0,
            errorRate: data.error_rate || 0,
            avgLatency: data.avg_latency || 0,
            activeEndpoints: data.active_endpoints || 0
          })
        }
      } catch (error) {
        console.error('Failed to fetch stats:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchStats()
  }, [])

  const statCards = [
    {
      title: 'Total Logs',
      value: stats.totalLogs.toLocaleString(),
      icon: Database,
      color: '#667eea',
      bgColor: '#667eea20'
    },
    {
      title: 'Error Rate',
      value: `${stats.errorRate.toFixed(1)}%`,
      icon: AlertTriangle,
      color: '#f56565',
      bgColor: '#f5656520'
    },
    {
      title: 'Avg Latency',
      value: `${stats.avgLatency}ms`,
      icon: Clock,
      color: '#ed8936',
      bgColor: '#ed893620'
    },
    {
      title: 'Active Endpoints',
      value: stats.activeEndpoints.toString(),
      icon: Activity,
      color: '#48bb78',
      bgColor: '#48bb7820'
    }
  ]

  const quickActions = [
    {
      title: 'Discover Logs',
      description: 'Search and filter API logs',
      icon: Database,
      link: '/discover',
      color: '#667eea'
    },
    {
      title: 'AI Analysis',
      description: 'Chat with AI about your logs',
      icon: Zap,
      link: '/ai-chat',
      color: '#764ba2'
    },
    {
      title: 'Detect Patterns',
      description: 'Find recurring failures',
      icon: TrendingUp,
      link: '/ai-chat',
      color: '#48bb78'
    }
  ]

  return (
    <div style={{ padding: '2rem', maxWidth: '1400px', margin: '0 auto' }}>
      <div style={{ marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: 'bold', color: '#1a202c', marginBottom: '0.5rem' }}>
          Dashboard
        </h1>
        <p style={{ color: '#718096' }}>
          Overview of your API logs and system health
        </p>
      </div>

      {/* Stats Grid */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
        gap: '1.5rem',
        marginBottom: '2rem'
      }}>
        {statCards.map((stat) => {
          const Icon = stat.icon
          return (
            <div
              key={stat.title}
              style={{
                background: 'white',
                padding: '1.5rem',
                borderRadius: '12px',
                boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
                display: 'flex',
                alignItems: 'center',
                gap: '1rem'
              }}
            >
              <div style={{
                width: '64px',
                height: '64px',
                borderRadius: '12px',
                background: stat.bgColor,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
                <Icon size={32} color={stat.color} />
              </div>
              <div>
                <div style={{ fontSize: '0.875rem', color: '#718096', marginBottom: '0.25rem' }}>
                  {stat.title}
                </div>
                <div style={{ fontSize: '2rem', fontWeight: 'bold', color: '#1a202c' }}>
                  {loading ? '...' : stat.value}
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Quick Actions */}
      <div style={{ marginBottom: '2rem' }}>
        <h2 style={{ fontSize: '1.5rem', fontWeight: 'bold', color: '#1a202c', marginBottom: '1rem' }}>
          Quick Actions
        </h2>
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
          gap: '1.5rem'
        }}>
          {quickActions.map((action) => {
            const Icon = action.icon
            return (
              <Link
                key={action.title}
                to={action.link}
                style={{
                  background: 'white',
                  padding: '1.5rem',
                  borderRadius: '12px',
                  boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
                  textDecoration: 'none',
                  display: 'flex',
                  alignItems: 'start',
                  gap: '1rem',
                  transition: 'all 0.2s',
                  border: '2px solid transparent'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = action.color
                  e.currentTarget.style.boxShadow = `0 4px 12px ${action.color}40`
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = 'transparent'
                  e.currentTarget.style.boxShadow = '0 1px 3px rgba(0,0,0,0.1)'
                }}
              >
                <div style={{
                  width: '48px',
                  height: '48px',
                  borderRadius: '8px',
                  background: `${action.color}20`,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  flexShrink: 0
                }}>
                  <Icon size={24} color={action.color} />
                </div>
                <div>
                  <h3 style={{ fontSize: '1.125rem', fontWeight: '600', color: '#1a202c', marginBottom: '0.25rem' }}>
                    {action.title}
                  </h3>
                  <p style={{ fontSize: '0.875rem', color: '#718096', margin: 0 }}>
                    {action.description}
                  </p>
                </div>
              </Link>
            )
          })}
        </div>
      </div>

      {/* Getting Started */}
      <div style={{
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '2rem',
        borderRadius: '12px',
        color: 'white'
      }}>
        <h2 style={{ fontSize: '1.5rem', fontWeight: 'bold', marginBottom: '1rem' }}>
          🚀 Getting Started
        </h2>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1.5rem' }}>
          <div>
            <h3 style={{ fontSize: '1.125rem', fontWeight: '600', marginBottom: '0.5rem' }}>
              1. Generate Demo Data
            </h3>
            <p style={{ fontSize: '0.875rem', opacity: 0.9 }}>
              Create sample logs to test the system
            </p>
          </div>
          <div>
            <h3 style={{ fontSize: '1.125rem', fontWeight: '600', marginBottom: '0.5rem' }}>
              2. Explore Discover
            </h3>
            <p style={{ fontSize: '0.875rem', opacity: 0.9 }}>
              Search and filter logs like Kibana
            </p>
          </div>
          <div>
            <h3 style={{ fontSize: '1.125rem', fontWeight: '600', marginBottom: '0.5rem' }}>
              3. Ask AI Questions
            </h3>
            <p style={{ fontSize: '0.875rem', opacity: 0.9 }}>
              Use AI to find patterns and root causes
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
