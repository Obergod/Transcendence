import { Link } from 'react-router-dom';

const GameSolo = () => {
  return (
    <main className="flex-1 flex flex-col items-center justify-center pt-24 pb-8 w-full">

      <div className="flex justify-between items-end w-full max-w-4xl mb-4 px-4">
        <h2 className="text-3xl font-black text-red-500 tracking-widest">SURVIE EN COURS...</h2>
        <div className="text-gray-400 font-mono">Score: 00000</div>
      </div>

      {/* ==========================================
          LA BOÎTE DU JEU (C'est ici que le Wasm s'injectera !)
          On lui donne un ID "game-container" pour que le script Go le trouve facilement.
          ========================================== */}
      <div
        id="game-container"
        className="w-full max-w-4xl aspect-video bg-black border-4 border-gray-700 rounded-xl shadow-[0_0_30px_rgba(220,38,38,0.1)] flex items-center justify-center relative overflow-hidden"
      >
        <div className="text-center animate-pulse">
          <p className="text-gray-500 font-mono text-lg mb-2">En attente du moteur Ebitengine...</p>
          <p className="text-gray-600 font-mono text-sm">[Le fichier game.wasm sera chargé ici]</p>
        </div>
      </div>

      {/* Bouton pour quitter la partie */}
      <Link
        to="/"
        className="mt-8 px-6 py-3 bg-transparent border-2 border-gray-600 text-gray-400 hover:border-red-600 hover:text-red-500 font-bold rounded-lg transition-all"
      >
        Abandonner la partie
      </Link>

    </main>
  );
};

export default GameSolo;