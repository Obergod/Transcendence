import { useState, useEffect, useCallback } from 'react';
import { useUser } from '../context/UserContext';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const Profile = ({ onLogout }: { onLogout: () => void }) => {
  const { user, login, logout, onlineUsers = [] } = useUser() as any;
  const navigate = useNavigate();
  const { t } = useTranslation();

  const [pseudo, setPseudo] = useState(user?.pseudo || "");
  const [email, setEmail] = useState(user?.email || "");
  const [avatarUrl, setAvatarUrl] = useState(user?.avatarUrl || "");
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [message, setMessage] = useState<{ text: string, type: 'success' | 'error' } | null>(null);
  const [stats, setStats] = useState({ best_score: 0, best_duration: 0, total_duration: 0 });
  const [friendships, setFriendships] = useState<any[]>([]);
  const [targetUsername, setTargetUsername] = useState("");
  const [socialMsg, setSocialMsg] = useState<{ text: string, type: 'success' | 'error' } | null>(null);

  useEffect(() => {
    if (!user) navigate('/');
  }, [user, navigate]);

  useEffect(() => {
    const fetchStats = async () => {
      const token = localStorage.getItem('jwt_token');
      if (!token) return;
      try {
        const res = await fetch("http://localhost:8081/api/user/stats", {
          headers: { "Authorization": `Bearer ${token}` }
        });
        if (res.ok) {
          const data = await res.json();
          setStats(data);
        }
      } catch (err) {}
    };
    if (user) fetchStats();
  }, [user]);

  const formatStatsTime = (seconds: number) => {
    if (seconds === 0) return "0s";
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    if (h > 0) return `${h}h ${m}m ${s}s`;
    if (m > 0) return `${m}m ${s}s`;
    return `${s}s`;
  };

  const validateImage = (file: File): Promise<string | null> => {
    return new Promise((resolve) => {
      const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
      if (!allowedTypes.includes(file.type)) {
        resolve("Format non autorisé.");
        return;
      }
      if (file.size > 2 * 1024 * 1024) {
        resolve("Image trop lourde. Maximum 2 MB.");
        return;
      }
      const img = new Image();
      img.onload = () => {
        URL.revokeObjectURL(img.src);
        if (img.width > 2000 || img.height > 2000) {
          resolve("Image trop grande. Maximum 2000x2000.");
        } else {
          resolve(null);
        }
      };
      img.onerror = () => resolve("Impossible de lire l'image.");
      img.src = URL.createObjectURL(file);
    });
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    const token = localStorage.getItem('jwt_token');

    const formData = new FormData();
    formData.append("username", pseudo);
    formData.append("email", email);
    if (avatarFile) formData.append("avatar", avatarFile);

    try {
      const response = await fetch("http://localhost:8081/api/user/update", {
        method: "PUT",
        headers: { "Authorization": `Bearer ${token}` },
        body: formData,
      });
      const data = await response.json();

      if (response.ok) {
        setMessage({ text: "Profil mis à jour !", type: 'success' });
        setAvatarUrl(data.user.avatarUrl);
        login({ ...user!, pseudo: data.user.username, email: data.user.email, avatarUrl: data.user.avatarUrl });
        setAvatarFile(null);
      } else {
        setMessage({ text: data.error || "Erreur de mise à jour", type: 'error' });
      }
    } catch (error) {
      setMessage({ text: "Erreur serveur", type: 'error' });
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('jwt_token');
    logout();
    onLogout();
    navigate('/');
  };

  const fetchFriends = useCallback(async () => {
    const token = localStorage.getItem('jwt_token');
    if (!token) return;
    try {
      const response = await fetch("http://localhost:8081/api/friends/list", {
        headers: { "Authorization": `Bearer ${token}` }
      });
      const data = await response.json();
      if (response.ok) setFriendships(data.data);
    } catch (error) {}
  }, []);

  useEffect(() => {
    if (user) fetchFriends();
  }, [user, fetchFriends]);

  const sendFriendRequest = async (e: React.FormEvent) => {
    e.preventDefault();
    setSocialMsg(null);
    const token = localStorage.getItem('jwt_token');
    try {
      const response = await fetch("http://localhost:8081/api/friends/request", {
        method: "POST",
        headers: { "Content-Type": "application/json", "Authorization": `Bearer ${token}` },
        body: JSON.stringify({ username: targetUsername }),
      });
      const data = await response.json();
      if (response.ok) {
        setSocialMsg({ text: data.message, type: 'success' });
        setTargetUsername("");
        fetchFriends();
      } else {
        setSocialMsg({ text: data.error, type: 'error' });
      }
    } catch (error) {
      setSocialMsg({ text: "Erreur réseau", type: 'error' });
    }
  };

  const respondToRequest = async (friendshipId: number, action: 'accept' | 'reject') => {
    const token = localStorage.getItem('jwt_token');
    try {
      await fetch("http://localhost:8081/api/friends/respond", {
        method: "PUT",
        headers: { "Content-Type": "application/json", "Authorization": `Bearer ${token}` },
        body: JSON.stringify({ friendship_id: friendshipId, action }),
      });
      fetchFriends();
    } catch (error) {}
  };

  if (!user) return null;

  const pendingRequests = friendships.filter(f => f.Status === 'pending' && f.FriendID === user.id);
  const acceptedFriends = friendships.filter(f => f.Status === 'accepted');
  const displayAvatar = avatarFile ? URL.createObjectURL(avatarFile) : avatarUrl;

  return (
    <div className="min-h-screen pt-32 pb-12 px-8 grid grid-cols-1 lg:grid-cols-2 gap-8 max-w-7xl mx-auto">

      {/* COLONNE GAUCHE : MON PROFIL */}
      <div className="bg-[#0f1423] border border-gray-700 p-8 rounded-2xl shadow-2xl h-fit">
        <h2 className="text-3xl font-black text-white mb-8 text-center uppercase tracking-widest">{t('profile.my_profile')}</h2>

        {message && (
          <div className={`p-4 mb-6 rounded-xl text-center font-bold ${message.type === 'success' ? 'bg-green-900/50 text-green-400 border-green-500' : 'bg-red-900/50 text-red-400 border-red-500'} border`}>
            {message.text}
          </div>
        )}

        <div className="grid grid-cols-3 gap-2 bg-[#1a2035] p-4 rounded-xl border border-gray-700 text-center mb-8 shadow-inner">
          <div>
            <div className="text-[9px] font-bold text-gray-400 uppercase tracking-wider">{t('profile.best_score')}</div>
            <div className="text-lg font-black text-yellow-500 mt-1">{stats.best_score}</div>
          </div>
          <div>
            <div className="text-[9px] font-bold text-gray-400 uppercase tracking-wider">{t('profile.best_time')}</div>
            <div className="text-lg font-black text-green-400 mt-1">{formatStatsTime(stats.best_duration)}</div>
          </div>
          <div>
            <div className="text-[9px] font-bold text-gray-400 uppercase tracking-wider">{t('profile.total_time')}</div>
            <div className="text-lg font-black text-blue-400 mt-1">{formatStatsTime(stats.total_duration)}</div>
          </div>
        </div>

        <div className="flex justify-center mb-8">
          <img src={displayAvatar} alt="Avatar" className="w-32 h-32 rounded-full object-cover border-4 border-red-600 shadow-[0_0_15px_rgba(220,38,38,0.5)]" />
        </div>

        <form onSubmit={handleUpdate} className="space-y-6">
          <div>
            <label className="block text-gray-400 text-xs font-bold uppercase mb-2">{t('profile.pseudo')}</label>
            <input
              type="text"
              maxLength={15}
              value={pseudo}
              onChange={(e) => setPseudo(e.target.value)}
              className="w-full bg-[#1a2035] border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-red-500 outline-none transition-colors"
            />
          </div>
          <div>
            <label className="block text-gray-400 text-xs font-bold uppercase mb-2">{t('profile.email')}</label>
            <input
              type="email"
              maxLength={50}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-[#1a2035] border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-red-500 outline-none transition-colors"
            />
          </div>

          <div>
            <label className="block text-gray-400 text-xs font-bold uppercase mb-2">{t('profile.avatar')}</label>
            <input
              type="file"
              accept="image/*"
              onChange={async (e) =>{
                const file = e.target.files?.[0] || null;
                if (file) {
                  const error = await validateImage(file);
                  if (error) {
                    setMessage({ text: error, type: 'error'});
                    e.target.value = '';
                    return;
                  }
                }
                setAvatarFile(file);
              }}
              className="w-full bg-[#1a2035] border border-gray-700 rounded-xl px-4 py-3 text-white focus:border-red-500 outline-none transition-colors file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-bold file:bg-gray-700 file:text-white hover:file:bg-gray-600 cursor-pointer"
            />
          </div>

          <button type="submit" className="w-full bg-blue-600 hover:bg-blue-500 text-white font-black py-4 rounded-xl mt-4 uppercase tracking-widest shadow-lg">{t('profile.save')}</button>
        </form>

        <button onClick={handleLogout} className="w-full mt-6 bg-red-950/50 border border-red-500 hover:bg-red-600 text-red-400 hover:text-white font-bold py-3 rounded-xl transition-colors">{t('profile.logout')}</button>
      </div>

      {/* COLONNE DROITE : SOCIAL & AMIS */}
      <div className="bg-[#0f1423] border border-gray-700 p-8 rounded-2xl shadow-2xl h-fit flex flex-col gap-8">
        <h2 className="text-3xl font-black text-white text-center uppercase tracking-widest">{t('profile.social')}</h2>

        <form onSubmit={sendFriendRequest} className="flex flex-col gap-2">
          <label className="text-gray-400 text-xs font-bold uppercase">{t('profile.add_player')}</label>
          <div className="flex gap-2">
            <input
              type="text"
              placeholder={t('profile.friend_placeholder')}
              value={targetUsername}
              onChange={(e) => setTargetUsername(e.target.value)}
              className="flex-1 bg-[#1a2035] border border-gray-700 rounded-xl px-4 text-white focus:border-red-500 outline-none"
              required
            />
            <button type="submit" className="bg-green-600 hover:bg-green-500 text-white font-bold px-6 py-3 rounded-xl transition-colors">{t('profile.add_btn')}</button>
          </div>
          {socialMsg && <span className={`text-sm ${socialMsg.type === 'success' ? 'text-green-400' : 'text-red-400'}`}>{socialMsg.text}</span>}
        </form>

        <hr className="border-gray-700" />

        {pendingRequests.length > 0 && (
          <div>
            <h3 className="text-xl font-bold text-yellow-500 mb-4">{t('profile.pending_requests')}</h3>
            <div className="space-y-3">
              {pendingRequests.map(req => (
                <div key={req.ID} className="flex items-center justify-between bg-[#1a2035] p-3 rounded-xl border border-gray-700">
                  <div className="flex items-center gap-3">
                    <img src={req.User.AvatarURL} alt="avatar" className="w-10 h-10 rounded-full object-cover" />
                    <span className="font-bold">{req.User.Username}</span>
                  </div>
                  <div className="flex gap-2">
                    <button onClick={() => respondToRequest(req.ID, 'accept')} className="bg-green-600/20 text-green-500 hover:bg-green-600 hover:text-white p-2 rounded-lg transition-colors">✓</button>
                    <button onClick={() => respondToRequest(req.ID, 'reject')} className="bg-red-600/20 text-red-500 hover:bg-red-600 hover:text-white p-2 rounded-lg transition-colors">✗</button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        <div>
          <h3 className="text-xl font-bold text-white mb-4">{t('profile.my_friends')}</h3>
          {acceptedFriends.length === 0 ? (
            <p className="text-gray-500 italic">{t('profile.no_friends')}</p>
          ) : (
            <div className="space-y-3">
              {acceptedFriends.map(f => {
                const isSender = f.UserID === user.id;
                const friendData = isSender ? f.Friend : f.User;
                const isOnline = onlineUsers.includes(friendData.ID);

                return (
                  <div key={f.ID} className="flex items-center justify-between bg-[#1a2035] p-3 rounded-xl border border-gray-700 hover:border-gray-500 transition-colors cursor-pointer">
                    <div className="flex items-center gap-4">
                      <div className="relative">
                        <img src={friendData.AvatarURL} alt="avatar" className="w-12 h-12 rounded-full object-cover border border-gray-600" />
                        <span className={`absolute bottom-0 right-0 w-3 h-3 border-2 border-[#1a2035] rounded-full ${isOnline ? 'bg-green-500' : 'bg-gray-500'}`}></span>
                      </div>
                      <span className="font-bold text-lg">{friendData.Username}</span>
                    </div>
                    <button
                      onClick={() => navigate(`/chat/${friendData.ID}/${friendData.Username}`)}
                      className="text-gray-400 hover:text-white px-4 py-2 bg-gray-800 hover:bg-red-600 transition-colors rounded-lg"
                    >
                      {t('profile.chat_btn')}
                    </button>
                  </div>
                )
              })}
            </div>
          )}
        </div>

      </div>
    </div>
  );
};

export default Profile;