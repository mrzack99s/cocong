import { IConfirmDialog } from "../interfaces/Common";
import {
  Dialog,
  DialogTrigger,
  Button,
  DialogSurface,
  DialogBody,
  DialogTitle,
  DialogContent,
  DialogActions,
} from "@fluentui/react-components";

const ConfirmDialog = (props: IConfirmDialog) => {
  return (
    <Dialog modalType="alert" open={props.open}>
      <DialogSurface>
        <DialogBody>
          <DialogTitle>{props.title}</DialogTitle>
          <DialogContent>{props.message}</DialogContent>
          <DialogActions fluid>
            <Button appearance="secondary" onClick={props.onReject}>
              Cancel
            </Button>
            <Button appearance="primary" onClick={props.onAccept}>
              Confirm
            </Button>
          </DialogActions>
        </DialogBody>
      </DialogSurface>
    </Dialog>
  );
};

export default ConfirmDialog;
