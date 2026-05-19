import { useState } from 'react';
import { useTranslation } from 'react-i18next';

// On définit les options que la modale attend
interface AuthModalProps {
  onClose: () => void;
  onLoginSuccess: () => void;
}

const AuthModal = ({ onClose, onLoginSuccess }: AuthModalProps) => {
  // Un "State" pour savoir si on affiche le formulaire de Connexion ou d'Inscription
  const [isRegistering, setIsRegistering] = useState(false);

  // Les "States" pour stocker ce que l'utilisateur tape
  const [pseudo, setPseudo] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const { t } = useTranslation();

  // Pour afficher des erreurs
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); // Empêche le rechargement de la page
    setError(null);     // On efface les anciennes erreurs à chaque nouvelle tentative

    // 1. SCÉNARIO : INSCRIPTION (REGISTER)
    if (isRegistering) {
      if (password !== confirmPassword) {
        setError(t('auth.error_password_mismatch'));
        return;
      }
      if (password.length < 6) {
        setError(t('auth.error_password_short'));
        return;
      }

      try {
        // Envoi de la requête au Backend Go
        const response = await fetch('/api/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ pseudo, email, password })
        });
        const data = await response.json();

        if (!response.ok) {
          throw new Error(data.error || data.message || t('auth.error_register'));
        }
        console.log("Succès :", data.message);
        onLoginSuccess();
        onClose();
      } catch (err: any) {
        setError(err.message);
      }

    // 2. SCÉNARIO : CONNEXION (LOGIN)
    } else {
      try {
        const response = await fetch('/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email, password })
        });
        const data = await response.json();

        if (!response.ok) {
          throw new Error(data.error || data.message || t('auth.error_login'));
        }

        // Plus tard, quand ton backend Go enverra un JWT, tu le stockeras ici :

        console.log("Connexion réussie ! Bienvenue", data.pseudo);
        onLoginSuccess();
        onClose();

      } catch (err: any) {
        setError(err.message);
      }
    }
  };

  return (
    // L'arrière-plan noir flouté qui recouvre tout le site
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">

      {/* La boîte de la modale */}
      <div className="bg-gray-900 border-2 border-gray-700 rounded-3xl w-full max-w-md overflow-hidden shadow-[0_0_50px_rgba(220,38,38,0.15)] relative">

        {/* Bouton Fermer */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-500 hover:text-white transition-colors"
        >
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>

        <div className="p-8">
          <h2 className="text-3xl font-black text-white text-center mb-8 tracking-widest">
            {isRegistering ? t('auth.register_title') : t('auth.login_title')}
          </h2>

          {error && (
            <div className="bg-red-500/20 border border-red-500 text-red-400 p-3 rounded-lg text-sm text-center mb-6 font-bold">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">

            {/* Le champ Pseudo (Uniquement pour l'inscription) */}
            {isRegistering && (
              <div>
                <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">{t('auth.pseudo')}</label>
                <input
                  type="text" required value={pseudo} onChange={(e) => setPseudo(e.target.value)}
                  className="w-full bg-gray-800 border border-gray-700 focus:border-red-500 rounded-xl px-4 py-3 text-white outline-none transition-colors"
                  placeholder="Maati_42"
                />
              </div>
            )}

            <div>
              <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">{t('auth.email')}</label>
              <input
                type="email" required value={email} onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-gray-800 border border-gray-700 focus:border-red-500 rounded-xl px-4 py-3 text-white outline-none transition-colors"
                placeholder="survivor@42.fr"
              />
            </div>

            <div>
              <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">{t('auth.password')}</label>
              <input
                type="password" required value={password} onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-gray-800 border border-gray-700 focus:border-red-500 rounded-xl px-4 py-3 text-white outline-none transition-colors"
                placeholder="••••••••"
              />
            </div>

            {/* Le champ Confirmation (Uniquement pour l'inscription) */}
            {isRegistering && (
              <div>
                <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">{t('auth.confirm_password')}</label>
                <input
                  type="password" required value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full bg-gray-800 border border-gray-700 focus:border-red-500 rounded-xl px-4 py-3 text-white outline-none transition-colors"
                  placeholder="••••••••"
                />
              </div>
            )}

            <button type="submit" className="w-full bg-red-600 hover:bg-red-700 text-white font-black py-4 rounded-xl mt-4 transition-all tracking-widest">
              {isRegistering ? t('auth.submit_register') : t('auth.submit_login')}
            </button>
          </form>

          {/* Séparateur OU */}
          <div className="my-6 flex items-center">
            <div className="flex-1 border-t border-gray-700"></div>
            <span className="px-4 text-gray-500 text-xs font-bold uppercase">{t('auth.or')}</span>
            <div className="flex-1 border-t border-gray-700"></div>
          </div>

          {/* Bouton de Mock Google */}
          <button
            type="button"
            onClick={() => {
              console.log("Mock Login via Google");
              onLoginSuccess();
              onClose();
            }}
            className="w-full bg-white text-gray-900 hover:bg-gray-200 font-bold py-3 rounded-xl transition-colors flex items-center justify-center space-x-2"
          >
            <svg className="w-5 h-5" viewBox="0 0 24 24">
              <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
            </svg>
            <span>{t('auth.google')}</span>
          </button>

          {/* Le bouton pour basculer entre Login et Sign Up (déjà existant) */}
          <div className="mt-6 text-center">
            <button
              onClick={() => { setIsRegistering(!isRegistering); setError(null); }}
              className="text-gray-500 hover:text-white text-sm font-bold transition-colors"
            >
              {isRegistering ? t('auth.already_account') : t('auth.no_account')}
            </button>
          </div>

        </div>
      </div>
    </div>
  );
};

export default AuthModal;
