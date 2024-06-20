import "../../styles/globals.css";
import type { AppProps } from "next/app";
import { Router } from "next/router";
import NProgress from "nprogress";
import "nprogress/nprogress.css";
import AuthProvider from "../context/auth.context";
import "react-datepicker/dist/react-datepicker.css";
import "../../styles/main.css";
import { registerLocale } from "react-datepicker";
import ru from "date-fns/locale/ru";
import store from '@store/index'
import { Provider } from "react-redux";

registerLocale("ru", ru);

Router.events.on("routeChangeStart", () => NProgress.start());
Router.events.on("routeChangeComplete", () => NProgress.done());
Router.events.on("routeChangeError", () => NProgress.done());
NProgress.configure({ showSpinner: false });

function MyApp({ Component, pageProps }: AppProps) {
  const env = process.env.NEXT_PUBLIC_ENV_NAME;
  return (
    <Provider store={store}>
      <AuthProvider>
      {env === "dev" && (
        <div className="absolute inset-x-0 flex items-center justify-center pointer-events-none">
          <div className="px-4 py-1 text-sm bg-emerald-400 font-medium text-white rounded-b-md">
            DEMO
          </div>
        </div>
      )}
      <Component {...pageProps} />
    </AuthProvider>
    </Provider>
  );
}

export default MyApp;
