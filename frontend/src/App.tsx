import { useState } from 'react';

function App() {
  // Nos mémoires (States)
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false); // Pour ouvrir/fermer la fenêtre de login
  const [isLoggedIn, setIsLoggedIn] = useState(false);   // Pour simuler si on est connecté ou non

  // Petite fonction pour "simuler" une connexion réussie
  const handleFakeLogin = (e: React.FormEvent) => {
    e.preventDefault(); // Empêche la page de se recharger (comportement par défaut d'un formulaire)
    setIsLoggedIn(true); // On dit au site "C'est bon, il est connecté !"
    setIsLoginOpen(false); // On ferme la fenêtre
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white flex flex-col relative overflow-hidden">

      <header className="absolute top-0 w-full p-6 flex justify-between items-center z-10">

        {/* Bouton Paramètres */}
        <button
          onClick={() => setIsSettingsOpen(true)}
          className="text-gray-400 hover:text-white transition-colors duration-300"
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </button>

        {/* AFFICHAGE CONDITIONNEL DU HEADER */}
        {/* Si on n'est PAS connecté, on affiche le bouton rouge. Sinon, on affiche le profil ! */}
        {!isLoggedIn ? (
          <button
            onClick={() => setIsLoginOpen(true)}
            className="bg-red-600 hover:bg-red-700 text-white font-bold py-3 px-8 rounded-lg transition-all shadow-[0_0_15px_rgba(220,38,38,0.5)] hover:shadow-[0_0_25px_rgba(220,38,38,0.8)]"
          >
            Login / Sign in
          </button>
        ) : (
          <div className="flex items-center space-x-4 cursor-pointer hover:bg-gray-800 p-2 rounded-lg transition-colors">
            <span className="font-bold text-lg">Maati_42</span>
            {/* Une fausse image de profil (Avatar par défaut) */}
            <div className="w-12 h-12 bg-gray-700 rounded-full border-2 border-red-500 overflow-hidden">
              <svg className="w-full h-full text-gray-400 mt-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
              </svg>
            </div>
          </div>
        )}

      </header>

      <main className="flex-1 flex flex-col items-center justify-center space-y-12">
        <h1 className="text-6xl md:text-8xl font-black text-transparent bg-clip-text bg-gradient-to-b from-red-400 to-red-800 tracking-widest drop-shadow-lg text-center">
          42<br />SURVIVOR
        </h1>

        <div className="flex flex-col space-y-6 w-72">
          <button className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-red-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all">Solo</button>
          <button className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-red-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all">Multi</button>
          <button className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-yellow-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all">Match History</button>
        </div>
      </main>

      {/* ==========================================
          MODALE PARAMÈTRES (Inchangée)
          ========================================== */}
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
              <div className="flex justify-between items-center pt-4"><span className="text-gray-400 font-semibold">Afficher les FPS</span><input type="checkbox" className="w-5 h-5 accent-red-600 rounded" /></div>
            </div>
          </div>
        </div>
      )}

      {/* ==========================================
          MODALE DE CONNEXION (LOGIN / SIGN IN)
          ========================================== */}
      {isLoginOpen && (
        <div className="absolute inset-0 bg-black/80 flex items-center justify-center z-50 backdrop-blur-sm">
          <div className="bg-gray-900 border-2 border-gray-700 p-8 rounded-2xl w-[400px] shadow-2xl relative">

            <button onClick={() => setIsLoginOpen(false)} className="absolute top-4 right-4 text-gray-500 hover:text-white transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
            </button>

            <h2 className="text-3xl font-bold text-center text-gray-200 mb-6">Connexion</h2>

            {/* Le formulaire classique */}
            <form onSubmit={handleFakeLogin} className="space-y-4">
              <div>
                <label className="block text-gray-400 mb-1 text-sm">Email</label>
                <input type="email" required className="w-full bg-gray-800 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-red-500" placeholder="joueur@42.fr" />
              </div>

              <div>
                <label className="block text-gray-400 mb-1 text-sm">Mot de passe</label>
                <input type="password" required className="w-full bg-gray-800 border border-gray-600 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-red-500" placeholder="••••••••" />
              </div>

              <button type="submit" className="w-full bg-red-600 hover:bg-red-700 text-white font-bold py-3 rounded-lg transition-colors mt-4">
                Se connecter
              </button>
            </form>

            <div className="my-6 flex items-center">
              <div className="flex-1 border-t border-gray-700"></div>
              <span className="px-4 text-gray-500 text-sm">OU</span>
              <div className="flex-1 border-t border-gray-700"></div>
            </div>

            {/* Le bouton Google */}
            <button onClick={handleFakeLogin} className="w-full bg-white text-gray-900 hover:bg-gray-200 font-bold py-3 rounded-lg transition-colors flex items-center justify-center space-x-2">
              <svg className="w-5 h-5" viewBox="0 0 24 24">
                <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
              </svg>
              <span>Continuer avec Google</span>
            </button>

            <p className="mt-4 text-center text-gray-500 text-sm">
              Pas encore de compte ? <span className="text-red-500 hover:underline cursor-pointer">S'inscrire</span>
            </p>

          </div>
        </div>
      )}

    </div>
  )
}

export default App;