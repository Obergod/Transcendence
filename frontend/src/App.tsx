import { useState } from 'react';
// NOUVEAU : On importe les outils de navigation et nos pages
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import Profile from './pages/Profile';
import Lobby from './pages/Lobby';
import Game from './pages/Game';
import MatchHistory from './pages/MatchHistory';
import Chat from './pages/Chat';

function App() {
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const handleFakeLogin = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoggedIn(true);
    setIsLoginOpen(false);
  };

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-900 text-white flex flex-col relative overflow-hidden">

        <header className="absolute top-0 w-full p-6 flex justify-between items-center z-10">

          <div className="flex space-x-4">
            <button onClick={() => setIsSettingsOpen(true)} className="text-gray-400 hover:text-white transition-colors duration-300">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            </button>
            <Link to="/" className="text-gray-400 hover:text-white font-bold text-xl flex items-center">
              Accueil
            </Link>
          </div>

          {isLoggedIn && (
          <Link to="/profile" className="flex items-center space-x-4 cursor-pointer hover:bg-gray-800 p-2 rounded-lg transition-colors">
            <span className="font-bold text-lg">Impots.gouv.fr</span>
            <div className="w-12 h-12 bg-gray-700 rounded-full border-2 border-red-500 overflow-hidden">
              <svg className="w-full h-full text-gray-400 mt-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
              </svg>
            </div>
          </Link>
        )}

        </header>

<Routes>
  <Route path="/" element={<Home isLoggedIn={isLoggedIn} onLoginClick={() => setIsLoginOpen(true)} />} />
  <Route path="/profile" element={<Profile onLogout={() => setIsLoggedIn(false)} />} />
  <Route path="/play" element={<Lobby />} />
  <Route path="/game" element={<Game />} />
  <Route path="/history" element={<MatchHistory />} />
  <Route path="/chat/:friendName" element={<Chat />} />
</Routes>

        {isSettingsOpen && (
          <div className="absolute inset-0 bg-black/80 flex items-center justify-center z-50 backdrop-blur-sm">
            <div className="bg-gray-900 border-2 border-gray-700 p-8 rounded-2xl w-96 shadow-2xl relative">
              <button onClick={() => setIsSettingsOpen(false)} className="absolute top-4 right-4 text-gray-500 hover:text-white transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
              <h2 className="text-3xl font-bold text-center text-gray-200 mb-8 border-b border-gray-700 pb-4">Paramètres</h2>
              <div className="space-y-6">
                <div><label className="block text-gray-400 mb-2 font-semibold">Volume de la musique</label><input type="range" className="w-full accent-red-600" /></div>
                <div><label className="block text-gray-400 mb-2 font-semibold">Volume des effets (SFX)</label><input type="range" className="w-full accent-red-600" /></div>
              </div>
            </div>
          </div>
        )}

        {isLoginOpen && (
          <div className="absolute inset-0 bg-black/80 flex items-center justify-center z-50 backdrop-blur-sm">
            <div className="bg-gray-900 border-2 border-gray-700 p-8 rounded-2xl w-[400px] shadow-2xl relative">
              <button onClick={() => setIsLoginOpen(false)} className="absolute top-4 right-4 text-gray-500 hover:text-white transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
              <h2 className="text-3xl font-bold text-center text-gray-200 mb-6">Connexion</h2>
              <form onSubmit={handleFakeLogin} className="space-y-4">
                <div><label className="block text-gray-400 mb-1 text-sm">Email</label><input type="email" required className="w-full bg-gray-800 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-red-500" placeholder="joueur@42.fr" /></div>
                <div><label className="block text-gray-400 mb-1 text-sm">Mot de passe</label><input type="password" required className="w-full bg-gray-800 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-red-500" placeholder="••••••••" /></div>
                <button type="submit" className="w-full bg-red-600 hover:bg-red-700 text-white font-bold py-3 rounded-lg transition-colors mt-4">Se connecter</button>
              </form>
              <div className="my-6 flex items-center"><div className="flex-1 border-t border-gray-700"></div><span className="px-4 text-gray-500 text-sm">OU</span><div className="flex-1 border-t border-gray-700"></div></div>
              <button onClick={handleFakeLogin} type="button" className="w-full bg-white text-gray-900 hover:bg-gray-200 font-bold py-3 rounded-lg transition-colors flex items-center justify-center space-x-2">
                <span>Continuer avec Google</span>
              </button>
            </div>
          </div>
        )}

      </div>
    </BrowserRouter>
  )
}

export default App;