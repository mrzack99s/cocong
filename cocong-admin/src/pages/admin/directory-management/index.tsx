import {
  Avatar,
  Button,
  createTableColumn,
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
  useTableFeatures,
  useTableSelection,
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
import { Autocomplete, Pagination, TextField } from "@mui/material";
import { json } from "stream/consumers";

const BW: NextPage = () => {
  const [mode, setMode] = useState<TabValue>("data");

  const [data, setData] = useState([] as any[]);
  const [dataBW, setBWData] = useState([] as any[]);

  const [pageCount, setPageCount] = useState(1);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const limit = 10;
  const apiConnector = useApiConnector();
  const [name, setName] = useState("");

  const [searchBWProfile, setSearchBWProfile] = useState("");
  const [addBtnDisabled, setAddBtnDisabled] = useState(false);
  const [id, setId] = useState("");
  const toast = useToast();
  const [refresh, setRefresh] = useState(false);
  const setConfirmDialog = useConfirmDialog();

  useEffect(() => {
    apiConnector
      .get("/op/directory/query", {
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

  useEffect(() => {
    apiConnector
      .get("/op/bandwidth/query", {
        params: {
          offset: 0,
          limit: 8,
          search: searchBWProfile,
        },
      })
      .then((res) => {
        setBWData(res.data.Data);
      })
      .catch(() => {});
  }, [searchBWProfile]);

  const columns = [
    { columnKey: "Name", label: "Name" },
    { columnKey: "Bandwidth", label: "Bandwidth Profile" },
  ];

  const clearInput = () => {
    setName("");
    setId("");
    setBwSelectedRows(new Set<TableRowId>([]));
  };

  const createSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (bwSelectedRows.size == 0) {
      toast("Warning", <>Please select bandwidth profile</>, "warning");
    } else {
      apiConnector
        .post("/op/directory/create", {
          Name: name,
          BandwidthID: dataBW[bwSelectedRows.values().next().value].ID,
        })
        .then(() => {
          toast("Success", <>Create a new directory success</>, "success");
          setRefresh(!refresh);
          setMode("data");
          clearInput();
        })
        .catch(() => {
          toast("Error", <>Cannot create a new directory</>, "error");
        });
    }
  };

  const updateSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    apiConnector
      .put("/op/directory/update", {
        ID: id,
        Name: name,
        BandwidthID: dataBW[bwSelectedRows.values().next().value].ID,
      })
      .then(() => {
        toast("Success", <>Update a directory success</>, "success");
        setRefresh(!refresh);
        setMode("data");
        clearInput();
      })
      .catch(() => {
        toast("Error", <>Cannot update a directory</>, "error");
      })
      .finally(() => {
        setAddBtnDisabled(false);
      });
  };

  const deleteDirectory = (id: string) => {
    apiConnector
      .delete("/op/directory/delete", {
        params: {
          id: id,
        },
      })
      .then(() => {
        toast("Success", <>Delete a directory success</>, "success");
        setRefresh(!refresh);
        setMode("data");
        clearInput();
      })
      .catch(() => {
        toast("Error", <>Cannot delete a directory</>, "error");
      });
  };

  type BwItem = {
    ID: string;
    Name: string;
    DownloadLimit: string;
    UploadLimit: string;
  };

  const bwColumns: TableColumnDefinition<BwItem>[] = [
    createTableColumn<BwItem>({
      columnId: "Name",
    }),
    createTableColumn<BwItem>({
      columnId: "DownloadLimit",
    }),
    createTableColumn<BwItem>({
      columnId: "UploadLimit",
    }),
  ];

  const [bwSelectedRows, setBwSelectedRows] = useState(
    () => new Set<TableRowId>([])
  );
  const {
    getRows,
    selection: { toggleRow, isRowSelected },
  } = useTableFeatures(
    {
      columns: bwColumns,
      items: dataBW,
    },
    [
      useTableSelection({
        selectionMode: "single",
        defaultSelectedItems: bwSelectedRows,
        onSelectionChange: (e, data) => setBwSelectedRows(data.selectedItems),
      }),
    ]
  );

  const bwRows = getRows((row) => {
    const selected = isRowSelected(row.rowId);
    return {
      ...row,
      onClick: (e: React.MouseEvent) => toggleRow(e, row.rowId),
      onKeyDown: (e: React.KeyboardEvent) => {
        if (e.key === " ") {
          e.preventDefault();
          toggleRow(e, row.rowId);
        }
      },
      selected,
      appearance: selected ? ("brand" as const) : ("none" as const),
    };
  });

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
          Directory Management{" "}
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
                Add Directory
              </Tab>
              <Tab
                icon={<AddSquareRegular />}
                value="edit"
                disabled={!addBtnDisabled}
              >
                Edit Directory
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
                    placeholder="Search (use snakecase for field), example: name like test"
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
                        {data.map((item, index) => (
                          <TableRow key={item.ID}>
                            <TableCell>
                              <TableCellLayout
                                media={<Avatar name={item.Name} />}
                              >
                                {item.Name}
                              </TableCellLayout>
                            </TableCell>
                            <TableCell>
                              {item.Bandwidth ? item.Bandwidth.Name : "None"}
                            </TableCell>
                            <TableCell>
                              <TableCellLayout>
                                <Button
                                  icon={<EditRegular />}
                                  aria-label="Edit"
                                  onClick={() => {
                                    setId(item.ID);
                                    setName(item.Name);
                                    setMode("edit");
                                    setAddBtnDisabled(true);
                                    setBwSelectedRows(new Set<TableRowId>([]));
                                  }}
                                />
                                <Button
                                  style={{ marginLeft: "0.5rem" }}
                                  icon={<DeleteRegular />}
                                  aria-label="Delete"
                                  onClick={() => {
                                    setConfirmDialog(
                                      "Delete directory",
                                      <>
                                        Are you sure to delete directory name:{" "}
                                        <b>{item.Name}</b>?
                                      </>,
                                      () => {
                                        deleteDirectory(item.ID);
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
                  style={{
                    display: "flex",
                    margin: "1rem 0",
                    height: "32px",
                    width: "100%",
                  }}
                >
                  <div
                    style={{
                      width: "180px",

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
                      value={name}
                      onChange={(_, d) => {
                        setName(d.value);
                      }}
                      required
                    />
                  </div>
                </div>

                <div
                  style={{ display: "flex", margin: "1rem 0", width: "100%" }}
                >
                  <div
                    style={{
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    Bandwidth Profile{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "calc(100% - 180px)" }}>
                    <Input
                      placeholder="Search (use snakecase for field), example: name like 100"
                      contentBefore={<SearchRegular />}
                      style={{ width: "400px" }}
                      value={searchBWProfile}
                      onChange={(_, d) => {
                        setSearchBWProfile(d.value);
                      }}
                    />
                    <Table aria-label="Table with single selection">
                      <TableHeader>
                        <TableRow>
                          <TableSelectionCell type="radio" hidden />
                          <TableHeaderCell>Name</TableHeaderCell>
                          <TableHeaderCell>
                            Download Limit (mbps)
                          </TableHeaderCell>
                          <TableHeaderCell>Upload Limit (mbps)</TableHeaderCell>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {bwRows.map(
                          ({
                            item,
                            selected,
                            onClick,
                            onKeyDown,
                            appearance,
                          }) => (
                            <TableRow
                              key={item.ID}
                              onClick={onClick}
                              onKeyDown={onKeyDown}
                              aria-selected={selected}
                              appearance={appearance}
                            >
                              <TableSelectionCell
                                checked={selected}
                                type="radio"
                                radioIndicator={{ "aria-label": "Select row" }}
                              />
                              <TableCell>
                                <TableCellLayout>{item.Name}</TableCellLayout>
                              </TableCell>

                              <TableCell>{item.DownloadLimit}</TableCell>
                              <TableCell>{item.UploadLimit}</TableCell>
                            </TableRow>
                          )
                        )}
                      </TableBody>
                    </Table>
                  </div>
                </div>

                <div style={{ marginTop: "2rem", paddingLeft: "1rem" }}>
                  <Button
                    appearance="outline"
                    onClick={() => {
                      clearInput();
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

            {mode == "edit" && (
              <form onSubmit={updateSubmit}>
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
                    Name{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={name}
                      onChange={(_, d) => {
                        setName(d.value);
                      }}
                      required
                    />
                  </div>
                </div>

                <div style={{ display: "flex", margin: "1rem 0" }}>
                  <div
                    style={{
                      width: "180px",
                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    Bandwidth Profile{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "calc(100% - 180px)" }}>
                    <Input
                      placeholder="Search (use snakecase for field), example: name like 100"
                      contentBefore={<SearchRegular />}
                      style={{ width: "400px" }}
                      value={searchBWProfile}
                      onChange={(_, d) => {
                        setSearchBWProfile(d.value);
                      }}
                    />
                    <Table aria-label="Table with single selection">
                      <TableHeader>
                        <TableRow>
                          <TableSelectionCell type="radio" hidden />
                          <TableHeaderCell>Name</TableHeaderCell>
                          <TableHeaderCell>
                            Download Limit (mbps)
                          </TableHeaderCell>
                          <TableHeaderCell>Upload Limit (mbps)</TableHeaderCell>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {bwRows.map(
                          ({
                            item,
                            selected,
                            onClick,
                            onKeyDown,
                            appearance,
                          }) => (
                            <TableRow
                              key={item.ID}
                              onClick={onClick}
                              onKeyDown={onKeyDown}
                              aria-selected={selected}
                              appearance={appearance}
                            >
                              <TableSelectionCell
                                checked={selected}
                                type="radio"
                                radioIndicator={{ "aria-label": "Select row" }}
                              />
                              <TableCell>
                                <TableCellLayout>{item.Name}</TableCellLayout>
                              </TableCell>

                              <TableCell>{item.DownloadLimit}</TableCell>
                              <TableCell>{item.UploadLimit}</TableCell>
                            </TableRow>
                          )
                        )}
                      </TableBody>
                    </Table>
                  </div>
                </div>

                <div style={{ marginTop: "2rem", paddingLeft: "1rem" }}>
                  <Button
                    appearance="outline"
                    onClick={() => {
                      clearInput();
                      setMode("data");
                      setAddBtnDisabled(false);
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    appearance="primary"
                    style={{ marginLeft: "0.5rem" }}
                    type="submit"
                  >
                    Update
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
