import { I18n } from "i18n-js";
import en from "../locales/en.json";
import be from "../locales/be.json";
import ru from "../locales/ru.json";

const i18n = new I18n();

const loadLocales = () => {
    i18n.enableFallback = true;
    i18n.defaultLocale = "en";
    i18n.translations = {
        en: en,
        be: be,
        ru: ru
    };

    const language = window.Telegram.WebApp.initDataUnsafe.user.language_code;
    i18n.locale = language ? language : "en";
};

loadLocales();
export { i18n };