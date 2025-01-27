import { themeColors } from "../constants/theme";
import { IMenuDrawer } from "../interfaces/Common";
import {
  Accordion,
  AccordionHeader,
  AccordionItem,
  AccordionPanel,
  AccordionToggleEventHandler,
  Avatar,
  Badge,
  Body1Strong,
  Body2,
  Button,
  Caption1,
  Dialog,
  DialogActions,
  DialogBody,
  DialogContent,
  DialogSurface,
  DialogTitle,
  DialogTrigger,
  Divider,
  DrawerBody,
  DrawerHeader,
  DrawerHeaderNavigation,
  DrawerHeaderTitle,
  Input,
  Link,
  MessageBar,
  MessageBarBody,
  MessageBarTitle,
  Popover,
  PopoverSurface,
  PopoverTrigger,
  Tag,
  Toolbar,
  ToolbarButton,
  ToolbarDivider,
  ToolbarToggleButton,
  Tooltip,
  makeStyles,
  shorthands,
} from "@fluentui/react-components";
import { DrawerOverlay } from "@fluentui/react-components/unstable";
import { useRouter } from "next/router";
import {
  Dispatch,
  SetStateAction,
  useState,
  useEffect,
  FormEvent,
} from "react";
import {
  TextUnderlineRegular,
  SettingsCogMultipleRegular,
  AppGenericFilled,
  HighlightFilled,
} from "@fluentui/react-icons";
import { useApiConnector, useProfile, useToast } from "../utils/AppProperties";
import { useCookies } from "react-cookie";

// Css in Typescript
const useStyles = makeStyles({
  navbar: {
    width: "100%",
    height: "60px",
    display: "flex",
    alignItems: "center",
    backgroundColor: themeColors.colorPrimary,
    // position: "fixed",
    color: "#fefefe",
    justifyContent: "space-between",
    // zIndex: 9999,
  },
  navbarLeft: {
    width: "auto",
    height: "60px",
    display: "flex",
    alignItems: "center",
    backgroundColor: themeColors.colorPrimary,
    color: "#fefefe",
    ...shorthands.padding("0", "0", "0", "1.5rem"),
  },
  navbarRight: {
    width: "auto",
    height: "60px",
    display: "flex",
    alignItems: "center",
    backgroundColor: themeColors.colorPrimary,
    color: "#fefefe",
    ...shorthands.padding("0", "1.5rem", "0", "0"),
  },
  btnDrawer: {
    ...shorthands.padding("0", "0", "0", "0"),
  },
  alertBtn: {
    color: "#fefefe",
    transitionDuration: "0.1s",
    "&:hover": {
      color: "red",
    },
    ...shorthands.margin("0", "1rem", "0", "0"),
  },
  avatar: {
    cursor: "pointer",
  },
  drawerToolbar: {
    justifyContent: "space-between",
  },
});

