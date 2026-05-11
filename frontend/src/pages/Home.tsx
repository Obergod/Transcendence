import { useState } from 'react';
import { Link } from 'react-router-dom';

interface HomeProps {
  isLoggedIn: boolean;
  onLoginClick: () => void;
}

const Home = ({ isLoggedIn, onLoginClick }: HomeProps) => {
  const [backendMessage, setBackendMessage] = useState<string | null>(null);

  const testBackendConnection = async () => {
    try {
      const response = await fetch('/api/hello');
      const data = await response.json();
      setBackendMessage(data.message);
    } catch (error) {
      console.error("Erreur de connexion", error);
      setBackendMessage("Erreur de connexion ❌");
    }
  };

  return (
    <main className="flex-1 flex flex-col items-center justify-center space-y-12">
      <h1 className="text-6xl md:text-8xl font-black text-transparent bg-clip-text bg-gradient-to-b from-red-400 to-red-800 tracking-widest drop-shadow-lg text-center">
        42<br />SURVIVOR
      </h1>

      <div className="bg-gray-800/80 p-4 rounded-xl border border-gray-700 flex flex-col items-center gap-2">
        <button onClick={testBackendConnection} className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition-colors text-sm">
          Tester la connexion Backend
        </button>
        {backendMessage && <p className="text-green-400 font-mono text-sm">{backendMessage}</p>}
      </div>

      <div className="flex flex-col space-y-6 w-72">
        {isLoggedIn ? (
          <>
            {/* NOUVEAU : Un seul gros bouton PLAY */}
            <Link to="/play" className="bg-red-600 hover:bg-red-700 border-2 border-red-500 py-4 rounded-xl text-2xl font-black uppercase tracking-widest transition-all text-center shadow-[0_0_20px_rgba(220,38,38,0.4)] text-white">
              PLAY
            </Link>
            <Link to="/history" className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-yellow-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center">
              Match History
            </Link>
          </>
        ) : (
          <button onClick={onLoginClick} className="bg-red-600 hover:bg-red-700 border-2 border-red-500 py-4 rounded-xl text-xl font-bold uppercase tracking-widest transition-all text-white shadow-[0_0_20px_rgba(220,38,38,0.5)]">
            Se connecter pour jouer
          </button>
        )}
      </div>
    </main>
  );
};

export default Home;