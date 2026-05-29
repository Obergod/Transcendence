import * as React from 'react';
import { createContext, useContext, useState, useEffect } from 'react';

interface User {
  id: number;
  pseudo: string;
  email: string;
  avatarUrl: string;
  status: string;
}

interface UserContextType {
  user: User | null;
  login: (userData: User) => void;
  logout: () => void;
  isReady: boolean;
  ws: WebSocket | null; //socket global
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isReady, setIsReady] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const checkToken = async () => {
      const token = localStorage.getItem('jwt_token');
      if (token) {
        try {
          const response = await fetch("http://localhost:8081/api/user/me", {
            headers: { "Authorization": `Bearer ${token}` }
          });
          if (response.ok) {
            const userData = await response.json();
            setUser(userData);
          } else {
            localStorage.removeItem('jwt_token');
          }
        } catch (error) {
          console.error("Erreur de restauration de session", error);
        }
      }
      setIsReady(true);
    };
    checkToken();
  }, []);

  // gestion du websoket
  useEffect(() => {
    let socket: WebSocket | null = null;

    if (user) {
      const token = localStorage.getItem('jwt_token');
      socket = new WebSocket(`ws://localhost:8081/ws?token=${token}`);

      socket.onopen = () => {
        console.log("🟢 WebSocket Connecté !");
      };

      socket.onclose = () => {
        console.log("🔴 WebSocket Déconnecté !");
      };

      setWs(socket);
    }

    // nettoyage
    return () => {
      if (socket) {
        socket.close();
        setWs(null);
      }
    };
  }, [user]);

  const login = (userData: User) => setUser(userData);

  const logout = () => {
    setUser(null);
    if (ws) ws.close();
  };

  if (!isReady) {
    return <div className="min-h-screen bg-gray-900 flex items-center justify-center text-white font-bold">Chargement du profil...</div>;
  }

  return (
    <UserContext.Provider value={{ user, login, logout, isReady, ws }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) throw new Error("useUser doit être utilisé dans un UserProvider");
  return context;
};