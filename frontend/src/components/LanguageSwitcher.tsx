import { useTranslation } from 'react-i18next';

const LanguageSwitcher = () => {
	const { i18n } = useTranslation();

	return (
	<div className="flex space-x-2">
		<button
			onClick={() => i18n.changeLanguage('fr')}
		className={i18n.language === 'fr'
			? 'text-white border-b-2 border-red-500 font-bold'
			: 'text-gray-400 hover:text-white transition-colors'}
		>
			FR
		</button>
		<button
			onClick={() => i18n.changeLanguage('en')}
		className={i18n.language === 'en'
			? 'text-white border-b-2 border-red-500 font-bold'
			: 'text-gray-400 hover:text-white transition-colors'}
		>
			EN
		</button>
		<button
			onClick={() => i18n.changeLanguage('ru')}
		className={i18n.language === 'ru'
			? 'text-white border-b-2 border-red-500 font-bold'
			: 'text-gray-400 hover:text-white transition-colors'}
		>
			RU
		</button>
	</div>
	);
};

export default LanguageSwitcher;
