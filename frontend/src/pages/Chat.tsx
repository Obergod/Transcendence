import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';

const Chat = () => {
  // On récupère le nom de l'ami depuis l'URL (ex: /chat/Dracula_xX)
  const { friendName } = useParams();

  // On adapte notre fausse base de données pour ce chat privé
  const [messages, setMessages] = useState([
    { id: 1, author: friendName, content: `Salut, tu es chaud pour un multi ?`, isMe: false, time: "11:20" },
    { id: 2, author: "Impots.gouv.fr", content: "Grave, je finis un truc et je t'invite !", isMe: true, time: "11:22" },
  ]);

  const [currentInput, setCurrentInput] = useState("");

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (currentInput.trim() === "") return;

    const newMessage = {
      id: messages.length + 1,
      author: "Impots.gouv.fr",
      content: currentInput,
      isMe: true,
      time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    };

    setMessages([...messages, newMessage]);
    setCurrentInput("");
  };

  return (
    <main className="flex-1 flex flex-col items-center pt-24 px-4 pb-12 w-full">
      <div className="w-full max-w-4xl bg-gray-900 border-2 border-gray-700 rounded-3xl overflow-hidden shadow-2xl flex flex-col h-[70vh]">

        {/* En-tête : Affiche dynamiquement le nom de l'ami */}
        <div className="bg-gray-800 p-4 border-b border-gray-700 flex justify-between items-center">
          <div className="flex items-center space-x-4">

            <div className="relative">
              <div className="w-10 h-10 bg-gray-700 rounded-full overflow-hidden">
                 <svg className="w-full h-full text-gray-500 mt-1" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" /></svg>
              </div>
              {/* Petite pastille verte de statut */}
              <div className="absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-gray-900 bg-green-500"></div>
            </div>

            <div>
              <h2 className="text-xl font-black text-white tracking-widest">{friendName?.toUpperCase()}</h2>
              <span className="text-green-400 text-xs font-bold">En ligne</span>
            </div>

          </div>

          {/* Le bouton fermer te ramène sur ton profil, là où est ta liste d'amis ! */}
          <Link to="/profile" className="text-gray-400 hover:text-white font-bold transition-colors">
            Fermer le chat
          </Link>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-6 bg-gray-900/50">
          {messages.map((msg) => (
            <div key={msg.id} className={`flex flex-col ${msg.isMe ? 'items-end' : 'items-start'}`}>
              <span className="text-gray-500 text-xs font-bold mb-1 ml-1">
                {msg.author} • {msg.time}
              </span>
              <div className={`px-5 py-3 rounded-2xl max-w-[80%] ${msg.isMe ? 'bg-red-600 text-white rounded-br-none' : 'bg-gray-800 border border-gray-700 text-gray-200 rounded-bl-none'}`}>
                {msg.content}
              </div>
            </div>
          ))}
        </div>

        <div className="bg-gray-800 p-4 border-t border-gray-700">
          <form onSubmit={handleSendMessage} className="flex gap-3">
            <input type="text" value={currentInput} onChange={(e) => setCurrentInput(e.target.value)} placeholder={`Envoyer un message à ${friendName}...`} className="flex-1 bg-gray-900 border border-gray-600 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-red-500 transition-colors" />
            <button type="submit" className="bg-red-600 hover:bg-red-700 text-white font-bold px-8 rounded-xl transition-all shadow-[0_0_15px_rgba(220,38,38,0.3)]">Envoyer</button>
          </form>
        </div>

      </div>
    </main>
  );
};

export default Chat;