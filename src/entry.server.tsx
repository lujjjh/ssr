import React from "react";
import { renderToString } from "react-dom/server";

export default () => {
  return renderToString(<div>Hello SSR</div>);
};
