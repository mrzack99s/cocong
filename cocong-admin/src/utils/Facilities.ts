export const sleep = (delay: number) =>
  new Promise((resolve) => setTimeout(resolve, delay));
export const useFile = () => {
  const doSaveAs = (fileName: string, text: string, type?: string) => {
    var element = document.createElement("a");
    if (!!type) {
      element.setAttribute(
        "href",
        `data:${type};charset=utf-8,` + encodeURIComponent(text)
      );
    } else {
      element.setAttribute(
        "href",
        "data:text/plain;charset=utf-8," + encodeURIComponent(text)
      );
    }

    element.setAttribute("download", fileName);
    element.style.display = "none";
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  };
  return [doSaveAs];
};
