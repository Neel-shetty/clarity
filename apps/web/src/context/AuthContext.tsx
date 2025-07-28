import { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';
import type { ReactNode } from 'react';

type User = {
  id : string;
  username : string;
}
type AuthContextType = {
  user: User | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => {},
  logout: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  useEffect(() => {
    axios.get<User>('/user/me')
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