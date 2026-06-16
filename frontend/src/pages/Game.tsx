import { useEffect, useState, useRef } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useUser } from '../context/UserContext';

const Game = () => {
  const [isGameLoaded, setIsGameLoaded] = useState(false);
  const [isGameOver, setIsGameOver] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const { t } = useTranslation();
  const { refreshLevel } = useUser();

  // onrécupère le mode choisi dans le Lobby (1 ou 2 joueurs). Par défaut 1.
  const location = useLocation();
  const gameMode = location.state?.mode || 1;

  useEffect(() => {
    (window as any).onGameover = async (durationInSeconds: number, score: number) => {
      setIsGameOver(true);
      const token = localStorage.getItem('jwt_token');
      try {
        await fetch("/api/match/save", {
          method: "POST",
          headers: { "Content-Type": "application/json", "Authorization": `Bearer ${token}` },
          body: JSON.stringify({ duration: durationInSeconds, score: score }),
        });
        await refreshLevel();
      } catch (err) {
        console.error("Erreur de sauvegarde du score", err);
      }
    };

    document.querySelectorAll('canvas').forEach(c => c.remove());

    const initWasm = () => {
      if (!document.querySelector('#game-container')) return;

      const go = new (window as any).Go();

      // LA MAGIE EST ICI : On passe les variables à Go juste avant de le lancer !
      (window as any).gameMode = gameMode;
      (window as any).timerLabel = t('game.timer_label');
      (window as any).scoreLabel = t('game.score_label');

      WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject)
        .then((result) => {
          setIsGameLoaded(true);
          go.run(result.instance);

          const moveCanvas = () => {
            const canvases = document.querySelectorAll('canvas');
            if (canvases.length > 0 && containerRef.current) {
               const canvas = canvases[canvases.length - 1];
               canvas.style.position = 'relative';
               canvas.style.width = '100%';
               canvas.style.height = '100%';
               containerRef.current.appendChild(canvas);

               canvases.forEach((c, index) => {
                 if (index !== canvases.length - 1) c.remove();
               });
            }
          };
          setTimeout(moveCanvas, 100);
        })
        .catch((err) => {
          console.error("Erreur Wasm:", err);
          setError(t('game.error_wasm'));
        });
    };

    if ((window as any).Go) {
        initWasm();
    } else {
        let script = document.createElement('script');
        script.src = '/wasm_exec.js';
        script.onload = initWasm;
        script.onerror = () => setError(t('game.error_script'));
        document.body.appendChild(script);
    }

    return () => {
      document.querySelectorAll('canvas').forEach(c => c.remove());
      delete (window as any).onGameover;
    };
  }, [t, gameMode]);

  const handleRetry = () => {
    setIsGameOver(false);
    if ((window as any).restartGame) {
        (window as any).restartGame();
    }
  };

  return (
    <main className="flex-1 flex flex-col items-center justify-center pt-24 pb-8 w-full">

      {/* HEADER DU JEU (Avec les IDs pour que Go puisse injecter le texte en temps réel) */}
      <div className="flex justify-between items-end w-full max-w-7xl mb-4 px-4">
        <h2 id="game-timer" className="text-3xl font-black text-red-500 tracking-widest uppercase">
            {t('game.timer_init')}
        </h2>
        <div id="game-score" className="text-gray-400 font-mono text-xl">{t('game.score_init')}</div>
      </div>

      {/* CONTENEUR DU JEU */}
      <div
        id="game-container"
        ref={containerRef}
        className="relative w-full max-w-7xl aspect-video bg-black border-4 border-gray-700 rounded-xl flex items-center justify-center overflow-hidden shadow-2xl"
      >
        {!isGameLoaded && !error && (
          <div className="text-center animate-pulse z-10">
            <p className="text-gray-500 font-mono text-lg mb-2">{t('game.loading')}</p>
            <div className="w-8 h-8 border-4 border-red-500 border-t-transparent rounded-full animate-spin mx-auto mt-4"></div>
          </div>
        )}

        {error && (
          <div className="text-center z-10 text-red-500 font-bold">
            <p>{t('game.error_prefix')} {error}</p>
          </div>
        )}
      </div>

      {/* BOUTON QUITTER */}
      <Link to="/" className="mt-8 px-8 py-3 bg-transparent border-2 border-gray-600 text-gray-500 hover:border-red-600 hover:text-red-500 font-black uppercase tracking-widest rounded-xl transition-all z-10">
        {t('game.quit')}
      </Link>

      {/*POPUP GAME OVER*/}
      {isGameOver && (
        <div className="fixed inset-0 bg-black/90 flex items-center justify-center z-[100] backdrop-blur-md p-4">
          <div className="bg-[#0f1423] border border-gray-700 rounded-2xl w-full max-w-md shadow-[0_0_50px_rgba(230,0,0,0.3)] relative overflow-hidden flex flex-col">

            <div className="h-2 w-full bg-[#e60000]"></div>

            <div className="p-10 text-center">
              <h2 className="text-4xl font-black text-white mb-2 uppercase tracking-tighter leading-tight">
                {t('game.over_title')}
              </h2>
              <p className="text-gray-500 font-bold text-xs uppercase tracking-widest mb-8">
                {t('game.over_subtitle')}
              </p>

              <div className="space-y-4">
                <button
                  onClick={handleRetry}
                  className="w-full bg-[#e60000] hover:bg-red-700 text-white font-black py-5 rounded-xl transition-all uppercase tracking-widest shadow-[0_0_20px_rgba(230,0,0,0.4)] hover:scale-[1.02] active:scale-95"
                >
                  {t('game.retry')}
                </button>

                <Link
                  to="/"
                  className="w-full inline-block bg-gray-800 hover:bg-gray-700 text-gray-300 font-black py-4 rounded-xl transition-all uppercase tracking-widest text-sm"
                >
                  {t('game.abandon')}
                </Link>
              </div>
            </div>

            <div className="bg-[#1a2035] p-4 text-center border-t border-gray-800">
                <span className="text-gray-600 text-[10px] font-black uppercase tracking-widest">{t('game.footer')}</span>
            </div>
          </div>
        </div>
      )}

    </main>
  );
};

export default Game;
