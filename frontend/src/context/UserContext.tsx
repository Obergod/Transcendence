import * as React from 'react';
import { createContext, useContext, useState } from 'react';
import mockData from '../mock.json';

// 1. On définit à quoi ressemble notre utilisateur (TypeScript)
interface User {
  id: number;
  pseudo: string;
  email: string;
  avatarUrl: string;
  status: string;
}

interface UserContextType {
  user: User | null;
}

// 2. On crée le contexte
const UserContext = createContext<UserContextType | undefined>(undefined);

// 3. Le "Provider" qui va englober ton application
export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  // On charge directement les fausses données avec la bonne clé "user"
  const [user] = useState<User>(mockData.user);

  return (
    <UserContext.Provider value={{ user }}>
      {children}
    </UserContext.Provider>
  );
};

// 4. Un petit outil pour récupérer l'utilisateur facilement dans tes pages
export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error("useUser doit être utilisé à l'intérieur d'un UserProvider");
  }
  return context;
};