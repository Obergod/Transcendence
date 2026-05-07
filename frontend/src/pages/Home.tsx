import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <main className="flex-1 flex flex-col items-center justify-center space-y-12">
      <h1 className="text-6xl md:text-8xl font-black text-transparent bg-clip-text bg-gradient-to-b from-red-400 to-red-800 tracking-widest drop-shadow-lg text-center">
        42<br />SURVIVOR
      </h1>

      <div className="flex flex-col space-y-6 w-72">
        <Link to="/play/solo" className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-red-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center">
          Solo
        </Link>
        <button className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-red-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all">
          Multi
        </button>
        <Link to="/history" className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-yellow-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center">
          Match History
        </Link>
      </div>
    </main>
  );
};

export default Home;









<Link
  to="/play/multi"
  className="bg-gray-800 hover:bg-gray-700 border-2 border-gray-600 hover:border-red-500 py-4 rounded-xl text-2xl font-bold uppercase tracking-widest transition-all text-center"
>
  Multi
</Link>