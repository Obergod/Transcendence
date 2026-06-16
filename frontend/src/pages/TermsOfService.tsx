import { useTranslation } from 'react-i18next';

const TermsOfService = () => {
  const { t } = useTranslation();

  return (
    <div className="min-h-screen pt-32 pb-12 px-8 max-w-4xl mx-auto text-gray-300">
      <h1 className="text-4xl font-black text-white mb-8">{t('tos.title')}</h1>

      <section className="mb-8">
        <h2 className="text-2xl font-bold text-white mb-4">{t('tos.s1_title')}</h2>
        <p>{t('tos.s1_text')}</p>
      </section>

      <section className="mb-8">
        <h2 className="text-2xl font-bold text-white mb-4">{t('tos.s2_title')}</h2>
        <ul className="list-disc pl-6 space-y-2">
          <li>{t('tos.s2_li1')}</li>
          <li>{t('tos.s2_li2')}</li>
          <li>{t('tos.s2_li3')}</li>
        </ul>
      </section>

      <section className="mb-8">
        <h2 className="text-2xl font-bold text-white mb-4">{t('tos.s3_title')}</h2>
        <p>{t('tos.s3_text')}</p>
      </section>

      <section className="mb-8">
        <h2 className="text-2xl font-bold text-white mb-4">{t('tos.s4_title')}</h2>
        <p>{t('tos.s4_text')}</p>
      </section>

      <section className="mb-8">
        <h2 className="text-2xl font-bold text-white mb-4">{t('tos.s5_title')}</h2>
        <p>{t('tos.s5_text')}</p>
      </section>
    </div>
  );
};

export default TermsOfService;