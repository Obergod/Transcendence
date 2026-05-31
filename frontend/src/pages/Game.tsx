import { useEffect, useState, useRef } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const Game = () => {
  const [isGameLoaded, setIsGameLoaded] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const { t } = useTranslation();

  useEffect(() => {
    // 1. NETTOYAGE EXTRÊME : On détruit tous les canvas existants (Fix du double écran)
    document.querySelectorAll('canvas').forEach(c => c.remove());

    let script = document.querySelector('script[src="/wasm_exec.js"]') as HTMLScriptElement;
    if (!script) {
        script = document.createElement('script');
        script.src = '/wasm_exec.js';
        document.body.appendChild(script);
    }

    // On force le rechargement de la fonction onload
    script.onload = () => {
      // Sécurité si on a quitté la page très vite
      if (!document.querySelector('#game-container')) return;

      const go = new (window as any).Go();
      WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject)
        .then((result) => {
          setIsGameLoaded(true);
          go.run(result.instance);

          const moveCanvas = () => {
            const canvases = document.querySelectorAll('canvas');
            if (canvases.length > 0 && containerRef.current) {
               // On prend le dernier canvas généré
               const canvas = canvases[canvases.length - 1];
               canvas.style.position = 'relative';
               canvas.style.width = '100%';
               canvas.style.height = '100%';
               containerRef.current.appendChild(canvas);

               // On détruit les clones s'il y en a eu !
               canvases.forEach((c, index) => {
                 if (index !== canvases.length - 1) c.remove();
               });
            }
          };
          setTimeout(moveCanvas, 100);
        })
        .catch((err) => {
          console.error("Erreur de chargement du Wasm:", err);
          setError(t('game.error_wasm'));
        });
    };

    script.onerror = () => setError(t('game.error_script'));

    // NETTOYAGE QUAND ON QUITTE LA PAGE
    return () => {
      document.querySelectorAll('canvas').forEach(c => c.remove());
    };
  }, [t]);

  return (
    <main className="flex-1 flex flex-col items-center justify-center pt-24 pb-8 w-full">
      <div className="flex justify-between items-end w-full max-w-4xl mb-4 px-4">
        <h2 className="text-3xl font-black text-red-500 tracking-widest">{t('game.in_progress')}</h2>
        <div className="text-gray-400 font-mono">Score: 00000</div>
      </div>

      <div
        id="game-container"
        ref={containerRef}
        className="w-full max-w-4xl aspect-video bg-black border-4 border-gray-700 rounded-xl flex items-center justify-center overflow-hidden"
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
