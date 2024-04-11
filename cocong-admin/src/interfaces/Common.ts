import { Dispatch, ReactElement, SetStateAction } from "react";

export interface IMenuDrawer {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
}

export interface ILoadingPage extends IMenuDrawer {
  message?: string;
}

export interface IFormValidate {
  [keys: string]: {
    duplicate?: boolean;
    valid: boolean;
    msg?: string;
  };
}

export interface IConfirmDialog {
  open?: boolean;
  message: ReactElement;
  title: string;
  onAccept: () => void;
  onReject: () => void;
}

export interface IPageMenu {
  menu: {
    name: string;
    icon: JSX.Element;
    description: string;
    onClick: () => void;
  }[];
}

export interface IAddProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  refresh: () => void;
}

export interface IEditProps {
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  refresh: () => void;
  onClose?: () => void;
  data: any;
}
