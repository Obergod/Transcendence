import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const Profile = ({ onLogout }: { onLogout: () => void }) => {
  const navigate = useNavigate();

  const [pseudo, setPseudo] = useState("Impots.gouv.fr");
  const [isEditing, setIsEditing] = useState(false);

  const { t } = useTranslation();

  // 1. NOTRE FAUSSE BASE DE DONNÉES D'AMIS
  const [friends] = useState([
    { id: 1, pseudo: "Rillettes de conde", isOnline: true },
    { id: 2, pseudo: "Obergod", isOnline: false },
    { id: 3, pseudo: "Arthur", isOnline: true },
    { id: 4, pseudo: "Jackysuma", isOnline: false },
  ]);

  const handleLogout = () => {
    onLogout();
    navigate("/");
  };

  return (
    <main className="flex-1 flex flex-col items-center pt-32 px-4 pb-12 w-full">

      {/* Le conteneur principal qui va gérer les 2 colonnes sur grand écran (lg:flex-row) */}
      <div className="w-full max-w-6xl flex flex-col lg:flex-row gap-8">

        {/* ==========================================
            COLONNE GAUCHE : PARAMÈTRES DU PROFIL
            ========================================== */}
        <div className="flex-1 bg-gray-800/50 border-2 border-gray-700 rounded-3xl p-8 backdrop-blur-md shadow-2xl">
          <h2 className="text-4xl font-black text-center mb-10 tracking-tight">
            {t('profile.title')} <span className="text-red-500"> {t('profile.title_highlight')}</span>
          </h2>

          <div className="flex flex-col md:flex-row items-center md:items-start gap-12">

            {/* PHOTO DE PROFIL */}
            <div className="relative group">
              <div className="w-40 h-40 bg-gray-700 rounded-full border-4 border-red-600 overflow-hidden shadow-[0_0_20px_rgba(220,38,38,0.3)]">
                <svg className="w-full h-full text-gray-500 mt-4" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
                </svg>
              </div>
              <button className="absolute bottom-1 right-1 bg-red-600 p-2 rounded-full hover:bg-red-700 transition-colors border-2 border-gray-900">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                </svg>
              </button>
            </div>

            {/* INFORMATIONS */}
            <div className="flex-1 w-full space-y-6">
              <div>
                <label className="block text-gray-400 text-sm font-bold uppercase tracking-wider mb-2"> {t('profile.pseudo')}</label>
                <div className="flex gap-3">
                  <input
                    type="text"
                    value={pseudo}
                    disabled={!isEditing}
                    onChange={(e) => setPseudo(e.target.value)}
                    className={`flex-1 bg-gray-900 border ${isEditing ? 'border-red-500 ring-1 ring-red-500' : 'border-gray-600'} rounded-xl px-4 py-3 focus:outline-none transition-all`}
                  />
                  <button
                    onClick={() => setIsEditing(!isEditing)}
                    className="bg-gray-700 hover:bg-gray-600 px-6 rounded-xl font-bold transition-colors"
                  >
                    {isEditing ?  t('profile.save') :  t('profile.edit')}
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-gray-400 text-sm font-bold uppercase tracking-wider mb-2"> {t('profile.email')}</label>
                <input type="text" disabled value="maati@student.42.fr" className="w-full bg-gray-900/50 border border-gray-700 rounded-xl px-4 py-3 text-gray-500 cursor-not-allowed" />
              </div>

              {/* ACTIONS DANGEREUSES */}
              <div className="pt-6 border-t border-gray-700 mt-8 flex flex-col gap-4">
                <button onClick={handleLogout} className="w-full bg-transparent border-2 border-gray-600 hover:border-red-600 hover:text-red-500 text-gray-400 font-bold py-3 rounded-xl transition-all">
                   {t('profile.logout')}
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* ==========================================
            COLONNE DROITE : LISTE D'AMIS
            ========================================== */}
        <div className="w-full lg:w-1/3 bg-gray-800/50 border-2 border-gray-700 rounded-3xl p-6 backdrop-blur-md shadow-2xl flex flex-col">

          <div className="flex justify-between items-center mb-6">
            <h3 className="text-2xl font-bold text-white"> {t('profile.friends')}</h3>
            {/* Bouton pour ajouter un ami (visuel) */}
            <button className="text-red-500 hover:text-red-400 transition-colors p-2 bg-gray-900 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
              </svg>
            </button>
          </div>

          {/* LA LISTE */}
          <div className="flex-1 space-y-3 overflow-y-auto pr-2">

            {/* 2. LA BOUCLE MAP : On parcourt notre liste d'amis pour générer l'affichage */}
            {friends.map((friend) => (
              <div key={friend.id} className="flex items-center justify-between p-3 bg-gray-900/80 rounded-xl border border-gray-700 hover:border-gray-500 transition-colors">

                <div className="flex items-center space-x-4">
                  {/* Avatar de l'ami */}
                  <div className="relative">
                    <div className="w-10 h-10 bg-gray-700 rounded-full overflow-hidden">
                      <svg className="w-full h-full text-gray-500 mt-1" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
                      </svg>
                    </div>
                    {/* Le point de statut (Vert = En ligne, Gris = Hors ligne) */}
                    <div className={`absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-gray-900 ${friend.isOnline ? 'bg-green-500' : 'bg-gray-500'}`}></div>
                  </div>

                  {/* Nom de l'ami */}
                  <span className={`font-bold ${friend.isOnline ? 'text-white' : 'text-gray-500'}`}>
                    {friend.pseudo}
                  </span>
                </div>

                {/* MODIFICATION ICI : Remplacement du <button> par <Link> */}
                <Link
                  to={`/chat/${friend.pseudo}`}
                  className="text-gray-400 hover:text-red-500 transition-colors p-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                  </svg>
                </Link>

              </div>
            ))}

          </div>

        </div>

      </div>
    </main>
  );
};

export default Profile;
