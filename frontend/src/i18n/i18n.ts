import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import fr from './fr.json';
import en from './en.json';
import ru from './ru.json';

i18n
	.use(initReactI18next)
	.init({
		resources: {
			fr: { translation: fr},
			en: { translation: en},
			ru: { translation: ru},
		},
		lng: 'fr',	// default language
		fallbackLng: 'en',	// if missing key, use english
	});

export default i18n;
