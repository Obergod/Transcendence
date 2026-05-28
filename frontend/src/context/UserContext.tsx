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
  isReady: boolean; // Pour éviter les clignotements d'écran pendant la vérification
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isReady, setIsReady] = useState(false);

  // VÉRIFICATION AUTOMATIQUE AU DÉMARRAGE
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
            setUser(userData); // On restaure l'utilisateur !
          } else {
            localStorage.removeItem('jwt_token'); // Le token est expiré ou faux
          }
        } catch (error) {
          console.error("Erreur de restauration de session", error);
        }
      }
      setIsReady(true); // L'application peut s'afficher
    };

    checkToken();
  }, []);

  const login = (userData: User) => setUser(userData);
  const logout = () => setUser(null);

  // On n'affiche rien tant qu'on ne sait pas si l'utilisateur est connecté ou non
  if (!isReady) {
    return <div className="min-h-screen bg-gray-900 flex items-center justify-center text-white">Chargement...</div>;
  }

  return (
    <UserContext.Provider value={{ user, login, logout, isReady }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) throw new Error("useUser doit être utilisé dans un UserProvider");
  return context;
};