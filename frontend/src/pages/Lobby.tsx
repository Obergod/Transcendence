import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const Lobby = () => {
  const navigate = useNavigate();
  const { t } = useTranslation();

  const handleStartGame = (playersCount: number) => {
    navigate('/game', { state: { mode: playersCount } });
  };

  return (
    <main className="flex-1 flex flex-col items-center pt-32 px-4 pb-12 w-full">
      <div className="w-full max-w-2xl bg-gray-800/50 border-2 border-gray-700 rounded-3xl p-8 backdrop-blur-md shadow-2xl">
        <div className="flex flex-col space-y-6">
          <h2 className="text-4xl font-black text-center mb-8 tracking-tight">
            {t('lobby.title')}<span className="text-red-500">{t('lobby.title_highlight')}</span>
          </h2>

          <button
            onClick={() => handleStartGame(1)}
            className="bg-red-600 hover:bg-red-700 text-white font-black py-6 rounded-2xl text-2xl transition-all shadow-[0_0_20px_rgba(220,38,38,0.4)] uppercase tracking-widest hover:scale-[1.02] active:scale-95"
          >
            {t('lobby.solo_mode')}
          </button>

          <button
            onClick={() => handleStartGame(2)}
            className="bg-blue-600 hover:bg-blue-700 text-white font-black py-6 rounded-2xl text-2xl transition-all shadow-[0_0_20px_rgba(37,99,235,0.4)] uppercase tracking-widest hover:scale-[1.02] active:scale-95"
          >
            {t('lobby.coop_mode')}
          </button>

          <Link to="/" className="text-center text-gray-500 hover:text-white transition-colors pt-6 font-bold uppercase tracking-widest text-sm">
            {t('lobby.back_to_menu')}
          </Link>
        </div>
      </div>
    </main>
  );
};

export default Lobby;