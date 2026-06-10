import { useTranslation } from 'react-i18next';

const LanguageSwitcher = () => {
	const { i18n } = useTranslation();

	const flags: Record<string, string> = {
		fr: '/flags/fr.svg',
		en: '/flags/gb.svg',
		ru: '/flags/ru.svg',
	};

	const langs = ['fr', 'en', 'ru'];

	return (
		<div className ="flex space-x-3">
			{langs.map((lng) => (
				<button
					key={lng}
					onClick={() => i18n.changeLanguage(lng)}
					className={`font-black text-4xl bg-clip-text text-transparent bg-cover bg-center transition-opacity ${
						i18n.language === lng ? 'opacity-100' : 'opacity-40 hover:opacity-70'
					}`}
					style={{ backgroundImage: `url(${flags[lng]})` }}
				>
					{lng.toUpperCase()}
				</button>
			))}
		</div>
	);
};

export default LanguageSwitcher;
