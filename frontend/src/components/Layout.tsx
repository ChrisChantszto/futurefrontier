import { Outlet, Link, useLocation } from 'react-router-dom'
import { LayoutDashboard, Search, MessageSquare, LogOut } from 'lucide-react'

interface LayoutProps {
  onLogout: () => void
}

export default function Layout({ onLogout }: LayoutProps) {
  const location = useLocation()

  const navItems = [
    { path: '/', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/discover', label: 'Discover', icon: Search },
    { path: '/ai-chat', label: 'AI Chat', icon: MessageSquare },
  ]

  return (
    <div style={{ display: 'flex', height: '100vh', background: '#f7fafc' }}>
      {/* Sidebar */}
      <div style={{
        width: '240px',
        background: 'linear-gradient(180deg, #667eea 0%, #764ba2 100%)',
        color: 'white',
        display: 'flex',
        flexDirection: 'column'
      }}>
        <div style={{ padding: '1.5rem', borderBottom: '1px solid rgba(255,255,255,0.1)' }}>
          <h1 style={{ fontSize: '1.25rem', fontWeight: 'bold', margin: 0 }}>
            🤖 AI Log Analysis
          </h1>
          <p style={{ fontSize: '0.75rem', opacity: 0.8, marginTop: '0.25rem' }}>
            Powered by Elasticsearch + Vertex AI
          </p>
        </div>

        <nav style={{ flex: 1, padding: '1rem' }}>
          {navItems.map((item) => {
            const Icon = item.icon
            const isActive = location.pathname === item.path
            return (
              <Link
                key={item.path}
                to={item.path}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  padding: '0.75rem 1rem',
                  marginBottom: '0.5rem',
                  borderRadius: '8px',
                  textDecoration: 'none',
                  color: 'white',
                  background: isActive ? 'rgba(255,255,255,0.2)' : 'transparent',
                  transition: 'all 0.2s'
                }}
                onMouseEnter={(e) => {
                  if (!isActive) e.currentTarget.style.background = 'rgba(255,255,255,0.1)'
                }}
                onMouseLeave={(e) => {
                  if (!isActive) e.currentTarget.style.background = 'transparent'
                }}
              >
                <Icon size={20} style={{ marginRight: '0.75rem' }} />
                <span style={{ fontWeight: isActive ? '600' : '400' }}>{item.label}</span>
              </Link>
            )
          })}
        </nav>

        <div style={{ padding: '1rem', borderTop: '1px solid rgba(255,255,255,0.1)' }}>
          <button
            onClick={onLogout}
            style={{
              display: 'flex',
              alignItems: 'center',
              width: '100%',
              padding: '0.75rem 1rem',
              background: 'rgba(255,255,255,0.1)',
              border: 'none',
              borderRadius: '8px',
              color: 'white',
              cursor: 'pointer',
              fontSize: '1rem',
              transition: 'all 0.2s'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'rgba(255,255,255,0.2)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'rgba(255,255,255,0.1)'}
          >
            <LogOut size={20} style={{ marginRight: '0.75rem' }} />
            <span>Logout</span>
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div style={{ flex: 1, overflow: 'auto' }}>
        <Outlet />
      </div>
    </div>
  )
}
