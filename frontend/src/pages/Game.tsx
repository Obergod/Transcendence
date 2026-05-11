import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

const Game = () => {
  const [isGameLoaded, setIsGameLoaded] = useState(false);

  useEffect(() => {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject)
      .then((result) => {
        setIsGameLoaded(true);
        go.run(result.instance);
      })
      .catch((err) => console.error("Erreur de chargement du Wasm:", err));
  }, []);

  return (
    <main className="flex-1 flex flex-col items-center justify-center pt-24 pb-8 w-full">
      <div className="flex justify-between items-end w-full max-w-4xl mb-4 px-4">
        <h2 className="text-3xl font-black text-red-500 tracking-widest">SURVIE EN COURS...</h2>
        <div className="text-gray-400 font-mono">Score: 00000</div>
      </div>

      <div id="game-container" className="w-full max-w-4xl aspect-video bg-black border-4 border-gray-700 rounded-xl shadow-[0_0_30px_rgba(220,38,38,0.1)] flex items-center justify-center relative overflow-hidden">
        {!isGameLoaded && (
          <div className="text-center animate-pulse z-10">
            <p className="text-gray-500 font-mono text-lg mb-2">Chargement du moteur Ebitengine...</p>
            <div className="w-8 h-8 border-4 border-red-500 border-t-transparent rounded-full animate-spin mx-auto mt-4"></div>
          </div>
        )}
      </div>

      <Link to="/" className="mt-8 px-6 py-3 bg-transparent border-2 border-gray-600 text-gray-400 hover:border-red-600 hover:text-red-500 font-bold rounded-lg transition-all z-10">
        Abandonner la partie
      </Link>
    </main>
  );
};

export default Game;