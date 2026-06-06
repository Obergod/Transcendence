import { useState, useEffect, useRef } from 'react';
import { Link, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useUser } from '../context/UserContext';

const Chat = () => {
  const { friendId, friendName } = useParams();
  const { user, ws } = useUser();
  const { t } = useTranslation();

  const [messages, setMessages] = useState<any[]>([]);
  const [currentInput, setCurrentInput] = useState("");

  // NOUVEAU : État pour afficher l'alerte de spam
  const [isSpamming, setIsSpamming] = useState(false);
  // NOUVEAU : Référence pour stocker l'heure du dernier envoi sans re-render le composant
  const lastMessageTime = useRef<number>(0);

  const [friendAvatar, setFriendAvatar] = useState<string | null>(null);
  const [friendLevel, setFriendLevel] = useState<number | null>(null);

  useEffect(() => {
    const fetchHistory = async () => {
      if (!friendId) return;
      const token = localStorage.getItem('jwt_token');
      try {
        const res = await fetch(`http://localhost:8081/api/chat/history/${friendId}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
          const data = await res.json();
          setMessages(data.data || []);
          if (data.friend) {
            if (data.friend.avatarUrl) setFriendAvatar(data.friend.avatarUrl);
            if (data.friend.level) setFriendLevel(data.friend.level);
          }
        }
      } catch (err) {
        console.error("Erreur chargement historique", err);
      }
    };
    fetchHistory();
  }, [friendId]);

  useEffect(() => {
    if (!ws) return;

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.type === "chat") {
          const isRelevant =
            (data.sender_id === user?.id && data.target_id === Number(friendId)) ||
            (data.sender_id === Number(friendId) && data.target_id === user?.id);

          if (isRelevant) {
            setMessages((prev) => [...prev, data]);
          }
        }
      } catch (err) {
        console.error("Erreur lecture message WS", err);
      }
    };
  }, [ws, user, friendId]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();

    if (currentInput.trim() === "" || currentInput.length > 300 || !ws) return;

    const now = Date.now();
    const cooldownMs = 1000; // 1 sec

    if (now - lastMessageTime.current < cooldownMs) {
      setIsSpamming(true);
      // Fait disparaître le message d'erreur après 1.5s
      setTimeout(() => setIsSpamming(false), 1500);
      return;
    }
    // ------------------------------------

    const payload = {
      type: "chat",
      target_id: Number(friendId),
      content: currentInput
    };

    ws.send(JSON.stringify(payload));

    // Réinitialisation après envoi
    setCurrentInput("");
    setIsSpamming(false);
    lastMessageTime.current = now; // On met à jour l'heure du dernier message
  };

  return (
    <main className="flex-1 flex flex-col items-center pt-24 px-4 pb-12 w-full">
      <div className="w-full max-w-4xl bg-[#0f1423] border border-gray-700 rounded-3xl overflow-hidden shadow-2xl flex flex-col h-[70vh]">

        {/* En-tête */}
        <div className="bg-[#1a2035] p-4 border-b border-gray-700 flex justify-between items-center">
          <div className="flex items-center space-x-4">
            <div className="relative">
              <div className="w-10 h-10 bg-gray-700 rounded-full overflow-hidden flex items-center justify-center font-bold text-gray-400">
                 {friendAvatar ? (
                  <img src={friendAvatar} alt={friendName} className="w-full h-full object-cover" />
                 ) : (
                  friendName?.charAt(0).toUpperCase()
                 )}
              </div>
              <div className="absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-[#1a2035] bg-green-500 shadow-[0_0_10px_rgba(34,197,94,0.8)]"></div>
            </div>
            <div>
              <div className="flex items-baseline gap-2">
                <h2 className="text-xl font-black text-white tracking-widest">{friendName?.toUpperCase()}</h2>
                {friendLevel && (
                  <span className="text-red-500 text-sm font-bold">{t('level.short')} {friendLevel}</span>
                )}
              </div>
              <span className="text-green-400 text-xs font-bold">{t('chat.online')}</span>
            </div>
          </div>

          <Link to="/profile" className="text-gray-400 hover:text-white font-bold transition-colors bg-gray-800 px-4 py-2 rounded-lg">
            {t('chat.back_to_profile')}
          </Link>
        </div>

        {/* Zone des messages */}
        <div className="flex-1 overflow-y-auto p-6 space-y-4 bg-[#0a0d17]">
          {messages.length === 0 && (
            <div className="text-center text-gray-500 italic mt-10">{t('chat.no_messages')}</div>
          )}

          {messages.map((msg, index) => {
            const isMe = msg.sender_id === user?.id;
            return (
              <div key={index} className={`flex flex-col ${isMe ? 'items-end' : 'items-start'}`}>
                <span className="text-gray-600 text-[10px] font-bold mb-1 ml-1 uppercase tracking-wider">
                  {msg.sender_name}
                </span>
                <div className={`px-5 py-3 rounded-2xl max-w-[80%] shadow-md ${isMe ? 'bg-red-600 text-white rounded-br-none' : 'bg-[#1a2035] border border-gray-700 text-gray-200 rounded-bl-none'}`}>
                  {msg.content}
                </div>
              </div>
            );
          })}
        </div>

        {/* Input */}
        <div className="bg-[#1a2035] p-4 border-t border-gray-700">
          <form onSubmit={handleSendMessage} className="flex gap-3 items-start">
            <div className="flex-1 flex flex-col">
              <input
                type="text"
                maxLength={300}
                value={currentInput}
                onChange={(e) => setCurrentInput(e.target.value)}
                placeholder={t('chat.placeholder', { name: friendName })}
                className={`w-full bg-[#0a0d17] border rounded-xl px-4 py-3 text-white focus:outline-none transition-colors ${
                  isSpamming ? 'border-orange-500 bg-orange-900/20' : 'border-gray-600 focus:border-red-500'
                }`}
              />

              {/* Zone d'infos sous l'input : Alerte de spam + Compteur */}
              <div className="flex justify-between items-center mt-1 px-1">
                <div className="text-orange-500 text-xs font-bold h-4">
                  {isSpamming ? t('chat.spam_warning') : ""}
                </div>
                <div className={`text-right text-xs font-bold ${currentInput.length >= 300 ? 'text-red-500' : 'text-gray-500'}`}>
                  {currentInput.length}/300
                </div>
              </div>

            </div>
            <button
              type="submit"
              disabled={isSpamming}
              className={`font-black uppercase tracking-widest px-8 py-3 h-[50px] rounded-xl transition-all shadow-[0_0_15px_rgba(220,38,38,0.3)] ${
                isSpamming
                  ? 'bg-gray-600 text-gray-400 cursor-not-allowed shadow-none'
                  : 'bg-red-600 hover:bg-red-700 text-white'
              }`}
            >
              {t('chat.send')}
            </button>
          </form>
        </div>

      </div> {/* <-- C'est cette balise qui manquait ! */}
    </main>
  );
};

export default Chat;
