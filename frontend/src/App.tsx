import { useState, useEffect } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Discover from './pages/Discover'
import AIChat from './pages/AIChat'
import Login from './pages/Login'
import { API_URL } from './config'
import './App.css'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Check if user is already logged in
    const checkAuth = async () => {
      try {
        const response = await fetch(`${API_URL}/api/healthz`, {
          credentials: 'include'
        })
        if (response.ok) {
          // User might be authenticated, but we'll verify on protected routes
          setIsAuthenticated(true)
        }
      } catch (error) {
        console.error('Auth check failed:', error)
      } finally {
        setLoading(false)
      }
    }
    checkAuth()
  }, [])

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <div>Loading...</div>
      </div>
    )
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={
          isAuthenticated ? <Navigate to="/" /> : <Login onLogin={() => setIsAuthenticated(true)} />
        } />
        <Route path="/" element={
          isAuthenticated ? <Layout onLogout={() => setIsAuthenticated(false)} /> : <Navigate to="/login" />
        }>
          <Route index element={<Dashboard />} />
          <Route path="discover" element={<Discover />} />
          <Route path="ai-chat" element={<AIChat />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
