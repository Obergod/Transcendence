import { useState } from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useUser } from './context/UserContext'; // <-- L'import manquant est ici !

// IMPORTS DE PAGES
import Home from './pages/Home';
import Profile from './pages/Profile';
import Lobby from './pages/Lobby';
import Game from './pages/Game';
import MatchHistory from './pages/MatchHistory';
import Chat from './pages/Chat';
import AuthModal from './components/AuthModal';

// IMPORTS DE WIDGETS
import LanguageSwitcher from './components/LanguageSwitcher';

function App() {
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);
  const { t } = useTranslation();

  // MAGIE : On récupère l'utilisateur depuis notre Contexte !
  const { user, levelInfo } = useUser();

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-900 text-white flex flex-col relative overflow-hidden">

        {/* --- HEADER --- */}
        <header className="absolute top-0 w-full p-6 flex justify-between items-center z-10">
          <div className="flex space-x-4">
            <button onClick={() => setIsSettingsOpen(true)} className="text-gray-400 hover:text-white transition-colors duration-300">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            </button>
            <Link to="/" className="text-gray-400 hover:text-white font-bold text-xl flex items-center">
              {t('nav.home')}
            </Link>
          </div>

          {/* <-- La correction est ici (on utilise 'user' au lieu de 'isLoggedIn') */}
          {user ? (
            <Link to="/profile" className="flex items-center space-x-4 cursor-pointer hover:bg-gray-800 p-2 rounded-lg transition-colors">
              <div className="flex flex-col items-end leading-tight">
                <span className="font-bold text-lg">{user.pseudo}</span>
                <span className="text-red-500 font-bold text-sm">
                  {levelInfo ? `${t('level.short')} ${levelInfo.level}` : ''}
                </span>
              </div>
              <div className="w-12 h-12 bg-gray-700 rounded-full border-2 border-red-500 overflow-hidden">
                <img src={user.avatarUrl} alt="Avatar" className="w-full h-full object-cover" />
              </div>
            </Link>
          ) : (
            <button
              onClick={() => setIsAuthModalOpen(true)}
              className="text-gray-400 hover:text-white font-bold text-xl flex items-center"
            >
              {t('nav.login')}
            </button>
          )}
        </header>

        {/* --- ROUTES --- */}
        <Routes>
          <Route path="/" element={<Home isLoggedIn={!!user} onLoginClick={() => setIsAuthModalOpen(true)} />} />
          <Route path="/profile" element={<Profile onLogout={() => console.log("Déconnexion en attente du contexte !")} />} />
          <Route path="/play" element={<Lobby />} />
          <Route path="/game" element={<Game />} />
          <Route path="/history" element={<MatchHistory />} />
		  <Route path="/chat/:friendId/:friendName" element={<Chat />} />
          <Route path="/chat/:friendName" element={<Chat />} />
        </Routes>

        {/* --- MODALES --- */}
        {isSettingsOpen && (
          <div className="absolute inset-0 bg-black/80 flex items-center justify-center z-50 backdrop-blur-sm">
            <div className="bg-gray-900 border-2 border-gray-700 p-8 rounded-2xl w-96 shadow-2xl relative">
              <button onClick={() => setIsSettingsOpen(false)} className="absolute top-4 right-4 text-gray-500 hover:text-white transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
              <h2 className="text-3xl font-bold text-center text-gray-200 mb-8 border-b border-gray-700 pb-4">{t('settings.title')}</h2>
              <div className="space-y-6">
                <div>
                  <label className="block text-gray-400 mb-2 font-semibold">{t('settings.language')}</label>
                  <LanguageSwitcher />
                </div>
              </div>
            </div>
          </div>
        )}

        {isAuthModalOpen && (
          <AuthModal
            onClose={() => setIsAuthModalOpen(false)}
            onLoginSuccess={() => setIsAuthModalOpen(false)}
          />
        )}

      </div>
    </BrowserRouter>
  )
}

export default App;
