import { useEffect, useState, useRef } from 'react';
import { Link } from 'react-router-dom';

const Game = () => {
  const [isGameLoaded, setIsGameLoaded] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let script: HTMLScriptElement | null = null;

    // Éviter de recharger le jeu s'il est déjà là
    if (document.querySelector('canvas')) {
       setIsGameLoaded(true);
       return;
    }

    script = document.createElement('script');
    script.src = '/wasm_exec.js';
    script.onload = () => {
      const go = new (window as any).Go();
      WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject)
        .then((result) => {
          setIsGameLoaded(true);
          go.run(result.instance);

          // FORCER LE CANVAS DANS LE CONTENEUR
          const moveCanvas = () => {
            const canvas = document.querySelector('canvas');
            if (canvas && containerRef.current) {
               canvas.style.position = 'relative';
               canvas.style.width = '100%';
               canvas.style.height = '100%';
               containerRef.current.appendChild(canvas);
            }
          };
          setTimeout(moveCanvas, 100);
        })
        .catch((err) => {
          console.error("Erreur Wasm:", err);
          setError("Impossible de charger le jeu.");
        });
    };
    script.onerror = () => setError("Impossible de charger wasm_exec.js");
    document.body.appendChild(script);

    // NETTOYAGE QUAND ON QUITTE LA PAGE
    return () => {
      if (script && document.body.contains(script)) {
        document.body.removeChild(script);
      }
      const canvas = document.querySelector('canvas');
      if (canvas) {
          canvas.remove();
      }
    };
  }, []);

  return (
    <main className="flex-1 flex flex-col items-center justify-center pt-24 pb-8 w-full">
      <div className="flex justify-between items-end w-full max-w-4xl mb-4 px-4">
        <h2 className="text-3xl font-black text-red-500 tracking-widest">SURVIE EN COURS...</h2>
        <div className="text-gray-400 font-mono">Score: 00000</div>
      </div>

      <div
        id="game-container"
        ref={containerRef}
        className="w-full max-w-4xl aspect-video bg-black border-4 border-gray-700 rounded-xl flex items-center justify-center overflow-hidden"
      >
        {!isGameLoaded && !error && (
          <div className="text-center animate-pulse z-10 text-white">
            <p>Chargement du moteur...</p>
          </div>
        )}
      </div>

      <Link to="/" className="mt-8 px-6 py-3 border-2 border-gray-600 text-gray-400 font-bold rounded-lg z-10">
        Abandonner la partie
      </Link>
    </main>
  );
};

export default Game;