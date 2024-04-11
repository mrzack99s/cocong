import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";
import { NextRouter, useRouter } from "next/router";
import https from "https";

export class ApiConnector {
  apiKey: string;
  instance: AxiosInstance;
  router: NextRouter;

  constructor(apiUrl: string, apiKey: string, router: NextRouter) {
    this.router = router;
    this.apiKey = apiKey;

    this.instance = axios.create({
      baseURL: apiUrl,
      headers: {
        "api-token": this.apiKey,
        "Cache-Control": "no-cache",
      },
      httpsAgent: new https.Agent({ rejectUnauthorized: false }),
    });
    // this.instance.interceptors.response.use(
    //   (response) => response,
    //   (error) => {
    //     if (error.response.status === 401) {
    //       if(!!error.response.data && error.response.data === "invalid token") {
    //         location.href = "/api/oauth2/refresh-token"
    //       }else{
    //         this.router.push("/unauthorized");
    //       }

    //     }
    //   }
    // );
  }
}
