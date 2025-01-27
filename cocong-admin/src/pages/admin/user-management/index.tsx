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
  EditRegular,
  TableRegular,
  DeleteRegular,
  KeyResetRegular,
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
  const [dataDirectory, setDataDirectory] = useState([] as any[]);

  const [pageCount, setPageCount] = useState(1);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const limit = 10;
  const apiConnector = useApiConnector();
  const [name, setName] = useState("");
  const [userid, setUserID] = useState("");
  const [username, setUsername] = useState("");
  const [statusEnable, setStatusEnable] = useState(true);

  const [searchDirectory, setSearchDirectory] = useState("");
  const [addBtnDisabled, setAddBtnDisabled] = useState(false);
  const [id, setId] = useState("");
  const toast = useToast();
  const [refresh, setRefresh] = useState(false);
  const setConfirmDialog = useConfirmDialog();

  useEffect(() => {
    apiConnector
      .get("/op/user/query", {
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
      .get("/op/directory/query", {
        params: {
          offset: 0,
          limit: 8,
          search: searchDirectory,
        },
      })
      .then((res) => {
        setDataDirectory(res.data.Data);
      })
      .catch(() => {});
  }, [searchDirectory]);

  const columns = [
    { columnKey: "Name", label: "Name" },
    { columnKey: "Username", label: "Username" },
    { columnKey: "UserID", label: "UserID" },
    { columnKey: "Status", label: "Status" },
    { columnKey: "Directory", label: "Directory Name" },
    // { columnKey: "UploadLimit", label: "Upload Limit" },
  ];

  const clearInput = () => {
    setName("");
    setId("");
    setUserID("");
    setUsername("");
    setSelectedRows(new Set<TableRowId>([]));
  };

  type DirectoryItem = {
    ID: string;
    Name: string;
    Enable: boolean;
  };

  const directoryColumns: TableColumnDefinition<DirectoryItem>[] = [
    createTableColumn<DirectoryItem>({
      columnId: "Name",
    }),
    createTableColumn<DirectoryItem>({
      columnId: "Enable",
    }),
  ];

  const [selectedRows, setSelectedRows] = useState(
    () => new Set<TableRowId>([])
  );
  const {
    getRows,
    selection: { toggleRow, isRowSelected },
  } = useTableFeatures(
    {
      columns: directoryColumns,
      items: dataDirectory,
    },
    [
      useTableSelection({
        selectionMode: "single",
        defaultSelectedItems: selectedRows,
        onSelectionChange: (e, data) => setSelectedRows(data.selectedItems),
      }),
    ]
  );

  const directoryRows = getRows((row) => {
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

  const createSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (selectedRows.size == 0) {
      toast("Warning", <>Please select directory</>, "warning");
    } else {
      let index = selectedRows.values().next().value;
      if (index != undefined) {
        apiConnector
          .post("/op/user/create", {
            Name: name,
            UserID: userid,
            Username: username,
            Enable: statusEnable,
            DirectoryID: dataDirectory[index as number].ID,
          })
          .then(() => {
            toast("Success", <>Create a new user success</>, "success");
            setRefresh(!refresh);
            setMode("data");
            clearInput();
          })
          .catch(() => {
            toast("Error", <>Cannot create a user directory</>, "error");
          });
      }
    }
  };

  const updateSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    apiConnector
      .put("/op/user/update", {
        ID: id,
        Name: name,
        UserID: userid,
        Enable: statusEnable,
      })
      .then(() => {
        toast("Success", <>Update a user success</>, "success");
        setRefresh(!refresh);
        setMode("data");
        clearInput();
        setAddBtnDisabled(false);
      })
      .catch(() => {
        toast("Error", <>Cannot update a user</>, "error");
      });
  };

  const deleteUser = (id: string) => {
    apiConnector
      .delete("/op/user/delete", {
        params: {
          id: id,
        },
      })
      .then(() => {
        toast("Success", <>Delete a user success</>, "success");
        setRefresh(!refresh);
        setMode("data");
        clearInput();
      })
      .catch(() => {
        toast("Error", <>Cannot delete a user</>, "error");
      });
  };
  const resetPasswordUser = (id: string) => {
    apiConnector
      .patch("/op/user/password-reset", {
        ID: id,
      })
      .then(() => {
        toast(
          "Success",
          <>
            Reset a user password success <br /> Default password:{" "}
            <b>P@ssw0rd</b>
          </>,
          "success"
        );
        setRefresh(!refresh);
        setMode("data");
        clearInput();
      })
      .catch(() => {
        toast("Error", <>Cannot reset a user password</>, "error");
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
          User Management{" "}
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
                Add User
              </Tab>
              <Tab
                icon={<AddSquareRegular />}
                value="edit"
                disabled={!addBtnDisabled}
              >
                Edit User
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
                            <TableCell>{item.Username}</TableCell>
                            <TableCell>{item.UserID}</TableCell>
                            <TableCell>
                              {item.Enable ? (
                                <Badge appearance="filled" color="brand">
                                  Enable
                                </Badge>
                              ) : (
                                <Badge appearance="filled" color="informative">
                                  Disable
                                </Badge>
                              )}
                            </TableCell>
                            <TableCell>
                              {item.Directory ? item.Directory.Name : "None"}
                            </TableCell>
                            <TableCell>
                              <TableCellLayout>
                                <Button
                                  icon={<KeyResetRegular />}
                                  aria-label="Reset"
                                  onClick={() => {
                                    setConfirmDialog(
                                      "Reset user password",
                                      <>
                                        Are you sure to reset password of user
                                        name: <b>{item.Name}</b>?
                                      </>,
                                      () => {
                                        resetPasswordUser(item.ID);
                                      },
                                      () => {}
                                    );
                                  }}
                                />
                                <Button
                                  style={{ marginLeft: "0.5rem" }}
                                  icon={<EditRegular />}
                                  aria-label="Edit"
                                  onClick={() => {
                                    setId(item.ID);
                                    setName(item.Name);
                                    setUserID(item.UserID);
                                    setUsername(item.Username);
                                    setMode("edit");
                                    setAddBtnDisabled(true);
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
                                        Are you sure to delete user name:{" "}
                                        <b>{item.Name}</b>?
                                      </>,
                                      () => {
                                        deleteUser(item.ID);
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
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    UserID{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={userid}
                      onChange={(_, d) => {
                        setUserID(d.value);
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
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    Username{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={username}
                      onChange={(_, d) => {
                        setUsername(d.value);
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
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    Default Password{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div
                    style={{
                      width: "383px",
                      border: "1px solid #d1d1d1",
                      borderRadius: "4px",
                      lineHeight: "32px",
                      paddingLeft: "1rem",
                      background: "#e2edf9",
                    }}
                  >
                    P@ssw0rd
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
                    Status{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Checkbox
                      checked={statusEnable}
                      onChange={(ev, data) =>
                        setStatusEnable(data.checked as boolean)
                      }
                      label="Enable"
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
                    Directory{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "calc(60% - 180px)" }}>
                    <Input
                      placeholder="Search (use snakecase for field), example: name like 100"
                      contentBefore={<SearchRegular />}
                      style={{ width: "400px" }}
                      value={searchDirectory}
                      onChange={(_, d) => {
                        setSearchDirectory(d.value);
                      }}
                    />
                    <Table aria-label="Table with single selection">
                      <TableHeader>
                        <TableRow>
                          <TableSelectionCell type="radio" hidden />
                          <TableHeaderCell>Name</TableHeaderCell>
                          <TableHeaderCell>Status</TableHeaderCell>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {directoryRows.map(
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
                              <TableCell>
                                <TableCellLayout>
                                  {item.Enable ? (
                                    <Badge appearance="filled" color="brand">
                                      Enable
                                    </Badge>
                                  ) : (
                                    <Badge
                                      appearance="filled"
                                      color="informative"
                                    >
                                      Disable
                                    </Badge>
                                  )}
                                </TableCellLayout>
                              </TableCell>
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
                    Username{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={username}
                      required
                      disabled
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
                  style={{ display: "flex", margin: "1rem 0", height: "32px" }}
                >
                  <div
                    style={{
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    UserID{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Input
                      style={{ width: "100%" }}
                      value={userid}
                      onChange={(_, d) => {
                        setUserID(d.value);
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
                      width: "180px",

                      fontWeight: "800",
                      lineHeight: "32px",
                    }}
                  >
                    Status{" "}
                    <span style={{ color: themeColors.colorDanger }}>*</span>
                  </div>
                  <div style={{ width: "400px" }}>
                    <Checkbox
                      checked={statusEnable}
                      onChange={(ev, data) =>
                        setStatusEnable(data.checked as boolean)
                      }
                      label="Enable"
                    />
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
