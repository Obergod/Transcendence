import { useEffect, useState, useRef } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const Game = () => {
  const [isGameLoaded, setIsGameLoaded] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const { t } = useTranslation();

  useEffect(() => {
    let script: HTMLScriptElement | null = null;

    script = document.createElement('script');
    script.src = '/wasm_exec.js';
    script.onload = () => {
      const go = new (window as any).Go();
      WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject)
        .then((result) => {
          setIsGameLoaded(true);
          go.run(result.instance);
          // Après le lancement, déplacer le canvas
          const moveCanvas = () => {
            const canvas = document.querySelector('canvas');
            if (canvas && containerRef.current && !containerRef.current.contains(canvas)) {
              // Styler le canvas pour qu'il prenne toute la place
              canvas.style.width = '100%';
              canvas.style.height = '100%';
              canvas.style.display = 'block';
              containerRef.current.appendChild(canvas);
            }
          };
          // Attendre un court instant que le canvas soit créé
          setTimeout(moveCanvas, 50);
        })
        .catch((err) => {
          console.error("Erreur de chargement du Wasm:", err);
          setError(t('game.error_wasm'));
        });
    };
    script.onerror = () => {
      setError(t('game.error_script'));
    };
    document.body.appendChild(script);

    return () => {
      if (script && document.body.contains(script)) {
        document.body.removeChild(script);
      }
    };
  }, []);

  return (
    <main className="flex-1 flex flex-col items-center justify-center pt-24 pb-8 w-full">
      <div className="flex justify-between items-end w-full max-w-4xl mb-4 px-4">
        <h2 className="text-3xl font-black text-red-500 tracking-widest">{t('game.in_progress')}</h2>
        <div className="text-gray-400 font-mono">Score: 00000</div>
      </div>

      <div
        id="game-container"
        ref={containerRef}
        className="w-full max-w-4xl aspect-video bg-black border-4 border-gray-700 rounded-xl shadow-[0_0_30px_rgba(220,38,38,0.1)] flex items-center justify-center relative overflow-hidden"
      >
        {!isGameLoaded && !error && (
          <div className="text-center animate-pulse z-10">
            <p className="text-gray-500 font-mono text-lg mb-2">{t('game.loading')}</p>
            <div className="w-8 h-8 border-4 border-red-500 border-t-transparent rounded-full animate-spin mx-auto mt-4"></div>
          </div>
        )}
        {error && (
          <div className="text-center z-10 text-red-500">
            <p>{t('game.error_prefix')} {error}</p>
          </div>
        )}
      </div>

      <Link to="/" className="mt-8 px-6 py-3 bg-transparent border-2 border-gray-600 text-gray-400 hover:border-red-600 hover:text-red-500 font-bold rounded-lg transition-all z-10">
        {t('game.quit')}
      </Link>
    </main>
  );
};

export default Game;
