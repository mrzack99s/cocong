import ConfirmDialog from "../components/ConfirmDialog";
import Navbar from "../components/Navbar";
import AppPropertiesProvider from "../utils/AppProperties";
import {
  createDOMRenderer,
  FluentProvider,
  GriffelRenderer,
  SSRProvider,
  RendererProvider,
  webLightTheme,
} from "@fluentui/react-components";
import type { AppProps } from "next/app";
import Head from "next/head";
import { useRouter } from "next/router";
import "office-ui-fabric-core/dist/css/fabric.min.css";
import { CookiesProvider } from "react-cookie";

type EnhancedAppProps = AppProps & { renderer?: GriffelRenderer };

function MyApp({ Component, pageProps, renderer }: EnhancedAppProps) {
  const router = useRouter();
  return (
    // 👇 Accepts a renderer from <Document /> or creates a default one
    //    Also triggers rehydration a client

    <RendererProvider renderer={renderer || createDOMRenderer()}>
      <SSRProvider>
        <CookiesProvider>
          <FluentProvider theme={webLightTheme}>
            <AppPropertiesProvider>
              <Head>
                <meta name="viewport" content="width=1376" />
              </Head>
              {router.pathname.startsWith("/admin") && <Navbar />}
              <Component {...pageProps} />
            </AppPropertiesProvider>
          </FluentProvider>
        </CookiesProvider>
      </SSRProvider>
    </RendererProvider>
  );
}

export default MyApp;
