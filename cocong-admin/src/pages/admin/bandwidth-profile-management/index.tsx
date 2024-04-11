import {
  Avatar,
  Button,
  Field,
  Input,
  Label,
  makeStyles,
  PresenceBadgeStatus,
  shorthands,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableCellLayout,
  TableHeader,
  TableHeaderCell,
  TableRow,
  TabList,
  TabValue,
  Title1,
  tokens,
} from "@fluentui/react-components";
import type { NextPage } from "next";
import Head from "next/head";

import {
  SearchRegular,
  AddSquareRegular,
  EditRegular,
  TableRegular,
  DeleteRegular,
  TagRegular,
} from "@fluentui/react-icons";
import { themeColors } from "../../../constants/theme";
import { FormEvent, useEffect, useState } from "react";
import {
  useApiConnector,
  useConfirmDialog,
  useToast,
} from "../../../utils/AppProperties";
import { Pagination } from "@mui/material";
import { json } from "stream/consumers";

const BW: NextPage = () => {
  const [mode, setMode] = useState<TabValue>("data");

  const [data, setData] = useState([] as any[]);
  const [pageCount, setPageCount] = useState(1);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const limit = 10;
  const apiConnector = useApiConnector();
  const [bwName, setBwName] = useState("");
  const [bwDownloadLimit, setBwDownloadLimit] = useState("");
  const [bwUploadLimit, setBwUploadLimit] = useState("");

  const [addBtnDisabled, setAddBtnDisabled] = useState(false);
  const [bwId, setBwId] = useState("");
  const toast = useToast();
  const [refresh, setRefresh] = useState(false);
  const setConfirmDialog = useConfirmDialog();

  useEffect(() => {
    apiConnector
      .get("/op/bandwidth/query", {
        params: {
          offset: (page - 1) * limit,
          limit: limit,
          search: search,
        },
      })
      .then((res) => {
        setData(res.data.Data);
        setPageCount(Math.ceil(res.data.Count / limit));
      })
      .catch(() => {});
  }, [page, search, refresh]);

  const columns = [
    { columnKey: "Name", label: "Name" },
    { columnKey: "DownloadLimit", label: "Download Limit" },
    { columnKey: "UploadLimit", label: "Upload Limit" },
  ];

  const clearInput = () => {
    setBwName("");
    setBwDownloadLimit("0");
    setBwUploadLimit("0");
    setBwId("");
  };

  const createSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    apiConnector
      .post("/op/bandwidth/create", {
        Name: bwName,
        DownloadLimit: parseInt(bwDownloadLimit),
        UploadLimit: parseInt(bwUploadLimit),
      })
      .then(() => {
        toast(
          "Success",
          <>Create a new bandwidth profile success</>,
          "success"
        );
        setRefresh(!refresh);
        setMode("data");
        clearInput();
      })
      .catch(() => {
        toast("Error", <>Cannot create a new bandwidth profile</>, "error");
      });
  };

  const deleteProfile = (id: string) => {
    apiConnector
      .delete("/op/bandwidth/delete", {
        params: {
          id: id,
        },
      })
      .then(() => {
        toast("Success", <>Delete a bandwidth profile success</>, "success");
        setRefresh(!refresh);
        setMode("data");
        clearInput();
      })
      .catch(() => {
        toast("Error", <>Cannot delete a bandwidth profile</>, "error");
      });
  };

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
          Bandwidth Profile Management{" "}
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
              <Tab icon={<TableRegular />} value="data">
                Data
              </Tab>
              <Tab
                icon={<AddSquareRegular />}
                value="add"
                disabled={addBtnDisabled}
              >
                Add Profile
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
            {mode == "data" && (
              <>
                <div style={{ marginBottom: "1rem" }}>
                  <Input
                    type="text"
                    contentBefore={<SearchRegular />}
                    placeholder="Search (use snakecase for field), example: name like 100  | upload_limit > 500"
                    value={search}
                    onChange={(_, d) => {
                      setSearch(d.value);
                    }}
                    style={{ width: "60%" }}
                  />
                </div>
                {data.length === 0 && (
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
                {data.length > 0 && (
                  <>
                    <Table style={{ width: "100%" }}>
                      <TableHeader>
                        <TableRow>
                          {columns.map((column) => (
                            <TableHeaderCell key={column.columnKey}>
                              <b>{column.label}</b>
                            </TableHeaderCell>
                          ))}
                          <TableHeaderCell />
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {data.map((item) => (
                          <TableRow key={item.ID}>
                            <TableCell>
                              <TableCellLayout
                                media={<Avatar name={item.Name} />}
                              >
                                {item.Name}
                              </TableCellLayout>
                            </TableCell>
                            <TableCell>{item.DownloadLimit}</TableCell>
                            <TableCell>{item.UploadLimit}</TableCell>
                            <TableCell>
                              <TableCellLayout>
                                <Button
                                  icon={<DeleteRegular />}
                                  aria-label="Delete"
                                  onClick={() => {
                                    setConfirmDialog(
                                      "Delete bandwidth profile",
                                      <>
                                        Are you sure to delete bandwidth profile
                                        name: <b>{item.Name}</b>?
                                      </>,
                                      () => {
                                        deleteProfile(item.ID);
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
                        count={pageCount}
                        defaultValue={page}
                        onChange={(_, page) => {
                          setPage(page);
                        }}
                      />
                    </div>
                  </>
                )}
              </>
            )}
            {mode == "add" && (
              <form onSubmit={createSubmit}>
                <div
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "150px",
                      padding: "0 1rem",
                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    Name{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={bwName}
                      onChange={(_, d) => {
                        setBwName(d.value);
                      }}
                      required
                    />
                  </div>
                </div>

                <div
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "150px",
                      padding: "0 1rem",
                      fontWeight: "800",
                    }}
                  >
                    Download Limit <br />
                    <span style={{ fontWeight: "normal", fontSize: "8pt" }}>
                      (Unit: mbps, 0 is unlimited)
                    </span>
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      type="number"
                      style={{ width: "100%" }}
                      value={bwDownloadLimit}
                      onChange={(_, d) => {
                        setBwDownloadLimit(d.value);
                      }}
                      required
                    />
                  </div>
                </div>

                <div
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "150px",
                      padding: "0 1rem",
                      fontWeight: "800",
                    }}
                  >
                    Upload Limit <br />
                    <span style={{ fontWeight: "normal", fontSize: "8pt" }}>
                      (Unit: mbps, 0 is unlimited)
                    </span>
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      type="number"
                      style={{ width: "100%" }}
                      value={bwUploadLimit}
                      onChange={(_, d) => {
                        setBwUploadLimit(d.value);
                      }}
                      required
                    />
                  </div>
                </div>

                <div style={{ marginTop: "2rem", paddingLeft: "1rem" }}>
                  <Button
                    appearance="outline"
                    onClick={() => {
                      setBwName("");
                      setBwDownloadLimit("0");
                      setBwUploadLimit("0");
                    }}
                  >
                    Clear
                  </Button>
                  <Button
                    appearance="primary"
                    style={{ marginLeft: "0.5rem" }}
                    type="submit"
                  >
                    Create
                  </Button>
                </div>
              </form>
            )}
          </div>
        </div>
      </div>
    </>
  );
};

export default BW;
