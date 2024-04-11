import {
  useApiConnector,
  usePageLoading,
  useToast,
} from "../utils/AppProperties";
import { Button, Input } from "@fluentui/react-components";

import { PersonRegular, PasswordRegular } from "@fluentui/react-icons";
import type { NextPage } from "next";
import Head from "next/head";
import { useRouter } from "next/router";
import { FormEvent, useState } from "react";
import { useCookies } from "react-cookie";

const Home: NextPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const apiConnector = useApiConnector();
  const toast = useToast();
  const router = useRouter();
  const [cookies, setCookie, removeCookie] = useCookies([
    "api-token",
    "refresh-token",
  ]);

  const [closePageLoading, setPageLoading] = usePageLoading();

  const submit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setPageLoading("Logging in to system");

    apiConnector
      .post("/op/login", {
        username: username,
        password: password,
      })
      .then((res) => {
        removeCookie("api-token");
        setCookie("api-token", res.data.AccessToken, {
          expires: new Date(res.data.AccessTokenExpired),
        });

        removeCookie("refresh-token");
        setCookie("refresh-token", res.data.RefreshToken, {
          expires: new Date(res.data.RefreshTokenExpired),
        });

        router.push("/admin");
      })
      .catch((err) => {
        toast("Error", <>{err.response.data}</>, "error");
      })
      .finally(() => {
        closePageLoading();
      });
  };

  return (
    <>
      <Head>
        <title>COCONG | Login</title>
      </Head>

      <div
        style={{
          width: "100vw",
          height: "100vh",
          background: "#fefefe",
        }}
      >
        <div
          style={{
            width: "300px",
            height: "400px",
            position: "absolute",
            top: "40%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            textAlign: "center",
          }}
        >
          <div
            style={{
              width: "100px",
              height: "100px",
              background: "#f0f0f0",
              textAlign: "center",
              marginTop: "2rem",
              fontSize: "32pt",
              margin: "1rem auto",
              borderRadius: "70px",
              position: "relative",
              left: "10%",
              padding: "1rem",
            }}
          >
            <div
              style={{
                fontSize: "24pt",
                fontWeight: 800,
                width: "77%",
                marginBottom: "2rem",
                color: "#59a4ff",
                position: "absolute",
                top: "40%",
                textAlign: "center",
              }}
            >
              COCO
            </div>
          </div>
          <div
            style={{
              marginTop: "1rem",
              width: "100%",
              height: "200px",
              borderRadius: "12px",
              padding: "2rem",
            }}
          >
            <form onSubmit={submit}>
              <Input
                placeholder="Username"
                contentBefore={<PersonRegular />}
                style={{ width: "100%" }}
                required
                onChange={(_, d) => {
                  setUsername(d.value);
                }}
              />
              <Input
                placeholder="Password"
                type="password"
                contentBefore={<PasswordRegular />}
                style={{ width: "100%", marginTop: "1rem" }}
                required
                onChange={(_, d) => {
                  setPassword(d.value);
                }}
              />
              <Button
                style={{ marginTop: "1rem" }}
                type="submit"
                appearance="primary"
              >
                Login
              </Button>
            </form>
          </div>
        </div>
      </div>
    </>
  );
};

export default Home;
