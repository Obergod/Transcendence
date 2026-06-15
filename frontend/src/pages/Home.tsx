import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

interface HomeProps {
  isLoggedIn: boolean;
  onLoginClick: () => void;
}

const Home = ({ isLoggedIn, onLoginClick }: HomeProps) => {
  const { t } = useTranslation();

  return (
    <main className="flex-1 flex flex-col items-center justify-center space-y-12">
      <h1 className="text-6xl md:text-8xl font-black text-transparent bg-clip-text bg-gradient-to-b from-red-400 to-red-800 tracking-widest drop-shadow-lg text-center">
        {t('home.title_line1')}<br />{t('home.title_line2')}
      </h1>

      <div className="flex flex-col space-y-6 w-72">
        {isLoggedIn ? (
          <>
            <Link to="/play" className="bg-red-600 hover:bg-red-700 border-2 border-red-500 py-4 rounded-xl text-2xl font-black uppercase tracking-widest transition-all text-center shadow-[0_0_20px_rgba(220,38,38,0.4)] text-white">
              {t('home.play')}
            </Link>
            <Link to="/history" className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-yellow-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center">
              {t('home.match_history')}
            </Link>
          </>
        ) : (
          <button onClick={onLoginClick} className="bg-red-600 hover:bg-red-700 border-2 border-red-500 py-4 rounded-xl text-4xl font-bold uppercase tracking-widest transition-all text-white shadow-[0_0_20px_rgba(220,38,38,0.5)]">
            {t('home.play')}
          </button>
        )}
      </div>
    </main>
  );
};

export default Home;
