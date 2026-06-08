import { useState } from 'react';
import { useUser } from '../context/UserContext';
import { useTranslation } from 'react-i18next';

interface AuthModalProps {
  onClose: () => void;
  onLoginSuccess: () => void;
}

const AuthModal = ({ onClose, onLoginSuccess }: AuthModalProps) => {
  const { login } = useUser();
  const { t } = useTranslation();

  const [isSignUp, setIsSignUp] = useState(true);
  const [pseudo, setPseudo] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);

    if (isSignUp && password !== confirmPassword) {
      setErrorMsg(t('auth.error_password_mismatch'));
      return;
    }

    const url = isSignUp ? "http://localhost:8081/api/signup" : "http://localhost:8081/api/signin";
    const payload = isSignUp
      ? { username: pseudo, email: email, password: password }
      : { login: email, password: password };

    try {
      const response = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      const text = await response.text();
      let data;

      try {
        data = JSON.parse(text);
      } catch (err) {
        setErrorMsg(t('auth.error_server_unexpected', { text }));
        return;
      }

      if (response.ok) {
        if (data.token) {
          localStorage.setItem('jwt_token', data.token);
        }

        if (data.user) {
          login({
            id: data.user.id,
            pseudo: data.user.username,
            email: data.user.email,
            avatarUrl: data.user.avatarUrl,
            status: data.user.status
          });
          onLoginSuccess();
        } else {
          setErrorMsg(t('auth.error_missing_user'));
        }
      } else {
        setErrorMsg(data.message || data.error || t('auth.error_login'));
      }

    } catch (error) {
      console.error("Fetch error:", error);
      setErrorMsg(t('auth.error_server_unreachable'));
    }
  };

  return (
    <div className="absolute inset-0 bg-black/80 flex items-center justify-center z-50 backdrop-blur-sm p-4">
      <div className="bg-[#0f1423] border border-gray-700 rounded-2xl w-full max-w-md shadow-2xl relative overflow-hidden flex flex-col">

        <button onClick={onClose} className="absolute top-4 right-4 text-gray-500 hover:text-white transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <div className="p-8 pb-6">
          <h2 className="text-3xl font-black text-center text-white mb-6 uppercase tracking-widest leading-tight">
            {isSignUp ? t('auth.register_title') : t('auth.login_title')}
          </h2>

          {errorMsg && (
            <div className="bg-red-950/50 border border-red-500 text-red-400 p-4 rounded-xl text-sm font-semibold text-center mb-6">
              {errorMsg}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {isSignUp && (
              <div>
                <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">{t('auth.pseudo')}</label>
                <input
                  type="text"
                  maxLength={15}
                  value={pseudo}
                  onChange={(e) => setPseudo(e.target.value)}
                  className="w-full bg-[#1a2035] border border-gray-700 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-red-500 transition-colors"
                  required={isSignUp}
                />
              </div>
            )}

            <div>
              <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">
                {isSignUp ? t('auth.email') : t('auth.pseudo_or_email')}
              </label>
              <input
                type="text"
                maxLength={50}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-[#1a2035] border border-gray-700 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-red-500 transition-colors"
                required
              />
            </div>

            <div>
              <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">{t('auth.password')}</label>
              <input
                type="password"
                maxLength={30}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-[#1a2035] border border-gray-700 rounded-xl px-4 py-3 text-white tracking-widest focus:outline-none focus:border-red-500 transition-colors"
                required
              />
            </div>

            {isSignUp && (
              <div>
                <label className="block text-gray-400 text-xs font-bold uppercase tracking-wider mb-2">{t('auth.confirm_password')}</label>
                <input
                  type="password"
                  maxLength={30}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full bg-[#1a2035] border border-gray-700 rounded-xl px-4 py-3 text-white tracking-widest focus:outline-none focus:border-red-500 transition-colors"
                  required={isSignUp}
                />
              </div>
            )}

            <button type="submit" className="w-full bg-[#e60000] hover:bg-red-700 text-white font-black py-4 rounded-xl mt-4 transition-colors uppercase tracking-widest shadow-[0_0_15px_rgba(220,38,38,0.3)]">
              {isSignUp ? t('auth.submit_register') : t('auth.login_btn')}
            </button>
          </form>
        </div>

        <div className="bg-[#1a2035] p-4 text-center mt-auto border-t border-gray-800">
          <button
            onClick={() => {
              setIsSignUp(!isSignUp);
              setErrorMsg(null);
            }}
            className="text-gray-400 hover:text-white font-bold text-sm transition-colors"
          >
            {isSignUp ? t('auth.already_account') : t('auth.no_account')}
          </button>
        </div>

      </div>
    </div>
  );
};

export default AuthModal;
