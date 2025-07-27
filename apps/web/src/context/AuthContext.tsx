import { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';
import type { ReactNode } from 'react';

type AuthContextType = {
  user: any;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => {},
  logout: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState(null);
  useEffect(() => {
    axios.get('/user/me')
      .then(res => setUser(res.data))
      .catch(() => setUser(null));
  }, []);

  const login = async (username:string, password:string) => {
    const res = await axios.post('/user/login', { username, password });
    setUser(res.data);
  };

  const logout = async () => {
    await axios.post('/user/logout');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}