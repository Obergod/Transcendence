import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const Lobby = () => {
  const navigate = useNavigate();
  const [view, setView] = useState<'menu' | 'create' | 'join'>('menu');
  const [joinCode, setJoinCode] = useState("");
  const [generatedCode, setGeneratedCode] = useState("");

  const { t } = useTranslation();

  const onlineFriends = [
    { id: 1, pseudo: "Rillettes de conde" },
    { id: 3, pseudo: "Arthur" },
  ];

  const handleCreateLobby = () => {
    const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789";
    let result = "";
    for (let i = 0; i < 5; i++) result += chars.charAt(Math.floor(Math.random() * chars.length));
    setGeneratedCode(result);
    setView('create');
  };

  const handleStartGame = () => {
    navigate('/game');
  };

  return (
    <main className="flex-1 flex flex-col items-center pt-32 px-4 pb-12 w-full">
      <div className="w-full max-w-2xl bg-gray-800/50 border-2 border-gray-700 rounded-3xl p-8 backdrop-blur-md shadow-2xl">

        {view === 'menu' && (
          <div className="flex flex-col space-y-6">
            <h2 className="text-4xl font-black text-center mb-6 tracking-tight">{t('lobby.title')}<span className="text-red-500">{t('lobby.title_highlight')}</span></h2>
            <button onClick={handleCreateLobby} className="bg-red-600 hover:bg-red-700 text-white font-bold py-6 rounded-2xl text-2xl transition-all shadow-[0_0_20px_rgba(220,38,38,0.3)]">
              {t('lobby.create')}
            </button>
            <button onClick={() => setView('join')} className="bg-gray-700 hover:bg-gray-600 text-white font-bold py-6 rounded-2xl text-2xl transition-all">
              {t('lobby.join_with_code')}
            </button>
            <Link to="/" className="text-center text-gray-500 hover:text-white transition-colors pt-4 font-bold">{t('lobby.back_to_menu')}</Link>
          </div>
        )}

        {view === 'create' && (
          <div className="space-y-8">
            <div className="text-center">
              <h3 className="text-gray-400 font-bold uppercase tracking-widest text-sm mb-2">{t('lobby.lobby_code_label')}</h3>
              <div className="text-5xl font-black text-red-500 tracking-[0.5em] ml-[0.5em]">{generatedCode}</div>
              <p className="text-gray-500 mt-2 text-sm">{t('lobby.share_code_hint')}</p>
            </div>

            <div className="border-t border-gray-700 pt-6">
              <h4 className="text-xl font-bold mb-4">{t('lobby.invite_friends')}</h4>
              <div className="space-y-3">
                {onlineFriends.map(friend => (
                  <div key={friend.id} className="flex justify-between items-center bg-gray-900 p-4 rounded-xl border border-gray-700">
                    <span className="font-bold">{friend.pseudo}</span>
                    <button className="bg-red-600 hover:bg-red-700 text-xs font-black px-4 py-2 rounded-lg transition-colors">{t('lobby.invite')}</button>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex gap-4 pt-6">
              <button onClick={() => setView('menu')} className="flex-1 bg-gray-700 py-3 rounded-xl font-bold">{t('lobby.cancel')}</button>
              <button onClick={handleStartGame} className="flex-1 bg-green-600 hover:bg-green-700 py-3 rounded-xl font-bold transition-colors">
                {t('lobby.start_game')}
              </button>
            </div>
          </div>
        )}

        {view === 'join' && (
          <div className="flex flex-col space-y-8 items-center">
            <h3 className="text-3xl font-black text-white">{t('lobby.join_title')}</h3>
            <div className="w-full">
              <label className="block text-gray-400 text-center mb-4 font-bold uppercase tracking-widest">{t('lobby.code_input_label')}</label>
              <input type="text" maxLength={5} value={joinCode} onChange={(e) => setJoinCode(e.target.value.toUpperCase())} className="w-full bg-gray-900 border-2 border-red-500/50 focus:border-red-500 rounded-2xl py-6 text-5xl text-center font-black tracking-[0.3em] outline-none transition-all text-white" placeholder="_____" />
            </div>
            <div className="flex gap-4 w-full">
              <button onClick={() => setView('menu')} className="flex-1 bg-gray-700 py-4 rounded-xl font-bold">{t('lobby.back')}</button>
              <button disabled={joinCode.length < 5} onClick={handleStartGame} className={`flex-1 py-4 rounded-xl font-bold transition-all ${joinCode.length === 5 ? 'bg-red-600 hover:bg-red-700 shadow-[0_0_20px_rgba(220,38,38,0.4)]' : 'bg-gray-800 text-gray-600 cursor-not-allowed'}`}>
                {t('lobby.join')}
              </button>
            </div>
          </div>
        )}

      </div>
    </main>
  );
};

export default Lobby;
