import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const MatchHistory = () => {
  // Voici exactement le format de données (JSON) que le serveur Go t'enverra un jour !
  const fakeBackendData = [
    { id: 101, date: "06 Mai 2026", duration: "15:42", kills: 1420, result: "Victoire" },
    { id: 102, date: "05 Mai 2026", duration: "08:12", kills: 312, result: "Défaite" },
    { id: 103, date: "04 Mai 2026", duration: "20:00", kills: 2890, result: "Victoire" },
    { id: 104, date: "02 Mai 2026", duration: "12:30", kills: 840, result: "Défaite" },
  ];

  const [matches] = useState(fakeBackendData);

  const { t } = useTranslation();

  return (
    <main className="flex-1 flex flex-col items-center pt-24 px-4 pb-12 w-full">
      <div className="w-full max-w-4xl bg-gray-800/50 border-2 border-gray-700 rounded-3xl p-8 backdrop-blur-md shadow-2xl">

        <div className="flex justify-between items-center mb-10 border-b border-gray-700 pb-6">
          <h2 className="text-4xl font-black tracking-tight text-white">
            {t('history.title')} <span className="text-red-500">{t('history.title_highlight')}</span>
          </h2>
          <Link to="/" className="text-gray-400 hover:text-white font-bold transition-colors">
            {t('history.back')}
          </Link>
        </div>

        {/* L'en-tête du tableau */}
        <div className="grid grid-cols-4 text-gray-400 font-bold uppercase tracking-wider mb-4 px-4 text-sm md:text-base">
          <span>{t('history.col_date')}</span>
          <span className="text-center">{t('history.col_duration')}</span>
          <span className="text-center">{t('history.col_kills')}</span>
          <span className="text-right">{t('history.col_result')}</span>
        </div>

        {/* La liste des matchs générée automatiquement avec .map() */}
        <div className="space-y-4">
          {matches.map((match) => (
            <div
              key={match.id}
              className={`grid grid-cols-4 items-center p-4 rounded-xl border ${
                match.result === t('history.victory')
                  ? 'bg-green-900/20 border-green-800'
                  : 'bg-red-900/20 border-red-800'
              } transition-colors hover:bg-gray-700/50`}
            >
              <span className="text-white font-semibold">{match.date}</span>

              <span className="text-gray-300 font-mono text-center">{match.duration}</span>

              <span className="text-center text-yellow-500 font-black">{match.kills}</span>

              <span className={`text-right font-black uppercase tracking-widest ${
                match.result === t('history.victory') ? 'text-green-500' : 'text-red-500'
              }`}>
                {match.result}
              </span>
            </div>
          ))}
        </div>

      </div>
    </main>
  );
};

export default MatchHistory;
