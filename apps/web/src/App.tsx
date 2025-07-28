import { Route,Routes} from 'react-router-dom'
import LoginPage from "./pages/LoginPage"
import { ThemeProvider } from './components/ui/theme-provider'
import './App.css'
import Signup from './pages/Signup'
import ProtectedRoute from './components/nav/ProtectedRoute'
import Home from "./pages/Home"
import NotFound from './pages/NotFound'
import { AuthProvider } from './context/AuthContext'

function App() {
  return (
    <AuthProvider>
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <Routes>
          <Route path="/" element={<LoginPage/>} />
          <Route path="/signup" element={<Signup/>} /> 
          <Route path="/home" element={<ProtectedRoute> <Home/> </ProtectedRoute>}></Route>
          <Route path="*" element={<NotFound/>}/>
        </Routes>
      </ThemeProvider>
    </AuthProvider>
  )
}

export default App
