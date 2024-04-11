import { useCookies } from "react-cookie";
import {
  Spinner,
  Toast,
  ToastBody,
  ToastIntent,
  ToastTitle,
  Toaster,
  useId,
  useToastController,
  Image,
} from "@fluentui/react-components";
import axios, { AxiosInstance } from "axios";
import { useRouter } from "next/router";
import {
  useContext,
  createContext,
  useState,
  useEffect,
  useRef,
  ReactElement,
} from "react";
import { ApiConnector } from "./ApiConnector";
import { IProfile } from "../interfaces/Profile";
import ConfirmDialog from "../components/ConfirmDialog";
import { IConfirmDialog } from "../interfaces/Common";
type ContextProps = {
  closePageLoading: () => void;
  setPageLoading: (msg: string) => void;
  apiConnector: ApiConnector;
  profile: IProfile;
  setToast: (
    title: string,
    msg: ReactElement,
    intent?: ToastIntent,
    bgColor?: string
  ) => void;
  setConfirmDialog: (
    title: string,
    msg: ReactElement,
    onAccept: () => void,
    onReject: () => void
  ) => void;
};

export const AppPropertiesContext = createContext<ContextProps | null>(null);

interface AppProperties {
  children: any;
}

const AppPropertiesProvider: React.FC<AppProperties> = ({ children }) => {
  const [pageWaiting, setPageWaiting] = useState(false);
  const [pageWaitingMsg, setPageWaitingMsg] = useState("");
  const toasterId = useId("toaster");
  const { dispatchToast } = useToastController(toasterId);
  const [allEnv, setAllEnv] = useState({} as any);
  const [profile, setProfile] = useState<IProfile>({ Name: "", ID: "" });
  const [confirmDialogProps, setComfirmDialogProps] =
    useState<IConfirmDialog>();
  const [confirmOpen, setConfirmOpen] = useState(false);

  const setToast = (
    title: string,
    msg: ReactElement,
    intent?: ToastIntent,
    bgColor?: string
  ) => {
    dispatchToast(
      <Toast style={{ background: bgColor ? bgColor : "#fefefe" }}>
        <ToastTitle>{title}</ToastTitle>
        <ToastBody style={{ wordWrap: "break-word", width: "85%" }}>
          {msg}
        </ToastBody>
      </Toast>,
      {
        pauseOnHover: true,
        intent: intent ? intent : "info",
        position: "top-end",
        timeout: 2000,
      }
    );
  };

  const setConfirmDialog = (
    title: string,
    msg: ReactElement,
    onAccept: () => void,
    onReject: () => void
  ) => {
    setConfirmOpen(true);
    setComfirmDialogProps({
      message: msg,
      onAccept: onAccept,
      onReject: onReject,
      title: title,
    });
  };

  const setPageLoading = (msg: string) => {
    setPageWaitingMsg(msg);
    setPageWaiting(true);
  };

  const closePageLoading = () => {
    setPageWaiting(false);
    setPageWaitingMsg("");
  };
  const router = useRouter();
  const [cookies, setCookie, removeCookie] = useCookies([
    "api-token",
    "refresh-token",
  ]);
  const apiConnector = new ApiConnector(
    process.env.NEXT_PUBLIC_API_URL || "",
    cookies["api-token"],
    router
  );

  useEffect(() => {
    if (!cookies["api-token"]) {
      removeCookie("api-token");
      if (cookies["refresh-token"]) {
        apiConnector.instance
          .post("/op/refresh-token", {
            RefreshToken: cookies["refresh-token"],
          })
          .then((res) => {
            setCookie("api-token", res.data.AccessToken, {
              expires: new Date(res.data.AccessTokenExpired),
            });

            removeCookie("refresh-token");
            setCookie("refresh-token", res.data.RefreshToken, {
              expires: new Date(res.data.RefreshTokenExpired),
            });

            router.push("/admin");
          })
          .catch(() => {
            removeCookie("refresh-token");
            removeCookie("api-token");
            router.push("/");
          });
      } else {
        router.push("/");
      }
    } else {
      apiConnector.instance
        .get("/op/me")
        .then((res) => {
          setProfile(res.data);
          if (router.pathname === "/") {
            router.push("/admin");
          }
        })
        .catch((err) => {
          if (err.response.status == 401) {
            if (cookies["refresh-token"]) {
              apiConnector.instance
                .post("/op/refresh-token", {
                  RefreshToken: cookies["refresh-token"],
                })
                .then((res) => {
                  setCookie("api-token", res.data.AccessToken, {
                    expires: new Date(res.data.AccessTokenExpired),
                  });

                  removeCookie("refresh-token");
                  setCookie("refresh-token", res.data.RefreshToken, {
                    expires: new Date(res.data.RefreshTokenExpired),
                  });

                  router.push("/admin");
                })
                .catch(() => {
                  removeCookie("refresh-token");
                  removeCookie("api-token");
                  router.push("/");
                });
            }
          }
        });
    }
  }, [router.asPath]);

  return (
    <AppPropertiesContext.Provider
      value={{
        setPageLoading: setPageLoading,
        closePageLoading: closePageLoading,
        setToast: setToast,
        apiConnector: apiConnector,
        profile: profile,
        setConfirmDialog: setConfirmDialog,
      }}
    >
      {confirmDialogProps && (
        <ConfirmDialog
          message={confirmDialogProps.message}
          open={confirmOpen}
          onAccept={() => {
            confirmDialogProps.onAccept();
            setConfirmOpen(false);
          }}
          onReject={() => {
            confirmDialogProps.onReject();
            setConfirmOpen(false);
          }}
          title={confirmDialogProps.title}
        />
      )}

      {pageWaiting && (
        <div
          style={{
            backgroundColor: "rgba(255, 255, 255, 0.2)",
            backdropFilter: "blur(10px)",
            width: "100vw",
            height: "100vh",
            position: "fixed",
            zIndex: 9999,
            justifyContent: "center",
            justifyItems: "center",
            display: "flex",
          }}
        >
          <Spinner
            labelPosition="below"
            size="huge"
            label={pageWaitingMsg ? pageWaitingMsg : "Wait a minute"}
          />
        </div>
      )}

      <Toaster toasterId={toasterId} />
      {!pageWaiting && <>{children}</>}
    </AppPropertiesContext.Provider>
  );
};

export default AppPropertiesProvider;
export const usePageLoading = (): [() => void, (msg: string) => void] => {
  const appContext = useContext(AppPropertiesContext) as ContextProps;
  return [appContext.closePageLoading, appContext.setPageLoading];
};

export const useProfile = () => {
  try {
    const appContext = useContext(AppPropertiesContext) as ContextProps;
    return appContext.profile;
  } catch (err) {
    console.log(err);

    return {} as IProfile;
  }
};

export const useToast = () => {
  const appContext = useContext(AppPropertiesContext) as ContextProps;
  return appContext.setToast;
};

export const useApiConnector = () => {
  const appContext = useContext(AppPropertiesContext) as ContextProps;
  return appContext.apiConnector.instance;
};

export const useConfirmDialog = () => {
  const appContext = useContext(AppPropertiesContext) as ContextProps;
  return appContext.setConfirmDialog;
};
