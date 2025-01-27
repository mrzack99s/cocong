import {
  Avatar,
  Button,
  createTableColumn,
  Field,
  Input,
  Label,
  makeStyles,
  PresenceBadgeStatus,
  Badge,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableCellLayout,
  TableColumnDefinition,
  TableHeader,
  TableHeaderCell,
  TableRow,
  TableRowId,
  TableSelectionCell,
  TabList,
  TabValue,
  Title1,
  tokens,
  Checkbox,
  useTableSelection,
  useTableFeatures,
} from "@fluentui/react-components";
import type { NextPage } from "next";
import Head from "next/head";

import {
  SearchRegular,
  AddSquareRegular,
  ArrowClockwiseRegular,
  TableRegular,
  DeleteRegular,
  ArrowDownloadRegular,
} from "@fluentui/react-icons";
import { themeColors } from "../../constants/theme";
import { FormEvent, useEffect, useState } from "react";
import {
  useApiConnector,
  useConfirmDialog,
  useProfile,
  useToast,
} from "../../utils/AppProperties";
import { Autocomplete, Pagination, TextField } from "@mui/material";
import { json } from "stream/consumers";
import { useFile } from "../../utils/Facilities";

const BW: NextPage = () => {
  const [mode, setMode] = useState<TabValue>("online-session");
  const toast = useToast();
  const setConfirmDialog = useConfirmDialog();
  const apiConnector = useApiConnector();
  const limit = 20;

  // Online-sesion
  const [dataOnlineSession, setDataOnlineSession] = useState([] as any[]);
  const [pageCountOnlineSession, setPageCountOnlineSession] = useState(1);
  const [countOnlineSession, setCountOnlineSession] = useState(0);
  const [searchOnlineSession, setSearchOnlineSession] = useState("");
  const [pageOnlineSession, setPageOnlineSession] = useState(1);
  const [refreshOnlineSession, setRefreshOnlineSession] = useState(false);

  useEffect(() => {
    if (mode == "online-session") {
      apiConnector
        .get("/op/session/query", {
          params: {
            offset: (pageOnlineSession - 1) * limit,
            limit: limit,
            search: searchOnlineSession,
          },
        })
        .then((res) => {
          setDataOnlineSession(res.data.Data);
          setCountOnlineSession(res.data.Count);
          setPageCountOnlineSession(Math.ceil(res.data.Count / limit));
        })
        .catch(() => {});
    }
  }, [pageOnlineSession, searchOnlineSession, refreshOnlineSession]);

  // login-log
  const [dataLoginLog, setDataLoginLog] = useState([] as any[]);
  const [pageCountLoginLog, setPageCountLoginLog] = useState(1);
  const [searchLoginLog, setSearchLoginLog] = useState("");
  const [pageLoginLog, setPageLoginLog] = useState(1);
  const [refreshLoginLog, setRefreshLoginLog] = useState(false);

  useEffect(() => {
    if (mode == "login-log") {
      apiConnector
        .get("/op/log/login", {
          params: {
            offset: (pageLoginLog - 1) * limit,
            limit: limit,
            search: searchLoginLog,
          },
        })
        .then((res) => {
          setDataLoginLog(res.data.Data);
          setPageCountLoginLog(Math.ceil(res.data.Count / limit));
        })
        .catch(() => {});
    }
  }, [pageLoginLog, searchLoginLog, refreshLoginLog]);

  // logout-log
  const [dataLogoutLog, setDataLogoutLog] = useState([] as any[]);
  const [pageCountLogoutLog, setPageCountLogoutLog] = useState(1);
  const [searchLogoutLog, setSearchLogoutLog] = useState("");
  const [pageLogoutLog, setPageLogoutLog] = useState(1);
  const [refreshLogoutLog, setRefreshLogoutLog] = useState(false);

  useEffect(() => {
    if (mode == "logout-log") {
      apiConnector
        .get("/op/log/logout", {
          params: {
            offset: (pageLogoutLog - 1) * limit,
            limit: limit,
            search: searchLogoutLog,
          },
        })
        .then((res) => {
          setDataLogoutLog(res.data.Data);
          setPageCountLogoutLog(Math.ceil(res.data.Count / limit));
        })
        .catch(() => {});
    }
  }, [pageLogoutLog, searchLogoutLog, refreshLogoutLog]);

  // net-log
  const [dataNetLog, setDataNetLog] = useState([] as any[]);
  const [pageCountNetLog, setPageCountNetLog] = useState(1);
  const [searchNetLog, setSearchNetLog] = useState("");
  const [pageNetLog, setPageNetLog] = useState(1);
  const [refreshNetLog, setRefreshNetLog] = useState(false);

  useEffect(() => {
    switch (mode) {
      case "online-session":
        setRefreshOnlineSession((r) => !r);
        break;
      case "login-log":
        setRefreshLoginLog((r) => !r);
        break;
      case "logout-log":
        setRefreshLogoutLog((r) => !r);
        break;
      case "net-log":
        setRefreshNetLog((r) => !r);
        break;
    }
  }, [mode]);

  useEffect(() => {
    if (mode == "net-log") {
      apiConnector
        .get("/op/log/net", {
          params: {
            offset: (pageNetLog - 1) * limit,
            limit: limit,
            search: searchNetLog,
          },
        })
        .then((res) => {
          setDataNetLog(res.data.Data);
          setPageCountNetLog(Math.ceil(res.data.Count / limit));
        })
        .catch(() => {});
    }
  }, [pageNetLog, searchNetLog, refreshNetLog]);

  const kickSession = (id: string) => {
    apiConnector
      .patch("/op/session/kick", {
        SessionID: id,
      })
      .then(() => {
        toast("Success", <>Kick a user success</>, "success");
        setRefreshOnlineSession(!refreshOnlineSession);
      })
      .catch(() => {
        toast("Error", <>Cannot kick a user</>, "error");
      });
  };

  const profile = useProfile();
  const [doSaveAs] = useFile();
  return (
    <>
      <Head>
        <title>COCONG | Admin Centre</title>
      </Head>

      <div style={{ margin: "2rem 1rem" }}>
        <span
          style={{
            background: themeColors.colorSecondary,
            fontSize: "16pt",
            width: "auto",
            //   position: "absolute",
            padding: "0.5rem 1rem ",
            color: "#fff",
            borderRadius: "8px",
          }}
        >
          Monitors{" "}
        </span>
      </div>

      <div className="ms-Grid" dir="ltr" style={{ margin: "1rem" }}>
        <div className="ms-Grid-row">
          <div className="ms-Grid-col ms-hiddenSm ms-md3 ms-lg3 ms-xl2 ms-xxl2">
            <TabList
              defaultSelectedValue="tab2"
              size="large"
              vertical
              selectedValue={mode}
              onTabSelect={(_, d) => {
                setMode(d.value);
              }}
            >
              <Tab icon={<TableRegular />} value="online-session">
                Online Sessions
              </Tab>
              <Tab icon={<TableRegular />} value="login-log">
                Logged in logs
              </Tab>
              <Tab icon={<TableRegular />} value="logout-log">
                Logged out logs
              </Tab>
              <Tab icon={<TableRegular />} value="net-log">
                Network logs
              </Tab>
            </TabList>
          </div>
          <div
            className="ms-Grid-col ms-sm12 ms-md9 ms-lg9 ms-xl10 ms-xxl10"
            style={{
              border: "1px solid #f5f5f5",
              borderRadius: "12px",
              padding: "1rem 1rem 2rem 1rem",
            }}
          >
            {mode == "online-session" && (
              <>
                <div style={{ display: "flex", width: "100%" }}>
                  <div
                    style={{
                      width: "260px",
                      border: "1px solid #f5f5f5",
                      borderRadius: "12px",
                      padding: "1rem 1rem 0 1rem",
                      marginBottom: "1rem",
                    }}
                  >
                    <div>
                      <b>Online Sessions</b>
                    </div>
                    <div
                      style={{
                        display: "flex",
                        fontSize: "14pt",
                        height: "64px",
                        justifyContent: "center",
                        alignItems: "center",
                      }}
                    >
                      {countOnlineSession.toLocaleString()}
                    </div>
                  </div>
                  <div
                    style={{
                      maxWidth: "260px",
                      padding: "0 1rem",
                      marginBottom: "1rem",
                    }}
                  >
                    <Button
                      style={{ height: "100%" }}
                      icon={<ArrowClockwiseRegular />}
                      onClick={() => {
                        setDataOnlineSession([]);
                        setRefreshOnlineSession(!refreshOnlineSession);
                      }}
                    >
                      Refresh
                    </Button>
                  </div>
                </div>

                <div style={{ marginBottom: "1rem", height: "32px" }}>
                  <Input
                    type="text"
                    contentBefore={<SearchRegular />}
                    placeholder="Search"
                    value={searchOnlineSession}
                    onChange={(_, d) => {
                      setSearchOnlineSession(d.value);
                    }}
                    style={{ width: "60%" }}
                  />
                </div>
                {dataOnlineSession.length === 0 && (
                  <div
                    style={{
                      textAlign: "left",
                      fontSize: "16pt",
                      padding: "2rem",
                    }}
                  >
                    Not found any data in database
                  </div>
                )}
                {dataOnlineSession.length > 0 && (
                  <>
                    <Table style={{ width: "100%" }} noNativeElements>
                      <TableHeader>
                        <TableRow>
                          <TableHeaderCell>
                            <b>User</b>
                          </TableHeaderCell>
                          <TableHeaderCell style={{ maxWidth: "200px" }}>
                            <b>IP Address</b>
                          </TableHeaderCell>
                          <TableHeaderCell style={{ maxWidth: "120px" }}>
                            <b>Auth Type</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>Last Seen</b>
                          </TableHeaderCell>

                          <TableHeaderCell style={{ maxWidth: "60px" }} />
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {dataOnlineSession.map((item, index) => (
                          <TableRow key={item.ID}>
                            <TableCell>
                              <TableCellLayout
                                media={<Avatar name={item.User} />}
                              >
                                {item.User}
                              </TableCellLayout>
                            </TableCell>
                            <TableCell style={{ maxWidth: "200px" }}>
                              {item.IPAddress}
                            </TableCell>
                            <TableCell style={{ maxWidth: "120px" }}>
                              {item.AuthType}
                            </TableCell>
                            <TableCell>
                              {new Date(item.LastSeen).toLocaleString()}
                            </TableCell>
                            <TableCell style={{ maxWidth: "60px" }}>
                              <TableCellLayout>
                                <Button
                                  style={{ marginLeft: "0.5rem" }}
                                  icon={<DeleteRegular />}
                                  aria-label="Delete"
                                  disabled={profile.ID === item.ID}
                                  onClick={() => {
                                    setConfirmDialog(
                                      "Kick session",
                                      <>
                                        Are you sure to kick session of user:{" "}
                                        <b>{item.User}</b>?
                                      </>,
                                      () => {
                                        kickSession(item.ID);
                                      },
                                      () => {}
                                    );
                                  }}
                                />
                              </TableCellLayout>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>

                    <div style={{ marginTop: "1rem" }}>
                      <Pagination
                        count={pageCountOnlineSession}
                        defaultValue={pageOnlineSession}
                        onChange={(_, page) => {
                          setPageOnlineSession(page);
                        }}
                      />
                    </div>
                  </>
                )}
              </>
            )}

            {mode == "login-log" && (
              <>
                <Button
                  style={{ marginBottom: "1rem" }}
                  icon={<ArrowClockwiseRegular />}
                  onClick={() => {
                    setDataLoginLog([]);
                    setRefreshLoginLog(!refreshLoginLog);
                  }}
                >
                  Refresh
                </Button>

                <Button
                  style={{ marginBottom: "1rem", marginLeft: "0.25rem" }}
                  icon={<ArrowDownloadRegular />}
                  onClick={() => {
                    doSaveAs(
                      `logged-in-log.json`,
                      JSON.stringify(dataLoginLog)
                    );
                  }}
                >
                  Download JSON
                </Button>

                <div style={{ marginBottom: "1rem", height: "32px" }}>
                  <Input
                    type="text"
                    contentBefore={<SearchRegular />}
                    placeholder="Search"
                    value={searchLoginLog}
                    onChange={(_, d) => {
                      setSearchLoginLog(d.value);
                    }}
                    style={{ width: "60%" }}
                  />
                </div>
                {dataLoginLog.length === 0 && (
                  <div
                    style={{
                      textAlign: "left",
                      fontSize: "16pt",
                      padding: "2rem",
                    }}
                  >
                    Not found any data in memory
                  </div>
                )}
                {dataLoginLog.length > 0 && (
                  <>
                    <Table style={{ width: "100%" }}>
                      <TableHeader>
                        <TableRow>
                          <TableHeaderCell>
                            <b>Transaction At</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>User</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>IP Address</b>
                          </TableHeaderCell>

                          <TableHeaderCell>
                            <b>Success</b>
                          </TableHeaderCell>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {dataLoginLog.map((item, index) => (
                          <TableRow key={item.ID}>
                            <TableCell>
                              {new Date(item.TransactionAt).toLocaleString()}
                            </TableCell>
                            <TableCell>{item.User}</TableCell>
                            <TableCell>{item.IPAddress}</TableCell>

                            <TableCell>
                              {item.Success ? (
                                <Badge appearance="filled" color="brand">
                                  Success
                                </Badge>
                              ) : (
                                <Badge appearance="filled" color="informative">
                                  Failed
                                </Badge>
                              )}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>

                    <div style={{ marginTop: "1rem" }}>
                      <Pagination
                        count={pageCountLoginLog}
                        defaultValue={pageLoginLog}
                        onChange={(_, page) => {
                          setPageLoginLog(page);
                        }}
                      />
                    </div>
                  </>
                )}
              </>
            )}

            {mode == "logout-log" && (
              <>
                <Button
                  style={{ marginBottom: "1rem" }}
                  icon={<ArrowClockwiseRegular />}
                  onClick={() => {
                    setDataLogoutLog([]);
                    setRefreshLogoutLog(!refreshLogoutLog);
                  }}
                >
                  Refresh
                </Button>
                <Button
                  style={{ marginBottom: "1rem", marginLeft: "0.25rem" }}
                  icon={<ArrowDownloadRegular />}
                  onClick={() => {
                    doSaveAs(
                      `logged-out-log.json`,
                      JSON.stringify(dataLogoutLog)
                    );
                  }}
                >
                  Download JSON
                </Button>

                <div style={{ marginBottom: "1rem", height: "32px" }}>
                  <Input
                    type="text"
                    contentBefore={<SearchRegular />}
                    placeholder="Search"
                    value={searchLogoutLog}
                    onChange={(_, d) => {
                      setSearchLogoutLog(d.value);
                    }}
                    style={{ width: "60%" }}
                  />
                </div>
                {dataLogoutLog.length === 0 && (
                  <div
                    style={{
                      textAlign: "left",
                      fontSize: "16pt",
                      padding: "2rem",
                    }}
                  >
                    Not found any data in database
                  </div>
                )}
                {dataLogoutLog.length > 0 && (
                  <>
                    <Table style={{ width: "100%" }}>
                      <TableHeader>
                        <TableRow>
                          <TableHeaderCell>
                            <b>Transaction At</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>User</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>IP Address</b>
                          </TableHeaderCell>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {dataLogoutLog.map((item, index) => (
                          <TableRow key={item.ID}>
                            <TableCell>
                              {new Date(item.TransactionAt).toLocaleString()}
                            </TableCell>
                            <TableCell>{item.User}</TableCell>
                            <TableCell>{item.IPAddress}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>

                    <div style={{ marginTop: "1rem" }}>
                      <Pagination
                        count={pageCountLogoutLog}
                        defaultValue={pageLogoutLog}
                        onChange={(_, page) => {
                          setPageLogoutLog(page);
                        }}
                      />
                    </div>
                  </>
                )}
              </>
            )}

            {mode == "net-log" && (
              <>
                <Button
                  style={{ marginBottom: "1rem" }}
                  icon={<ArrowClockwiseRegular />}
                  onClick={() => {
                    setDataNetLog([]);
                    setRefreshNetLog(!refreshNetLog);
                  }}
                >
                  Refresh
                </Button>

                <Button
                  style={{ marginBottom: "1rem", marginLeft: "0.25rem" }}
                  icon={<ArrowDownloadRegular />}
                  onClick={() => {
                    doSaveAs(`net-log.json`, JSON.stringify(dataNetLog));
                  }}
                >
                  Download JSON
                </Button>

                <div style={{ marginBottom: "1rem", height: "32px" }}>
                  <Input
                    type="text"
                    contentBefore={<SearchRegular />}
                    placeholder="Search"
                    value={searchNetLog}
                    onChange={(_, d) => {
                      setSearchNetLog(d.value);
                    }}
                    style={{ width: "60%" }}
                  />
                </div>
                {dataNetLog.length === 0 && (
                  <div
                    style={{
                      textAlign: "left",
                      fontSize: "16pt",
                      padding: "2rem",
                    }}
                  >
                    Not found any data in database
                  </div>
                )}
                {dataNetLog.length > 0 && (
                  <>
                    <Table style={{ width: "100%" }}>
                      <TableHeader>
                        <TableRow>
                          <TableHeaderCell>
                            <b>Transaction At</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>Protocol</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>Source</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>Destination</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>Traffic From Internet</b>
                          </TableHeaderCell>
                          <TableHeaderCell>
                            <b>Source</b>
                          </TableHeaderCell>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {dataNetLog.map((item, index) => (
                          <TableRow key={item.ID}>
                            <TableCell>
                              {new Date(item.TransactionAt).toLocaleString()}
                            </TableCell>
                            <TableCell>{item.Protocol}</TableCell>
                            <TableCell>
                              {item.SourceNetwork}, {item.SourcePort}
                            </TableCell>
                            <TableCell>
                              {item.DestinationNetwork}, {item.DestinationPort}
                            </TableCell>
                            <TableCell>
                              {item.TrafficFromInternet ? "True" : "False"}
                            </TableCell>
                            <TableCell>{item.User}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>

                    <div style={{ marginTop: "1rem" }}>
                      <Pagination
                        count={pageCountNetLog}
                        defaultValue={pageNetLog}
                        onChange={(_, page) => {
                          setPageNetLog(page);
                        }}
                      />
                    </div>
                  </>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </>
  );
};

export default BW;
