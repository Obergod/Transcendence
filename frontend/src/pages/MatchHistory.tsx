import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const MatchHistory = () => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<'history' | 'leaderboard'>('history');

  const [myHistory, setMyHistory] = useState<any[]>([]);
  const [leaderboard, setLeaderboard] = useState<any[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      const token = localStorage.getItem('jwt_token');
      if (!token) return;

      try {
        const histRes = await fetch("https://localhost:8081/api/match/history", {
          headers: { "Authorization": `Bearer ${token}` }
        });
        if (histRes.ok) {
          const histData = await histRes.json();
          setMyHistory(histData.data || []);
        }

        const leadRes = await fetch("https://localhost:8081/api/match/leaderboard", {
          headers: { "Authorization": `Bearer ${token}` }
        });
        if (leadRes.ok) {
          const leadData = await leadRes.json();
          setLeaderboard(leadData.data || []);
        }
      } catch (err) {
        console.error("Erreur de récupération des matchs", err);
      }
    };

    fetchData();
  }, []);

  const formatTime = (totalSeconds: number) => {
    const m = Math.floor(totalSeconds / 60).toString().padStart(2, '0');
    const s = (totalSeconds % 60).toString().padStart(2, '0');
    return `${m}:${s}`;
  };

  return (
    <main className="flex-1 flex flex-col items-center pt-24 px-4 pb-12 w-full">
      <div className="w-full max-w-4xl bg-gray-800/50 border-2 border-gray-700 rounded-3xl p-8 backdrop-blur-md shadow-2xl">

        <div className="flex justify-between items-center mb-8 border-b border-gray-700 pb-6">
          <h2 className="text-4xl font-black tracking-tight text-white uppercase">
            {t('history.data_center_pre')} <span className="text-red-500">{t('history.data_center_highlight')}</span>
          </h2>
          <Link to="/" className="text-gray-400 hover:text-white font-bold transition-colors">
            {t('history.back')}
          </Link>
        </div>

        {/* ONGLETS */}
        <div className="flex space-x-4 mb-8">
          <button
            onClick={() => setActiveTab('history')}
            className={`px-6 py-3 rounded-xl font-bold uppercase tracking-widest transition-all ${activeTab === 'history' ? 'bg-red-600 text-white shadow-[0_0_15px_rgba(220,38,38,0.5)]' : 'bg-[#1a2035] text-gray-400 hover:text-white border border-gray-700'}`}
          >
            {t('history.tab_history')}
          </button>
          <button
            onClick={() => setActiveTab('leaderboard')}
            className={`px-6 py-3 rounded-xl font-bold uppercase tracking-widest transition-all ${activeTab === 'leaderboard' ? 'bg-yellow-600 text-white shadow-[0_0_15px_rgba(202,138,4,0.5)]' : 'bg-[#1a2035] text-gray-400 hover:text-white border border-gray-700'}`}
          >
            {t('history.tab_leaderboard')}
          </button>
        </div>

        {/* CONTENU DES ONGLETS */}
        {activeTab === 'history' && (
          <div>
            <div className="grid grid-cols-2 text-gray-400 font-bold uppercase tracking-wider mb-4 px-4 text-sm">
              <span>{t('history.col_mission_date')}</span>
              <span className="text-right">{t('history.col_time_survived')}</span>
            </div>
            <div className="space-y-4">
              {myHistory.length === 0 && <p className="text-center text-gray-500 italic py-8">{t('history.no_missions')}</p>}
              {myHistory.map((match, i) => (
                <div key={i} className="grid grid-cols-2 items-center p-4 rounded-xl border bg-gray-900/40 border-gray-700 transition-colors hover:bg-gray-700/50">
                  <span className="text-gray-400 font-semibold">{new Date(match.CreatedAt).toLocaleDateString()}</span>
                  <span className="text-right text-white font-mono text-xl font-black tracking-widest">{formatTime(match.Duration)}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'leaderboard' && (
          <div>
            <div className="grid grid-cols-3 text-yellow-500/70 font-bold uppercase tracking-wider mb-4 px-4 text-sm">
              <span>{t('history.col_rank')}</span>
              <span className="text-center">{t('history.col_player')}</span>
              <span className="text-right">{t('history.col_record')}</span>
            </div>
            <div className="space-y-4">
              {leaderboard.length === 0 && <p className="text-center text-gray-500 italic py-8">{t('history.no_scores')}</p>}
              {leaderboard.map((match, index) => (
                <div key={index} className={`grid grid-cols-3 items-center p-4 rounded-xl border ${index === 0 ? 'bg-yellow-900/30 border-yellow-600 shadow-[0_0_20px_rgba(202,138,4,0.2)]' : 'bg-gray-900/40 border-gray-700'} transition-colors hover:bg-gray-700/50`}>
                  <div className="flex items-center space-x-3">
                    <span className={`text-2xl font-black ${index === 0 ? 'text-yellow-500' : index === 1 ? 'text-gray-400' : index === 2 ? 'text-amber-700' : 'text-gray-600'}`}>#{index + 1}</span>
                  </div>
                  <span className={`text-center font-black uppercase tracking-widest ${index === 0 ? 'text-yellow-500' : 'text-white'}`}>{match.username}</span>
                  <span className="text-right text-white font-mono text-xl font-black tracking-widest">{formatTime(match.duration)}</span>
                </div>
              ))}
            </div>
          </div>
        )}

      </div>
    </main>
  );
};

export default MatchHistory;