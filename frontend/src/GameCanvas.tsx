import { useEffect } from 'react';

// Déclarer le constructeur Go qui vient de wasm_exec.js
declare global {
  const Go: any;
}

export default function GameCanvas() {
  useEffect(() => {
    // Éviter de recharger le script plusieurs fois
    if (!document.querySelector('script[src="/wasm/wasm_exec.js"]')) {
      const script = document.createElement('script');
      script.src = '/wasm/wasm_exec.js';
      script.onload = () => {
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch('/wasm/main.wasm'), go.importObject)
          .then((result) => go.run(result.instance))
          .catch(console.error);
      };
      document.body.appendChild(script);
    } else {
      const go = new Go();
      WebAssembly.instantiateStreaming(fetch('/wasm/main.wasm'), go.importObject)
        .then((result) => go.run(result.instance))
        .catch(console.error);
    }
  }, []);

  return (
    <div>
      <h3>Zone de jeu</h3>
      <div id="ebiten-game"></div>
    </div>
  );
}