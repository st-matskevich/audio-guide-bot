import { I18n } from "i18n-js";
import en from "../locales/en.json";
import be from "../locales/be.json";
import ru from "../locales/ru.json";
import { getTelegramLanguage } from "./telegram";

const DEFAULT_LANGUAGE = "en";
const i18n = new I18n();

const loadLocales = () => {
    i18n.enableFallback = true;
    i18n.defaultLocale = DEFAULT_LANGUAGE;
    i18n.translations = {
        en: en,
        be: be,
        ru: ru
    };

    i18n.locale = getTelegramLanguage();
};

loadLocales();
export { i18n };