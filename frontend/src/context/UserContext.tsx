import * as React from 'react';
import { createContext, useContext, useState, useEffect, useCallback } from 'react';

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
  ws: WebSocket | null;
  levelInfo: { level: number; xpInLevel: number; xpForNext: number } | null;
  refreshLevel : () => Promise<void>;
  onlineUsers: number[];
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isReady, setIsReady] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [levelInfo, setLevelInfo] = useState<{ level: number; xpInLevel: number; xpForNext: number } | null>(null);

  // État qui va contenir les IDs des joueurs connectés
  const [onlineUsers, setOnlineUsers] = useState<number[]>([]);

  useEffect(() => {
    const checkToken = async () => {
      const token = localStorage.getItem('jwt_token');
      if (token) {
        try {
          const response = await fetch("/api/user/me", {
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

  // gestion du websocket
  useEffect(() => {
    let socket: WebSocket | null = null;
    let connectTimer: ReturnType<typeof setTimeout>;

    if (user) {
        const token = localStorage.getItem('jwt_token');

        connectTimer = setTimeout(() => {
            socket = new WebSocket(`wss://localhost:5173/ws?token=${token}`);

            socket.onopen = () => {
                console.log("🟢 WebSocket Connecté !");
            };

            // CORRECTION : On remet l'écouteur de messages qui met à jour la liste !
            socket.onmessage = (event: MessageEvent) => {
              try {
                const data = JSON.parse(event.data);
                if (data.type === "online_users") {
                  setOnlineUsers(data.users || []);
                }
              } catch (err) {
                console.error("Erreur de parsing WS", err);
              }
            };

            socket.onclose = () => {
                console.log("🔴 WebSocket Déconnecté !");
                setOnlineUsers([]); // On vide la liste à la déconnexion
            };

            setWs(socket);
        }, 150);
    }

    return () => {
        clearTimeout(connectTimer);
        if (socket) {
            socket.close();
            setWs(null);
        }
    };
  }, [user]);

  const refreshLevel = useCallback(async () => {
    if (!user) { setLevelInfo(null); return; }
    const token = localStorage.getItem('jwt_token');
    try {
      const res = await fetch("/api/user/level", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data = await res.json();
        setLevelInfo({ level: data.level, xpInLevel: data.xp_in_level, xpForNext: data.xp_for_next });
      }
    } catch (err) {
      console.error("Erreur chargement niveau", err);
    }
  }, [user]);

  useEffect(() => { refreshLevel(); }, [refreshLevel]);

  const login = (userData: User) => setUser(userData);

  const logout = () => {
    setUser(null);
    if (ws) ws.close();
  };

  if (!isReady) {
    return <div className="min-h-screen bg-gray-900 flex items-center justify-center text-white font-bold">Chargement du profil...</div>;
  }

  return (
    <UserContext.Provider value={{ user, login, logout, isReady, ws, levelInfo, refreshLevel, onlineUsers }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) throw new Error("useUser doit être utilisé dans un UserProvider");
  return context;
};