import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useUser } from '../context/UserContext';

const Chat = () => {
  // On récupère l'ID et le nom depuis l'URL
  const { friendId, friendName } = useParams();
  const { user, ws } = useUser();
  const { t } = useTranslation();

  const [messages, setMessages] = useState<any[]>([]);
  const [currentInput, setCurrentInput] = useState("");

  // Écoute des messages en provenance de Go
  useEffect(() => {
    if (!ws) return;

    // On définit ce qu'il se passe quand un message arrive
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        // On vérifie que c'est un message de chat et qu'il correspond à la conversation actuelle
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

    // Note : Pour l'historique complet, il faudra créer une route Go "GET /api/chat/history/:friendId"
    // et faire un fetch() ici au chargement du composant !

  }, [ws, user, friendId]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (currentInput.trim() === "" || !ws) return;

    // On prépare le colis pour le Hub Go
    const payload = {
      type: "chat",
      target_id: Number(friendId),
      content: currentInput
    };

    // On envoie au serveur !
    ws.send(JSON.stringify(payload));
    setCurrentInput("");
  };

  return (
    <main className="flex-1 flex flex-col items-center pt-24 px-4 pb-12 w-full">
      <div className="w-full max-w-4xl bg-[#0f1423] border border-gray-700 rounded-3xl overflow-hidden shadow-2xl flex flex-col h-[70vh]">

        {/* En-tête */}
        <div className="bg-[#1a2035] p-4 border-b border-gray-700 flex justify-between items-center">
          <div className="flex items-center space-x-4">
            <div className="relative">
              <div className="w-10 h-10 bg-gray-700 rounded-full flex items-center justify-center font-bold text-gray-400">
                 {friendName?.charAt(0).toUpperCase()}
              </div>
              <div className="absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-[#1a2035] bg-green-500 shadow-[0_0_10px_rgba(34,197,94,0.8)]"></div>
            </div>
            <div>
              <h2 className="text-xl font-black text-white tracking-widest">{friendName?.toUpperCase()}</h2>
              <span className="text-green-400 text-xs font-bold">{t('chat.online')}</span>
            </div>
          </div>

          <Link to="/profile" className="text-gray-400 hover:text-white font-bold transition-colors bg-gray-800 px-4 py-2 rounded-lg">
            Retour au profil
          </Link>
        </div>

        {/* Zone des messages */}
        <div className="flex-1 overflow-y-auto p-6 space-y-4 bg-[#0a0d17]">
          {messages.length === 0 && (
            <div className="text-center text-gray-500 italic mt-10">Aucun message. Lance la discussion !</div>
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
          <form onSubmit={handleSendMessage} className="flex gap-3">
            <input
              type="text"
              value={currentInput}
              onChange={(e) => setCurrentInput(e.target.value)}
              placeholder={`Écrire à ${friendName}...`}
              className="flex-1 bg-[#0a0d17] border border-gray-600 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-red-500 transition-colors"
            />
            <button type="submit" className="bg-red-600 hover:bg-red-700 text-white font-black uppercase tracking-widest px-8 rounded-xl transition-all shadow-[0_0_15px_rgba(220,38,38,0.3)]">
              Envoyer
            </button>
          </form>
        </div>

      </div>
    </main>
  );
};

export default Chat;