interface INavbar extends IMenuDrawer {}
const Navbar = () => {
  const styles = useStyles();

  const router = useRouter();
  const apiConnector = useApiConnector();
  const [cookies, setCookie, removeCookie] = useCookies([
    "api-token",
    "refresh-token",
  ]);

  const menus = [
    {
      name: "Monitor",
      icon: <AppGenericFilled />,
      link: "/admin",
      description: "Show log and monitor of your system",
    },
    {
      set_divider: true,
    },
    // {
    //   name: "Configure",
    //   icon: <SettingsCogMultipleRegular />,
    //   link: "/admin/configure",
    //   description: "System configure",
    // },
    // {
    //   set_divider: true,
    // },
    {
      name: "BW Profile Management",
      icon: <AppGenericFilled />,
      link: "/admin/bandwidth-profile-management",
      description: "To manage your bandwidth profiles",
    },
    {
      name: "Directory Management",
      icon: <AppGenericFilled />,
      link: "/admin/directory-management",
      description: "To manage your directories",
    },
    {
      name: "User Management",
      icon: <AppGenericFilled />,
      link: "/admin/user-management",
      description: "To manage your users",
    },
    {
      name: "Administrator Management",
      icon: <AppGenericFilled />,
      link: "/admin/administrator-management",
      description: "To manage your administrator users",
    },
  ];

  const profile = useProfile();
  const [openChangePassword, setOpenChangePassword] = useState(false);
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newAgainPassword, setNewAgainPassword] = useState("");
  const toast = useToast();

  const changePassword = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (newPassword === newAgainPassword) {
      apiConnector
        .post("/op/change-password", {
          CurrentPassword: currentPassword,
          NewPassword: newPassword,
        })
        .then(() => {
          toast("Success", <>Password changed</>, "success");
          setCurrentPassword("");
          setNewPassword("");
          setNewAgainPassword("");
          setOpenChangePassword(false);
        })
        .catch((err) => {
          toast(
            "Error",
            <>Somethig wrong! please check your password</>,
            "error"
          );
        });
    } else {
      toast("Warning", <>Your new password mismatch</>, "warning");
    }
  };

  return (
    <>
      <Dialog modalType="alert" open={openChangePassword}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>Change password</DialogTitle>
            <form onSubmit={changePassword}>
              <DialogContent>
                <div
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    Current Password{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "300px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={currentPassword}
                      onChange={(_, d) => {
                        setCurrentPassword(d.value);
                      }}
                      type="password"
                      required
                    />
                  </div>
                </div>
                <div
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    New Password{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "300px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={newPassword}
                      onChange={(_, d) => {
                        setNewPassword(d.value);
                      }}
                      type="password"
                      required
                    />
                  </div>
                </div>
                <div
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    New Password Again{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "300px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={newAgainPassword}
                      onChange={(_, d) => {
                        setNewAgainPassword(d.value);
                      }}
                      type="password"
                      required
                    />
                  </div>
                </div>
              </DialogContent>
              <DialogActions>
                <DialogTrigger disableButtonEnhancement>
                  <Button
                    appearance="secondary"
                    type="button"
                    onClick={() => {
                      setCurrentPassword("");
                      setNewPassword("");
                      setNewAgainPassword("");
                      setOpenChangePassword(false);
                    }}
                  >
                    Close
                  </Button>
                </DialogTrigger>
                <Button appearance="primary" type="submit">
                  Change
                </Button>
              </DialogActions>
            </form>
          </DialogBody>
        </DialogSurface>
      </Dialog>
      <div className={styles.navbar}>
        <div className={styles.navbarLeft}>
          <div>
            <Link
              href="/"
              style={{
                color: "#fff",
                fontSize: "20pt",
                textDecoration: "none",
              }}
              appearance="subtle"
            >
              CoCoNG
            </Link>
          </div>
        </div>
        <div className={styles.navbarRight}>
          <Popover>
            <PopoverTrigger disableButtonEnhancement>
              <Avatar
                size={40}
                name={profile.Name}
                color="neutral"
                className={styles.avatar}
              />
            </PopoverTrigger>

            <PopoverSurface style={{ width: "300px", padding: 0 }}>
              <div
                style={{
                  display: "flex",
                  justifyContent: "center",
                  marginTop: "2rem",
                }}
              >
                <Avatar size={96} name={profile.Name} color="brand" />
              </div>

              <div style={{ textAlign: "center", marginTop: "1rem" }}>
                <Body2>{profile.Name}</Body2>
              </div>
              <Divider style={{ margin: "1rem 0" }} alignContent="start">
                Actions
              </Divider>
              <div style={{ padding: "0.25rem 1rem 1rem 1rem" }}>
                <div>
                  <Button
                    appearance="secondary"
                    style={{ marginLeft: "0.25rem" }}
                    size="small"
                    onClick={() => {
                      setOpenChangePassword(true);
                    }}
                  >
                    Change password
                  </Button>
                  <Button
                    size="small"
                    style={{
                      marginLeft: "0.25rem",
                      background: themeColors.colorDanger,
                      color: "#fff",
                    }}
                    onClick={() => {
                      apiConnector
                        .delete("/op/logout")
                        .then(() => {
                          removeCookie("api-token");
                          removeCookie("refresh-token");
                          router.push("/");
                        })
                        .catch(() => {});
                    }}
                  >
                    Logout
                  </Button>
                </div>
              </div>
            </PopoverSurface>
          </Popover>
        </div>
      </div>

      <div>
        <Toolbar aria-label="with Tooltip" size="large">
          {menus.map((item) => (
            <>
              {!item.set_divider && (
                <>
                  <Tooltip
                    content={item.description!}
                    relationship="description"
                    withArrow
                  >
                    <ToolbarButton
                      vertical
                      appearance={
                        router.pathname === item.link ? "primary" : "subtle"
                      }
                      icon={item.icon}
                      style={{ margin: "1px" }}
                      onClick={() => {
                        if (router.pathname !== item.link) {
                          router.push(item.link!);
                        }
                      }}
                    >
                      {item.name}
                    </ToolbarButton>
                  </Tooltip>
                </>
              )}

              {item.set_divider && <ToolbarDivider />}
            </>
          ))}
        </Toolbar>
      </div>
    </>
  );
};

export default Navbar;
