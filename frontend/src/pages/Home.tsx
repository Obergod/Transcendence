import { Link } from 'react-router-dom';

// On définit les informations que Home s'attend à recevoir depuis App.tsx
interface HomeProps {
  isLoggedIn: boolean;
  onLoginClick: () => void;
}

const Home = ({ isLoggedIn, onLoginClick }: HomeProps) => {
  return (
    <main className="flex-1 flex flex-col items-center justify-center space-y-12">
      <h1 className="text-6xl md:text-8xl font-black text-transparent bg-clip-text bg-gradient-to-b from-red-400 to-red-800 tracking-widest drop-shadow-lg text-center">
        42<br />SURVIVOR
      </h1>

      <div className="flex flex-col space-y-6 w-72">
        {/* AFFICHAGE CONDITIONNEL : Si connecté, on montre les boutons du jeu */}
        {isLoggedIn ? (
          <>
            <Link to="/play/solo" className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-red-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center">
              Solo
            </Link>
            <Link to="/play/multi" className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-red-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center">
              Multi
            </Link>
            <Link to="/history" className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-yellow-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center">
              Match History
            </Link>
          </>
        ) : (
          <button
            onClick={onLoginClick}
            className="bg-red-600 hover:bg-red-700 border-2 border-red-500 py-4 rounded-xl text-xl font-bold uppercase tracking-widest transition-all text-white shadow-[0_0_20px_rgba(220,38,38,0.5)]"
          >
            Login / Signin
          </button>
        )}
      </div>
    </main>
  );
};

export default Home;